#!/bin/bash

set -e

cat << "EOF" > build.sh
#!/bin/bash

set -xe

files=main.c
outname=main

flags="-g -Werror=declaration-after-statement -Wall -Wextra -pedantic -std=c99"
incl=
libs=

if [[ $1 = "prod" ]]; then
    flags=${flags/-g/-O2}
fi

gcc $flags -o $outname $files $incl $libs
EOF

chmod u+x build.sh
