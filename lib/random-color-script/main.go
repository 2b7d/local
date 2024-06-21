package main

import (
	"log"
	"math/rand/v2"
	"os"
	"path"
	"syscall"
)

func main() {
	home := os.Getenv("HOME")
	if home == "" {
		log.Fatal("$HOME environment variable is not set")
	}

	dirpath := path.Join(home, ".local/share/color-scripts")

	dirents, err := os.ReadDir(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	if len(dirents) == 0 {
		log.Fatal("color-scripts directory is empty")
	}

	script := dirents[rand.IntN(len(dirents))]
	execpath := path.Join(dirpath, script.Name())

	if err := syscall.Exec(execpath, os.Args, os.Environ()); err != nil {
		log.Fatal(err)
	}
}
