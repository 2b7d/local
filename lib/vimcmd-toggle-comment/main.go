package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type Comment struct {
	prefix      []byte
	suffix      []byte
	isMultiLine bool
}

type Line struct {
	buf     []byte
	isEmpty bool
}

type Input struct {
	src            []byte
	len            int
	lines          []Line
	lineCount      int
	indentLevel    int
	isUncommenting bool
}

func main() {
	args := os.Args[1:]
	argslen := len(args)

	if argslen < 1 {
		log.Fatalln("need at least 1 argument")
	}

	comment := Comment{
		prefix:      []byte(args[0]),
		isMultiLine: false,
	}

	if argslen > 1 {
		comment.isMultiLine = true
		comment.suffix = []byte(args[1])
	}

	src, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	var (
		input  Input
		output []byte
	)

	parseInput(&input, src, &comment)

	if input.isUncommenting {
		output = make([]byte, 0, input.len)

		if comment.isMultiLine {
			output = uncommentMultiLine(output, &input, &comment)
		} else {
			output = uncommentSingleLine(output, &input, &comment)
		}
	} else {
		output = make([]byte, 0, input.len*3)

		if comment.isMultiLine {
			output = commentOutMultiLine(output, &input, &comment)
		} else {
			output = commentOutSingleLine(output, &input, &comment)
		}
	}

	fmt.Print(string(output))
}

func parseInput(input *Input, src []byte, comment *Comment) {
	input.src = src
	input.len = len(src)
	input.indentLevel = -1
	input.isUncommenting = true

	pad := 0
	line := Line{make([]byte, 0, 80), false}
	isLineStart := true
	isModeSet := false

	for i := 0; i < input.len; i++ {
		if isLineStart {
			for input.src[i] == ' ' || input.src[i] == '\t' {
				line.buf = append(line.buf, input.src[i])
				pad++
				i++
			}

			if input.src[i] == '\n' {
				line.isEmpty = true
				pad--
			}

			if !isModeSet && !line.isEmpty {
				if comment.isMultiLine {
					input.isUncommenting = bytes.HasPrefix(input.src[i:], comment.prefix)
					isModeSet = true
				} else {
					// if at least one line is not a comment then we are commenting out
					if !bytes.HasPrefix(input.src[i:], comment.prefix) {
						input.isUncommenting = false
						isModeSet = true
					}
				}
			}

			isLineStart = false
		}

		line.buf = append(line.buf, input.src[i])

		if input.src[i] == '\n' {
			if input.indentLevel == -1 {
				input.indentLevel = pad
			} else {
				if !comment.isMultiLine {
					if pad < input.indentLevel && !line.isEmpty {
						input.indentLevel = pad
					}
				}
			}

			input.lines = append(input.lines, line)

			pad = 0
			line = Line{make([]byte, 0, 80), false}
			isLineStart = true
		}
	}

	input.lineCount = len(input.lines)
}

func commentOutMultiLine(output []byte, input *Input, comment *Comment) []byte {
	isPrefixInserted := false
	lastLine := findLastNonEmptyLine(input.lines)

	for lineNum, line := range input.lines {
		if line.isEmpty {
			output = append(output, line.buf...)
			continue
		}

		for i, ch := range line.buf {
			if !isPrefixInserted && i == input.indentLevel {
				output = append(output, comment.prefix...)
				isPrefixInserted = true
			}

			if lineNum == lastLine && ch == '\n' {
				output = append(output, comment.suffix...)
			}

			output = append(output, ch)
		}
	}

	return output
}

func uncommentMultiLine(output []byte, input *Input, comment *Comment) []byte {
	isPrefixDeleted := false
	lastLine := findLastNonEmptyLine(input.lines)

	for lineNum, line := range input.lines {
		if line.isEmpty {
			output = append(output, line.buf...)
			continue
		}

		for i := 0; i < len(line.buf); i++ {
			if !isPrefixDeleted && i == input.indentLevel {
				i += len(comment.prefix)
				isPrefixDeleted = true
			}

			if lineNum == lastLine && line.buf[i] == '\n' {
				end := len(output) - len(comment.suffix)
				output = output[:end]
			}

			output = append(output, line.buf[i])
		}
	}

	return output
}

func findLastNonEmptyLine(lines []Line) int {
	for i := len(lines) - 1; i >= 0; i-- {
		if !lines[i].isEmpty {
			return i
		}
	}

	return 0
}

func commentOutSingleLine(output []byte, input *Input, comment *Comment) []byte {
	for _, line := range input.lines {
		if line.isEmpty {
			output = append(output, line.buf...)
			continue
		}

		for i, ch := range line.buf {
			if i == input.indentLevel {
				output = append(output, comment.prefix...)
			}
			output = append(output, ch)
		}
	}

	return output
}

func uncommentSingleLine(output []byte, input *Input, comment *Comment) []byte {
	for _, line := range input.lines {
		if line.isEmpty {
			output = append(output, line.buf...)
			continue
		}

		for i := 0; i < len(line.buf); i++ {
			if i == input.indentLevel {
				i += len(comment.prefix)
			}
			output = append(output, line.buf[i])
		}
	}

	return output
}
