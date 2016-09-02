package main

import (
	"fmt"
	"os"

	"github.com/GreenDelta/matlib"
)

func main() {

	if len(os.Args) < 2 {
		help()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "-h", "help":
		help()
	case "-i", "invert":
		invert()
	default:
		fmt.Println("Unknown command:", cmd, " (try help)")
	}

}

func help() {
	text := `
matlib

Usage: matlib <command> <args>

-h, help                     prints this help
-i, invert <input> <output>  inverts the matrix in the input file and writes it
                             to the output file
`
	fmt.Println(text)
}

func invert() {
	if len(os.Args) < 4 {
		fmt.Println("Not enough arguments: invert <input> <output>")
		return
	}
	m, err := matlib.Load(os.Args[2])
	if err != nil {
		fmt.Println("Failed to read matrix from", os.Args[2], err.Error())
		return
	}
	err = m.InvertInPlace()
	if err != nil {
		fmt.Println("Failed to invert matrix", os.Args[2], err.Error())
		return
	}
	err = matlib.Save(m, os.Args[3])
	if err != nil {
		fmt.Println("Failed to write inverse to", os.Args[3], err.Error())
	}
}
