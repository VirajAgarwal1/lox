package main

import (
	"bufio"
	"fmt"
	"os"

	error_pak "github.com/VirajAgarwal1/lox/error_pak"
)

type Token struct {
}

type Lox struct {
	source []byte
	tokens []Token
}

func (this *Lox) RunFile(filePath string) error {

	sourceText, err := os.ReadFile(filePath)
	if err != nil {
		return error_pak.RetErr(err.Error(), nil)
	}
	this.source = sourceText

	// TODO: Add the run command here
	fmt.Println(string(this.source))
	return nil
}

func (this *Lox) Repl() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("[Lox]> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		// TODO: Add the run command here
		fmt.Printf("%v\n", line)
	}
	if err := scanner.Err(); err != nil {
		return error_pak.RetErr("", err)
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		// Wrong usage of cli
		fmt.Println("Usage: jlox [script]")
		return
	}

	program := Lox{}

	if len(args) == 0 {
		// Run in REPL mode (Read, Execute, Print, Loop)
		err := program.Repl()
		if err != nil {
			fmt.Println(error_pak.RetErr("", err))
		}
	} else {
		// Run the source file
		err := program.RunFile(args[0])
		if err != nil {
			fmt.Println(error_pak.RetErr("", err))
		}
	}
}
