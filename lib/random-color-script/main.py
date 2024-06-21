#!/bin/python3

import subprocess
import random
import os
import sys

dir = os.path.expanduser("~/.local/share/color-scripts")
try:
    scripts = os.listdir(dir)
except Exception as err:
    print(err)
    sys.exit(1)

if not scripts:
    print("scripts folder is empty")
    sys.exit(0)

script = random.choice(scripts)
execpath = os.path.join(dir, script)
try:
    subprocess.run(execpath)
except Exception as err:
    print(err)
    sys.exit(1)
