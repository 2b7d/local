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
		log.Fatal("not enough arguments")
	}

	home := os.Getenv("HOME")
	if home == "" {
		log.Fatal("HOME is not set")
	}

	errs := []error{}
	success := []string{}
	args := os.Args[1:]

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
