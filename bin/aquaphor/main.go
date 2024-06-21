package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"strconv"
)

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Your system does not have $HOME environment variable set")
	}

	statedir := path.Join(homedir, ".local", "share", "aquaphor")
	statefile := path.Join(statedir, "state")

	if err := os.Mkdir(statedir, 0700); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			log.Fatal(err)
		}
	}

	state, err := os.OpenFile(statefile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer state.Close()

	var val int
	fmt.Fscanf(state, "%v", &val)

	if len(os.Args) > 1 {
		if os.Args[1] == "info" {
			fmt.Println(val, "/ 350")
			fmt.Println("left:", 350-val)
			return
		} else {
			v, err := strconv.Atoi(os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			val += v
		}
	} else {
		val += 2
	}

	if _, err := state.Seek(0, io.SeekStart); err != nil {
		log.Fatal(err)
	}

	if err := state.Truncate(0); err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(state, val)
}
