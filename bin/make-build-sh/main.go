package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
)

const src = `#!/bin/bash

set -xe

files=main.c
outname=main

flags="-g -Werror=declaration-after-statement -Wall -Wextra -pedantic -std=c99"
incl=
libs=

#if [[ -z $1 ]]; then
#    echo "provide build option"
#    exit 1
#fi

#case $1 in
#    build_option)
#        files=
#        outname=
#        incl=
#        libs=
#        ;;
#    *)
#        echo "unknown build option $1"
#        exit 1
#        ;;
#esac

#if [[ $2 = "prod" ]]; then
if [[ $1 = "prod" ]]; then
    flags=${flags/-g/-O2}
fi

gcc $flags -o $outname $files $incl $libs`

func main() {
	file, err := os.OpenFile("build.sh", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0700)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			fmt.Println("build.sh already exist")
			os.Exit(0)
		}
		log.Fatal(err)
	}

	file.WriteString(src)
	file.Close()
}
