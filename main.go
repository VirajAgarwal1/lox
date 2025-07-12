/*
^ NOTE:
	Currently, sourcecode of Lox only supports reading it in **UTF-8** encoding. If provided sourcecode in other encodings then the program may behave in unexpected ways..
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	errorhandler "github.com/VirajAgarwal1/lox/errorhandler"
)

const READER_BUFFER_SIZE = 4096

type Lox struct {
}

func (program *Lox) Run(source *bufio.Reader) error {

	for {
		r, _, err := source.ReadRune()
		if err != nil && err != io.EOF {
			return errorhandler.RetErr("", err)
		}
		// TODO: Got each individual char, now to get tokens from them
		fmt.Printf("%s", string(r))
		if err == io.EOF {
			break
		}
	}
	fmt.Println()
	// TODO: Get the tokens and do something with those tokens
	return nil
}
func (program *Lox) RunFile(filePath string) error {
	fr, err := os.Open(filePath)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	buf_reader := bufio.NewReaderSize(fr, READER_BUFFER_SIZE)
	err = program.Run(buf_reader)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	err = fr.Close()
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	return nil
}

// TODO: This Repl will require some improvements since, it cannot detect when an instruction is not complete
func (program *Lox) Repl() error {
	input := os.Stdin
	buf_input := bufio.NewReaderSize(input, READER_BUFFER_SIZE)
	for {
		err := program.Run(buf_input)
		if err != nil && err != io.EOF {
			return errorhandler.RetErr("", err)
		}
		if err == io.EOF {
			break
		}
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
			err = errorhandler.RetErr("", err)
			errorhandler.ReportErr(err)
		}
	} else {
		// Run the source file
		err := program.RunFile(args[0])
		if err != nil {
			errorhandler.ReportErr(err)
		}
	}
}
