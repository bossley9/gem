package main

import (
	"fmt"
	"os"

	"git.sr.ht/~bossley9/gem"
)

func main() {
	args := os.Args

	switch len(args) {
	case 3:
		err := convertFile(args[1], args[2])
		if err != nil {
			fmt.Println(err)
		}
	default:
		usage()
	}
}

func usage() {
	fmt.Println("gem converts files from Gemtext into HTML.\n")
	fmt.Println("Usage: gem [input gemtext file] [output HTML file]\n")
}

func convertFile(input string, output string) error {
	gemtext, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	html := gem.ToHTML(string(gemtext))

	if err := os.WriteFile(output, []byte(html), 0600); err != nil {
		return err
	}

	return nil
}
