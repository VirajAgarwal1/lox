package streamable_parser_demos

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/VirajAgarwal1/lox/lexer"
	ebnf_to_bnf "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/ebnf_to_bnf"
	grammar_file_parser "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/grammar_file_parser"
)

func Sample_EBNF_to_BNF() {
	file_reader, err := os.Open("/Users/viraj.agarwal/Projects/lox/parser/lox.grammar")
	if err != nil {
		panic(err)
	}
	source := bufio.NewReader(file_reader)

	scanner := lexer.LexicalAnalyzer{}
	scanner.Initialize(source)
	ebnf_grammar, err := grammar_file_parser.ProcessGrammarDefinition(&scanner)
	if err != nil && err != io.EOF {
		panic(err)
	}

	bnf_grammar := ebnf_to_bnf.EbnfToBnfConverter(ebnf_grammar)
	fmt.Println(bnf_grammar)
}
