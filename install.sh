#!/bin/bash

function handle_dir() {
    cd $1

    for program_dir in *; do
        install_program $program_dir $1 &
    done

    wait
}

function install_program() {
    local local_dir=~/.local
    local dest=$local_dir/$2

    cd $1

    echo "installing $1"
    if [[ $1 = "backlight-control" ]]; then
        ./build.sh prod
        mv backlight-control backlight-control-notify $dest
    elif [[ -f go.mod ]]; then
        go build
        mv $1 $dest
    elif [[ -f build.sh ]]; then
        ./build.sh prod
        mv $1 $dest
    else
        echo "$2/$1 unhandled program install"
    fi
}

for dir in bin lib; do
    handle_dir $dir &
done

wait
