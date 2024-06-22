package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "not enough arguments\n")
		help(os.Stderr)
		os.Exit(1)
	}

	args := os.Args[1:]

	for _, arg := range args {
		if arg == "--help" {
			help(os.Stdout)
			os.Exit(0)
		}
	}

	home := os.Getenv("HOME")
	if home == "" {
		log.Fatal("HOME is not set")
	}

	errs := []error{}
	success := []string{}

	for _, arg := range args {
		parts := strings.Split(arg, "/")

		index, found := findLocalLib(parts)

		if !found && path.IsAbs(arg) {
			errs = append(errs, fmt.Errorf("%s is not .local/lib path", arg))
			continue
		}

		if found {
			parts = parts[index+1:]
		}

		execpath := path.Join(parts...)
		target := path.Join(home, ".local", "lib", execpath)

		info, err := os.Stat(target)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				errs = append(errs, fmt.Errorf("%s does not exist in .local/lib", arg))
				continue
			}
			log.Fatal(err)
		}

		if info.IsDir() {
			errs = append(errs, fmt.Errorf("%s is directory", arg))
			continue
		}

		if info.Mode()&0111 == 0 {
			errs = append(errs, fmt.Errorf("%s is not executable", arg))
			continue
		}

		newname := path.Join(home, ".local/bin", path.Base(execpath))
		oldname := path.Join("../lib/", execpath)

		if err := os.Symlink(oldname, newname); err != nil {
			log.Fatal(err)
		}

		success = append(success, arg)
	}

	if len(success) > 0 {
		fmt.Println("Success:")
		for _, s := range success {
			fmt.Printf("    %s\n", s)
		}
	}

	if len(errs) > 0 {
		fmt.Println("Errors:")
		for _, e := range errs {
			fmt.Printf("    %s\n", e)
		}
	}
}

func help(stream *os.File) {
	stream.WriteString(fmt.Sprintf(`Usage: %s <TARGETS>
TARGETS:
    list of executable files in lib directory to make symbolic links of
`, os.Args[0]))
}

func findLocalLib(pathParts []string) (int, bool) {
	localIndex := -1

	for i, p := range pathParts {
		if p == ".local" {
			localIndex = i
			break
		}
	}

	if localIndex == -1 {
		return -1, false
	}

	libIndex := localIndex + 1

	if libIndex >= len(pathParts) {
		return -1, false
	}

	if pathParts[libIndex] != "lib" {
		return -1, false
	}

	return libIndex, true
}
