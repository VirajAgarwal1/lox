package parser_demos

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/parser"
)

func ParserDemo() {
	const scanner_buf_cap uint32 = 2

	// Sample input:  42, "hello", true
	sample_input := bufio.NewReader(strings.NewReader("42,\"hello\",true,identifier,false,2.89"))
	buf_scanner := lexer.BufferedLexer{}
	buf_scanner.Initialize(sample_input, uint32(scanner_buf_cap))

	// Run the parser
	_, _, err := parser.Parse_expression(&buf_scanner)
	if err != nil {
		fmt.Println("Parser error:", err)
		fmt.Println()
		fmt.Println()
	}
}

// var a = 1 + (2 - 3) * 4
// a = a + 1
// fun(a + 2)
