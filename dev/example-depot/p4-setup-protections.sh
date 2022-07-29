#!/usr/bin/env bash

SCRIPT_ROOT="$(dirname "${BASH_SOURCE[0]}")"
cd "${SCRIPT_ROOT}"

SCRATCH=$(mktemp -d -t p4_setup_protections_XXXXXXX)
cleanup() {
  rm -rf "$SCRATCH"
}
trap cleanup EXIT

set -euxo pipefail

export P4USER="${P4USER:-"admin"}"                   # the name of the Perforce superuser that the script will use to create the depot
export P4PORT="${P4PORT:-"perforce.sgdev.org:1666"}" # the address of the Perforce server to connect to

export P4_TEST_USERNAME="${P4_TEST_USERNAME:-"test-perforce"}"           # the name of the client that the script will use while it creates the depot
export P4_TEST_EMAIL="${P4_TEST_EMAIL:-"test-perforce@sourcegraph.com"}" # the name of the client that the script will use while it creates the depot

export DEPOT_NAME="${DEPOT_NAME:-"integration-test-depot"}" # the name of the depot that the script will create on the server

# check to see if user is logged in
if ! p4 ping >/dev/null; then
  printf "'%s' command failed. This indicates that you might not be logged into %s@%s.\nTry using %s to generate a session ticket. See %s for more information.\n" \
    "p4 ping" \
    "${P4USER}" \
    "${P4PORT}" \
    "p4 -u ${P4USER} login -a" \
    "https://handbook.sourcegraph.com/departments/ce-support/support/process/p4-enablement/#generate-a-session-ticket"
  exit 1
fi

# delete test user (if it exists already)
if p4 users | awk '{print $1}' | grep -q "${P4_TEST_USERNAME}"; then
  p4 user -yD "${P4_TEST_USERNAME}"
fi

# create test user
envsubst <"${SCRIPT_ROOT}/user_template.txt" | p4 user -i -f

PROTECTION_RULES="$(envsubst <"${SCRIPT_ROOT}/p4_protects.txt")"
mapfile -t INTEGRATION_TEST_GROUPS < <(awk '/# AWK-START/{flag=1; next} /# AWK-END/{flag=0} { if (flag && $2 == "group") {print $3} }' <<<"${PROTECTION_RULES}" | sort)

mapfile -t groups_to_delete < <(comm -12 <(p4 groups | sort) <(printf "%s\n" "${INTEGRATION_TEST_GROUPS[@]}"))
for group in "${groups_to_delete[@]}"; do
  p4 group -dF "${group}"
done

for group in "${INTEGRATION_TEST_GROUPS[@]}"; do
  GROUP="$group" envsubst <"${SCRIPT_ROOT}/group_template.txt" | p4 group -i
done

p4 protect -i <<<"${PROTECTION_RULES}"
