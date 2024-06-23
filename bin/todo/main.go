package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"syscall"
)

const dirname = "todo"
const filename = "list"
const textEditor = "/bin/nvim"

func main() {
	home := os.Getenv("HOME")
	if home == "" {
		log.Fatal("$HOME is not set")
	}

	dirpath := path.Join(home, ".local/share", dirname)
	filepath := path.Join(dirpath, filename)

	if err := os.Mkdir(dirpath, 0700); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			log.Fatal(err)
		}
	}

	f, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if len(os.Args) >= 2 && os.Args[1] == "list" {
		fmt.Printf("━━━ Todo List\n\n")
		if _, err := io.Copy(os.Stdout, f); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\n━━━\n\n")
		os.Exit(0)
	}

	args := []string{textEditor, filepath}
	if err := syscall.Exec(args[0], args, os.Environ()); err != nil {
		log.Fatal(err)
	}
}
