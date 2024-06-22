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

	for _, filePath := range args {
		parts := strings.Split(filePath, "/")
		index, found := findLocalLib(parts)

		if !found && path.IsAbs(filePath) {
			errs = append(errs, fmt.Errorf("%s is not .local/lib path", filePath))
			continue
		}

		if found {
			parts = parts[index+1:]
		}

		relFilePath := path.Join(parts...)
		absFilePath := path.Join(home, ".local/lib", relFilePath)

		info, err := os.Stat(absFilePath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				errs = append(errs, fmt.Errorf("%s does not exist in .local/lib", filePath))
				continue
			}
			log.Fatal(err)
		}

		if info.IsDir() {
			errs = append(errs, fmt.Errorf("%s is directory", filePath))
			continue
		}

		if info.Mode()&0111 == 0 {
			errs = append(errs, fmt.Errorf("%s is not executable", filePath))
			continue
		}

		linkName := path.Join(home, ".local/bin", path.Base(relFilePath))
		target := path.Join("../lib", relFilePath)

		if err := os.Symlink(target, linkName); err != nil {
			if errors.Is(err, fs.ErrExist) {
				errs = append(errs, fmt.Errorf("%s already exist in .local/bin", filePath))
				continue
			}
			log.Fatal(err)
		}

		success = append(success, filePath)
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
