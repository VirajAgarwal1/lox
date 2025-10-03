package streamable_parser_demos

import (
	"bufio"
	"io"
	"os"

	"github.com/VirajAgarwal1/lox/lexer"
	ebnf_to_bnf "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/ebnf_to_bnf"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/first_follow"
	grammar_file_parser "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/grammar_file_parser"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/parser_writer"
)

func Sample_parser_writer() {
	current_dir, _ := os.Getwd()

	file_reader, err := os.Open("parser/lox.grammar")
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

	firsts := first_follow.ComputeFirstSets(bnf_grammar)

	follow := first_follow.ComputeFollowSets(bnf_grammar)

	parser_writer.WriteParser(current_dir+"/streamable_parser/generated_parser.go", bnf_grammar, "expression", firsts, follow)
}
