#!/usr/bin/env bash

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")"
cd "${SCRIPT_ROOT}"

set -euo pipefail

export P4USER="${P4USER:-"admin"}"                   # the name of the Perforce superuser that the script will use to create the depot
export P4PORT="${P4PORT:-"perforce.sgdev.org:1666"}" # the address of the Perforce server to connect to

export P4_TEST_USERNAME="${P4_TEST_USERNAME:-"test-perforce"}"           # the username of the fake user that the script will create for integration testing purposes
export P4_TEST_EMAIL="${P4_TEST_EMAIL:-"test-perforce@sourcegraph.com"}" # the email address of the fake user that the script will create for integration testing purposes

export DEPOT_NAME="${DEPOT_NAME:-"integration-test-depot"}" # the name of the depot that the script will create on the server

declare -A dependencies=(
  ["p4"]="$(printf "Please install '%s' by:\n\t- (macOS): running %s\n\t- (Linux): installing it via your distribution's package manager\nSee %s for more information.\n" \
    "p4" \
    "brew install p4" \
    "https://www.perforce.com/downloads/helix-command-line-client-p4")"

  ["fzf"]="$(printf "Please install '%s' by:\n\t- (macOS): running %s\n\t- (Linux): installing it via your distribution's package manager\nSee %s for more information.\n" \
    "fzf" \
    "brew install fzf" \
    "https://github.com/junegunn/fzf#installation")"
)

for d in "${!dependencies[@]}"; do
  if ! command -v "$d" >/dev/null 2>&1; then
    instructions="${dependencies[$d]}"
    printf "command %s is not installed.\n%s" "$d" "$instructions"
    exit 1
  fi
done

function join {
  local d=${1-} f=${2-}
  if shift 2; then
    printf %s "$f" "${@/#/$d}"
  fi
}
export -f join

my_chronic() {
  tmp=$(mktemp) || return # this will be the temp file w/ the output
  "$@" >"$tmp" 2>&1       # this should run the command, respecting all arguments
  ret=$?
  [ "$ret" -eq 0 ] || (echo && cat "$tmp") # if $? (the return of the last run command) is not zero, cat the temp file
  rm -f "$tmp"
}
export -f my_chronic

# check to see if user is logged in
if ! p4 ping &>/dev/null; then
  printf "'%s' command failed. This indicates that you might not be logged into %s@%s.\nTry using %s to generate a session ticket.\nSee %s for more information.\n" \
    "p4 ping" \
    "${P4USER}" \
    "${P4PORT}" \
    "p4 -u ${P4USER} login -a" \
    "https://handbook.sourcegraph.com/departments/ce-support/support/process/p4-enablement/#generate-a-session-ticket"
  exit 1
fi

{
  printf "(re)creating test user '%s' ..." "$P4_TEST_USERNAME"

  # delete test user (if it exists already)
  if p4 users | awk '{print $1}' | grep -Fxq "$P4_TEST_USERNAME"; then
    my_chronic p4 user -yD "$P4_TEST_USERNAME"
  fi

  # create test user
  envsubst <"${SCRIPT_ROOT}/user_template.txt" | my_chronic p4 user -i -f

  printf "done\n"
}

{
  printf "loading protection rules file ..."

  protection_rules_text="$(envsubst <"${SCRIPT_ROOT}/p4_protects.txt")"

  printf "done\n"
}

{
  # parse the protection rules file to discover all the names of the groups
  # for the integration tests
  awk_program=$(
    cat <<-'END'
/# AWK-START/ {flag=1; next}
/# AWK-END/ {flag=0}
{ if (flag && $2 == "group") {print $3} }
END
  )

  mapfile -t all_integration_test_groups < <(awk "$awk_program" <<<"${protection_rules_text}" | sort | uniq)

  # ask the user which groups they'd like the test user to be a member of
  printf "Which group(s) would you like '%s' to be a member of? (tab to select, enter to continue)\n" "$P4_TEST_USERNAME"
  selected_groups="$(fzf --multi --height=20% --layout=reverse <<<"$(join $'\n' "${all_integration_test_groups[@]}")" | sort | uniq)"

  printf "(re)creating test groups (%s) with appropriate memberships (%s) ..." \
    "$(join ', ' "${all_integration_test_groups[@]}")" \
    "${selected_groups/$'\n'/, }"

  # delete any pre-existing test groups from the server
  mapfile -t groups_to_delete < <(comm -12 <(p4 groups | sort) <(printf "%s\n" "${all_integration_test_groups[@]}"))
  for group in "${groups_to_delete[@]}"; do
    my_chronic p4 group -dF "$group"
  done

  # create all the test groups, making sure to add the test user
  # as members of all the groups the user selected
  for group in "${all_integration_test_groups[@]}"; do
    user=""
    if grep -Fxq "$group" <<<"$selected_groups"; then
      user="$P4_TEST_USERNAME"
    fi

    USERNAME="$user" GROUP="$group" envsubst <"${SCRIPT_ROOT}/group_template.txt" | my_chronic p4 group -i
  done

  printf "done\n"
}

{
  printf "uploading protections table ..."

  my_chronic p4 protect -i <<<"${protection_rules_text}"

  printf "done\n"
}
