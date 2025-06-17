package main

import (
	"fmt"
	"os"
)

type Token struct {
}

type Lox struct {
	source *string
	tokens []Token
}

func main() {
	args := os.Args[1:]
	fmt.Println(args)
}
