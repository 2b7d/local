package main

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	wallpapersDir := homedir + "/.local/share/wallpapers"

	wallpapers, err := os.ReadDir(wallpapersDir)
	if err != nil {
		log.Fatal(err)
	}

	dirinfo, err := os.Stat(wallpapersDir)
	if err != nil {
		log.Fatal(err)
	}

	savedWallpapers := wallpapers
	savedModTime := dirinfo.ModTime()

	for {
		dirinfo, err := os.Stat(wallpapersDir)
		if err != nil {
			log.Fatal(err)
		}

		if dirinfo.ModTime() != savedModTime {
			var err error

			wallpapers, err = os.ReadDir(wallpapersDir)
			if err != nil {
				log.Fatal(err)
			}

			savedWallpapers = wallpapers
			savedModTime = dirinfo.ModTime()
		}

		index := rand.Intn(len(wallpapers))

		cmd := exec.Command("nitrogen", "--set-zoom", wallpapersDir+"/"+wallpapers[index].Name())
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		if len(wallpapers) == 1 {
			wallpapers = savedWallpapers
		} else {
			tmp := make([]os.DirEntry, 0, len(wallpapers))
			for i, w := range wallpapers {
				if i != index {
					tmp = append(tmp, w)
				}
			}
			wallpapers = tmp
		}

		time.Sleep(5 * time.Minute)
	}
}
