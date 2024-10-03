package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

type Mode int

const (
	COMMENT Mode = iota
	UNCOMMENT
)

var stdout *bufio.Writer

func getInputInfo(src, commentPrefix []byte) (lines [][]byte, lowestIndentLevel int, mode Mode) {
	lines = bytes.Split(src, []byte{'\n'})
	lines = lines[:len(lines)-1] // art: Split leaves garbage empty []byte at the end
	lowestIndentLevel = 100
	mode = COMMENT

	isModeSet := false

	for _, line := range lines {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		indentLevel := 0
		for i, char := range line {
			// art: we are assuming that indentation is either spaces or tabs, not mix of them
			if char == ' ' || char == '\t' {
				continue
			}

			// art: we are setting mode by checking first line
			if !isModeSet {
				if bytes.HasPrefix(line[i:], commentPrefix) {
					mode = UNCOMMENT
				}
				isModeSet = true
			}

			indentLevel = i
			break
		}

		if indentLevel < lowestIndentLevel {
			lowestIndentLevel = indentLevel
			if lowestIndentLevel == 0 {
				break
			}
		}
	}

	return lines, lowestIndentLevel, mode
}

func handleMultiLine(mode Mode, lines [][]byte, prefix, suffix []byte, lowestIndentLevel int) {
	lastNonEmptyLine := 0

	for i := len(lines) - 1; i >= 0; i-- {
		if len(bytes.TrimSpace(lines[i])) == 0 {
			continue
		}

		lastNonEmptyLine = i
		break
	}

	if mode == COMMENT {
		isPrefixInserted := false

		for i, line := range lines {
			if len(bytes.TrimSpace(line)) == 0 {
				stdout.Write(line)
				stdout.WriteByte('\n')
				continue
			}

			if !isPrefixInserted {
				stdout.Write(line[:lowestIndentLevel])
				stdout.Write(prefix)
				stdout.Write(line[lowestIndentLevel:])
				isPrefixInserted = true
			} else {
				stdout.Write(line)
				if i == lastNonEmptyLine {
					stdout.Write(suffix)
				}
			}

			stdout.WriteByte('\n')
		}
	} else {
		isPrefixRemoved := false

		for i, line := range lines {
			if len(bytes.TrimSpace(line)) == 0 {
				stdout.Write(line)
				stdout.WriteByte('\n')
				continue
			}

			if !isPrefixRemoved {
				left, right, _ := bytes.Cut(line, prefix)
				stdout.Write(left)
				stdout.Write(right)
				isPrefixRemoved = true
			} else {
				if i == lastNonEmptyLine {
					left, right, found := bytes.Cut(line, suffix)
					stdout.Write(left)
					if found {
						stdout.Write(right)
					}
				} else {
					stdout.Write(line)
				}
			}

			stdout.WriteByte('\n')
		}
	}

	stdout.Flush()
}

func handleSingleLine(mode Mode, lines [][]byte, prefix []byte, lowestIndentLevel int) {
	if mode == COMMENT {
		for _, line := range lines {
			if len(bytes.TrimSpace(line)) == 0 {
				stdout.Write(line)
				stdout.WriteByte('\n')
				continue
			}

			stdout.Write(line[:lowestIndentLevel])
			stdout.Write(prefix)
			stdout.Write(line[lowestIndentLevel:])
			stdout.WriteByte('\n')
		}
	} else {
		for _, line := range lines {
			if len(bytes.TrimSpace(line)) == 0 {
				stdout.Write(line)
				stdout.WriteByte('\n')
				continue
			}

			left, right, found := bytes.Cut(line, prefix)
			stdout.Write(left)
			if found {
				stdout.Write(right)
			}
			stdout.WriteByte('\n')
		}
	}

	stdout.Flush()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("comment prefix is required")
	}

	var (
		prefix, suffix []byte
		hasSuffix      bool
	)

	prefix = []byte(os.Args[1])
	if len(os.Args) > 2 {
		suffix = []byte(os.Args[2])
		hasSuffix = true
	}

	src, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	stdout = bufio.NewWriterSize(os.Stdout, len(src)*2)
	lines, lowestIndentLevel, mode := getInputInfo(src, prefix)

	if hasSuffix {
		handleMultiLine(mode, lines, prefix, suffix, lowestIndentLevel)
	} else {
		handleSingleLine(mode, lines, prefix, lowestIndentLevel)
	}
}
