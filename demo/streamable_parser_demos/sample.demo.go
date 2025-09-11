package streamable_parser_demos

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/VirajAgarwal1/lox/lexer"
	parser_generator "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator"
)

func Sample_fun() {
	file_reader, err := os.Open("/Users/viraj.agarwal/Projects/lox/parser/lox.grammar")
	if err != nil {
		panic(err)
	}
	source := bufio.NewReader(file_reader)

	scanner := lexer.LexicalAnalyzer{}
	scanner.Initialize(source)
	prod, err := parser_generator.ProcessGrammarDefinition(&scanner)
	if err != nil && err != io.EOF {
		panic(err)
	}

	sanitized_grammar := parser_generator.EbnfToBnfConverter(prod)
	fmt.Println(sanitized_grammar)
}
