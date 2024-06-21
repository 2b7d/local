#!/bin/bash

path=~/.local/share
dir=$path/todo
file=$dir/list

if [[ ! -d $dir ]]; then
    mkdir $dir
    touch $file
elif [[ ! -e $file ]]; then
    touch $file
fi

if [[ $1 = "print" ]]; then
    echo -e "━━━ Todo List\n"
    cat $file
    echo -e "\n━━━\n"
    exit
fi

nvim $file
