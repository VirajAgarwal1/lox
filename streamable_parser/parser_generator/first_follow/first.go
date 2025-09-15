package first_follow

import (
	"fmt"

	"github.com/VirajAgarwal1/lox/lexer/dfa"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/utils"
)

var first = map[string]([]dfa.TokenType){}

func process_sequence_first(definition []utils.Grammar_element) []dfa.TokenType {

	final_first_set := []dfa.TokenType{}

	for i, elem := range definition {

		temp_first_set := []dfa.TokenType{}
		if elem.IsNonTerminal {
			temp_first_set = process_non_terminal_first(elem.Non_term_name)
		} else {
			temp_first_set = []dfa.TokenType{elem.Terminal_type}
		}

		if utils.Contains(temp_first_set, utils.Epsilon) {
			final_first_set = union_first_sets_wo_epsilon(final_first_set, temp_first_set)
			if i == len(definition)-1 {
				final_first_set = append(final_first_set, utils.Epsilon)
			}
			continue
		}
		final_first_set = union_first_sets(final_first_set, temp_first_set)
		break
	}

	return final_first_set
}
func process_or_first(definitions [][]utils.Grammar_element) []dfa.TokenType {

	final_first_set := []dfa.TokenType{}
	for _, def := range definitions {
		temp_first_set := process_sequence_first(def)
		final_first_set = union_first_sets(final_first_set, temp_first_set)
	}
	return final_first_set
}
func process_non_terminal_first(non_term string) []dfa.TokenType {

	first_set, processed_already := first[non_term]
	if processed_already { // Caching
		return first_set
	}

	definitions, found := bnf_grammar_global[non_term]
	if !found {
		panic(fmt.Sprintf("Trying to access a non-terminal %s which doesnt exist in the BNF grammar", non_term))
	}
	first[non_term] = process_or_first(definitions)

	return first[non_term]
}

func ComputeFirstFromBNF(bnf_grammar map[string]([][]utils.Grammar_element)) map[string]([]dfa.TokenType) {

	bnf_grammar_global = bnf_grammar

	for non_term := range bnf_grammar {
		first[non_term] = process_non_terminal_first(non_term)
	}

	return first
}
