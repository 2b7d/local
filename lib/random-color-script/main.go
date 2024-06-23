package main

import (
	"log"
	"math/rand/v2"
	"os"
	"path"
	"syscall"
)

const scriptsDirName = "color-scripts"

func main() {
	home := os.Getenv("HOME")
	if home == "" {
		log.Fatal("$HOME is not set")
	}

	dirpath := path.Join(home, ".local/share", scriptsDirName)
	dirents, err := os.ReadDir(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	if len(dirents) == 0 {
		log.Fatalf("%s directory is empty", scriptsDirName)
	}

	script := dirents[rand.IntN(len(dirents))]
	execpath := path.Join(dirpath, script.Name())

	args := []string{execpath}
	if err := syscall.Exec(args[0], args, os.Environ()); err != nil {
		log.Fatal(err)
	}
}
