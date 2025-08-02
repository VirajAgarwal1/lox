package parser_demos

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/parser"
)

func ParserDemo() {
	// Sample input:  42, "hello", true
	sample_input := bufio.NewReader(strings.NewReader("42,\"hello\",true,identifier,false,2.89"))
	buf_scanner := lexer.BufferedLexicalAnalyzer{}
	buf_scanner.Initialize(sample_input)

	// Run the parser
	_, _, err := parser.Parse_expression(&buf_scanner)
	if err != nil {
		fmt.Println("Parser error:", err)
		fmt.Println()
		fmt.Println()
	}
}
