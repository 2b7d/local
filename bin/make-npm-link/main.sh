#!/bin/bash

set -e

if [[ $# -eq 0 ]]; then
    echo "provide executable name[s]"
    exit 1
fi

error=0

for f in $@; do
    f=$HOME/.local/lib/npm/bin/$f

    if [[ ! -f $f ]]; then
        echo "$f does not exist"
        error=1
        continue
    fi

    if [[ ! -x $f ]]; then
        echo "$f is not executable or user \`$USER\` does not have permission"
        error=1
        continue
    fi
done

if [[ $error -eq 1 ]]; then
    exit 1
fi

for f in $@; do
    target=$HOME/.local/lib/npm/bin/$f
    dest=$HOME/.local/bin/$f

    if [[ -e $dest ]]; then
        rm $dest
    fi

    ln -sr $target $dest
done
