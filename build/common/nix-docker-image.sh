#!/usr/bin/env bash

set -e

if [[ $NIX_BUILD_NATIVE -eq 1 ]]; then
    case $(uname -m) in
    "arm64")
        SYSTEM="aarch64-linux"
        ;;
    *)
        SYSTEM="x86_64-linux"
        ;;
    esac

    nix build .\#packages.${SYSTEM}.docker-image --print-build-logs
    exit $?
fi

which nix > /dev/null

if [[ ( $? -eq 0 ) && ( `uname` == "Linux" ) ]]; then
    nix build .\#docker-image --print-build-logs
    exit $?
fi

./build/common/nix.sh build .\#docker-image --print-build-logs
