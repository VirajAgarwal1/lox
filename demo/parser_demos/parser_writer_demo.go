package parser_demos

import (
	"bufio"
	"io"
	"os"

	"github.com/VirajAgarwal1/lox/errorhandler"
	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/parser/grammar"
)

func WriteGrammarParserDemo() {
	file_reader, err := os.Open("parser/lox.grammar")
	if err != nil {
		errorhandler.ReportErr(err)
		return
	}
	buf_file_reader := bufio.NewReader(file_reader)

	scanner := lexer.LexicalAnalyzer{}
	scanner.Initialize(buf_file_reader)

	processed_grammar, err := grammar.ProcessGrammarDefinition(&scanner)
	if err != nil && err != io.EOF {
		errorhandler.ReportErr(err)
		return
	}

	err = grammar.GenerateGrammarParserFile(processed_grammar, "parser/generated_parser.go")
	if err != nil {
		errorhandler.ReportErr(err)
		return
	}
}
