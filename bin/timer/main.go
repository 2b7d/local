package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	duration := time.Minute * 45

	if len(os.Args) > 1 {
		var err error

		if duration, err = time.ParseDuration(os.Args[1]); err != nil {
			fmt.Println("timer duration is invalid, must be 3h4m5s")
			os.Exit(1)
		}
	}

	startingDurationText := duration.String()

	for {
		fmt.Print("\033[2J\033[H")

		h := duration / time.Hour
		m := duration % time.Hour / time.Minute
		s := duration % time.Hour % time.Minute / time.Second

		fmt.Printf("  %02d:%02d:%02d\033[H", h, m, s)

		if duration <= 0 {
			cmd := exec.Command("notify-send")
			cmd.Args = append(cmd.Args, "Time is Up", fmt.Sprintf("Timer %s has finished", startingDurationText))

			if err := cmd.Run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println()
			break
		}

		duration -= time.Second
		time.Sleep(time.Second)
	}
}
