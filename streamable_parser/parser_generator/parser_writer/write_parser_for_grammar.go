package parser_writer

import (
	"bufio"
	"os"

	"github.com/VirajAgarwal1/lox/lexer/dfa"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/first_follow"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/parser_writer/code_snippets"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/utils"
)

func WriteParser(path string, bnf_grammar map[string][][]utils.Grammar_element, starting_non_terminal string, firstSet map[string]first_follow.FirstSetInfo, followSet map[string][]dfa.TokenType) error {
	// Always overwrite the file cleanly
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	code := ""
	code += code_snippets.Package_and_Imports + "\n\n"
	code += code_snippets.Consts_code(starting_non_terminal) + "\n\n"
	code += code_snippets.GrammarRules_code(bnf_grammar, firstSet, followSet) + "\n"

	if _, err := writer.WriteString(code); err != nil {
		return err
	}
	return writer.Flush()
}
