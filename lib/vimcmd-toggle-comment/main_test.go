package main

import (
	"testing"
)

type tcase struct {
	id int
	input string
	want string
}

func TestSingleLineCommentOut(t *testing.T) {
	comment := Comment{
		prefix: []byte("//"),
		isMultiLine: false,
	}

	tests := []tcase{
		{1, "A\n",              "//A\n"},
		{2, "  A\n",            "  //A\n"},
		{3, "\n",               "\n"},
		{4, "  \n",             "  \n"},
		{5, "A\n  B\n",         "//A\n//  B\n"},
		{6, "  A\n  B\n",       "  //A\n  //B\n"},
		{7, "  A\nB\n",         "//  A\n//B\n"},
		{8, "  A\n\n  \n  B\n", "  //A\n\n  \n  //B\n"},
		{9, "\tA\n",            "\t//A\n"},
		{10, "  //A\n",         "  ////A\n"},
		{11, "//\n",            "////\n"},
		{12, "//A\n//B\nC\n",   "////A\n////B\n//C\n"},
	}

	for _, tc := range tests {
		var (
			input Input
			output []byte
		)
		parseInput(&input, []byte(tc.input), &comment)

		output = commentOutSingleLine(output, &input, &comment)
		if string(output) != tc.want {
			t.Errorf("\ncase: %d\ninput: %q\nwant:  %q\ngot:   %q\n", tc.id, tc.input, tc.want, output)
		}
	}
}

func TestSingleLineUncomment(t *testing.T) {
	comment := Comment{
		prefix: []byte("//"),
		isMultiLine: false,
	}

	tests := []tcase{
		{1, "//A\n",                "A\n"},
		{2, "  //A\n",              "  A\n"},
		{3, "\n",                   "\n"},
		{4, "  \n",                 "  \n"},
		{5, "//A\n//  B\n",         "A\n  B\n"},
		{6, "  //A\n  //B\n",       "  A\n  B\n"},
		{7, "//  A\n//B\n",         "  A\nB\n"},
		{8, "  //A\n\n  \n  //B\n", "  A\n\n  \n  B\n"},
		{9, "\t//A\n",              "\tA\n"},
		{10, "  ////A\n",           "  //A\n"},
		{11, "////\n",              "//\n"},
		{12, "////A\n////B\n//C\n", "//A\n//B\nC\n"},
	}

	for _, tc := range tests {
		var (
			input Input
			output []byte
		)
		parseInput(&input, []byte(tc.input), &comment)

		output = uncommentSingleLine(output, &input, &comment)
		if string(output) != tc.want {
			t.Errorf("\ncase: %d\ninput: %q\nwant:  %q\ngot:   %q\n", tc.id, tc.input, tc.want, output)
		}
	}
}

func TestMultiLineCommentOut(t *testing.T) {
	comment := Comment{
		prefix: []byte("/*"),
		suffix: []byte("*/"),
		isMultiLine: true,
	}

	tests := []tcase{
		{1, "A\n",              "/*A*/\n"},
		{2, "  A\n",            "  /*A*/\n"},
		{3, "\n",               "\n"},
		{4, "  \n",             "  \n"},
		{5, "A\n  B\n",         "/*A\n  B*/\n"},
		{6, "  A\n  B\n",       "  /*A\n  B*/\n"},
		{7, "  A\nB\n",         "  /*A\nB*/\n"},
		{8, "  A\n\n  \n  B\n", "  /*A\n\n  \n  B*/\n"},
		{9, "\tA\n",            "\t/*A*/\n"},
		{10, "\n  \n\n",        "\n  \n\n"},
	}

	for _, tc := range tests {
		var (
			input Input
			output []byte
		)
		parseInput(&input, []byte(tc.input), &comment)

		output = commentOutMultiLine(output, &input, &comment)
		if string(output) != tc.want {
			t.Errorf("\ncase: %d\ninput: %q\nwant:  %q\ngot:   %q\n", tc.id, tc.input, tc.want, output)
		}
	}
}

func TestMultiLineUncomment(t *testing.T) {
	comment := Comment{
		prefix: []byte("/*"),
		suffix: []byte("*/"),
		isMultiLine: true,
	}

	tests := []tcase{
		{1, "/*A*/\n",              "A\n"},
		{2, "  /*A*/\n",            "  A\n"},
		{3, "\n",                   "\n"},
		{4, "  \n",                 "  \n"},
		{5, "/*A\n  B*/\n",         "A\n  B\n"},
		{6, "  /*A\n  B*/\n",       "  A\n  B\n"},
		{7, "  /*A\nB*/\n",         "  A\nB\n"},
		{8, "  /*A\n\n  \n  B*/\n", "  A\n\n  \n  B\n"},
		{9, "\t/*A*/\n",            "\tA\n"},
		{10, "\n  \n\n",            "\n  \n\n"},
		{11, "/**/\n",              "\n"},
		{12, "/*\n  \n*/\n",        "\n  \n\n"},
	}

	for _, tc := range tests {
		var (
			input Input
			output []byte
		)
		parseInput(&input, []byte(tc.input), &comment)

		output = uncommentMultiLine(output, &input, &comment)
		if string(output) != tc.want {
			t.Errorf("\ncase: %d\ninput: %q\nwant:  %q\ngot:   %q\n", tc.id, tc.input, tc.want, output)
		}
	}
}
