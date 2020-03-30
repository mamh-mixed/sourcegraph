#!/usr/bin/env bash

## This script will ensure that the libsqlite3-pcre dynamic library exists in the root of this
## repository (either libsqlite3-pcre.dylib for Darwin or libsqlite3-pcre.so for linux). This
## script is used by run the symbol service locally, which compiles against the shared library.
##
## Invocation:
## - `./libsqlite3-pcre/build.sh`         : build the library
## - `./libsqlite3-pcre/build.sh libpath` : output its path

cd "$(dirname "${BASH_SOURCE[0]}")/../.."
set -ux

OUTPUT=`mktemp -d -t sgdockerbuild_XXXXXXX`
cleanup() {
    rm -rf "$OUTPUT"
}
trap cleanup EXIT

function libpath() {
    case "$OSTYPE" in
        darwin*)
            echo "$PWD/libsqlite3-pcre.dylib"
            ;;

        linux*)
            echo "$PWD/libsqlite3-pcre.so"
            ;;

        *)
            echo "Unknown platform $OSTYPE"
            exit 1
            ;;
    esac
}

function build() {
    libsqlite3PcrePath=$(libpath)
    if [ -f "$libsqlite3PcrePath" ]; then
        # Already exists
        exit 0
    fi

    if ! command -v pkg-config >/dev/null 2>&1 || ! command -v pkg-config --cflags sqlite3 libpcre >/dev/null 2>&1; then
        echo "Missing sqlite dependencies."
        case "$OSTYPE" in
            darwin*)
                echo "Install them by running 'brew install pkg-config sqlite pcre FiloSottile/musl-cross/musl-cross'"
                ;;

            linux*)
                echo "Install them by running 'apt-get install libpcre3-dev libsqlite3-dev pkg-config musl-tools'"
                ;;

            *)
                echo "See the local development documentation: https://github.com/sourcegraph/sourcegraph/blob/master/doc/dev/local_development.md#step-2-install-dependencies"
                ;;
        esac

        exit 1
    fi

    echo "--- $libsqlite3PcrePath build"
    curl -fsSL https://codeload.github.com/ralight/sqlite3-pcre/tar.gz/c98da412b431edb4db22d3245c99e6c198d49f7a | tar -C "$OUTPUT" -xzvf - --strip 1
    cd "$OUTPUT"

    case "$OSTYPE" in
        darwin*)
            # pkg-config spits out multiple arguments and must not be quoted.
            gcc -fno-common -dynamiclib pcre.c -o "$libsqlite3PcrePath" $(pkg-config --cflags sqlite3 libpcre) $(pkg-config --libs libpcre) -fPIC
            exit 0
            ;;

        linux*)
            # pkg-config spits out multiple arguments and must not be quoted.
            gcc -shared -o "$libsqlite3PcrePath" $(pkg-config --cflags sqlite3 libpcre) -fPIC -W -Werror pcre.c $(pkg-config --libs libpcre) -Wl,-z,defs
            exit 0
            ;;

        *)
            echo "See the local development documentation: https://github.com/sourcegraph/sourcegraph/blob/master/doc/dev/local_development.md#step-2-install-dependencies"Uecho "Unknown platform $OSTYPE"
            exit 1
            ;;
    esac
}

"${1:-build}"
