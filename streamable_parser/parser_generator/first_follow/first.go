package first_follow

import (
	"fmt"

	"github.com/VirajAgarwal1/lox/lexer/dfa"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/utils"
)

type FirstSetInfo struct {
	FirstForNonTerminal []dfa.TokenType
	FirstForDefinitions [][]dfa.TokenType
}

var firstSets = map[string]FirstSetInfo{}

func ComputeFirstForSequence(definition []utils.Grammar_element) []dfa.TokenType {

	final_first_set := []dfa.TokenType{}

	for i, elem := range definition {

		temp_first_set := FirstSetInfo{}
		if elem.IsNonTerminal {
			temp_first_set = ComputeFirstForNonTerminal(elem.Non_term_name)
		} else {
			temp_first_set = FirstSetInfo{[]dfa.TokenType{elem.Terminal_type}, nil}
		}

		if utils.Contains(temp_first_set.FirstForNonTerminal, utils.Epsilon) {
			final_first_set = union_sets_wo_epsilon(final_first_set, temp_first_set.FirstForNonTerminal)
			if i == len(definition)-1 {
				final_first_set = append(final_first_set, utils.Epsilon)
			}
			continue
		}
		final_first_set = union_sets(final_first_set, temp_first_set.FirstForNonTerminal)
		break
	}

	return final_first_set
}
func ComputeFirstForAlternatives(definitions [][]utils.Grammar_element) FirstSetInfo {

	final_first_set := []dfa.TokenType{}
	first_sets := [][]dfa.TokenType{}
	for _, def := range definitions {
		temp_first_set := ComputeFirstForSequence(def)
		final_first_set = union_sets(final_first_set, temp_first_set)
		first_sets = append(first_sets, temp_first_set)
	}
	return FirstSetInfo{FirstForNonTerminal: final_first_set, FirstForDefinitions: first_sets}
}
func ComputeFirstForNonTerminal(non_term string) FirstSetInfo {

	first_set, processed_already := firstSets[non_term]
	if processed_already { // Caching
		return first_set
	}

	definitions, found := bnfGrammar[non_term]
	if !found {
		panic(fmt.Sprintf("Trying to access a non-terminal %s which doesnt exist in the BNF grammar", non_term))
	}
	firstSets[non_term] = ComputeFirstForAlternatives(definitions)

	return firstSets[non_term]
}

func ComputeFirstSets(bnf_grammar map[string]([][]utils.Grammar_element)) map[string]FirstSetInfo {

	bnfGrammar = bnf_grammar

	for non_term := range bnf_grammar {
		firstSets[non_term] = ComputeFirstForNonTerminal(non_term)
	}

	return firstSets
}
