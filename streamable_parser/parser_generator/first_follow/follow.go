package first_follow

import (
	"github.com/VirajAgarwal1/lox/lexer/dfa"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/utils"
)

var FollowSets = map[string]([]dfa.TokenType){}

func ComputeFollowFromSequence(non_term_on_lhs string, choices_idx int, non_term_idx int) []dfa.TokenType {

	follow_set := []dfa.TokenType{dfa.EOF}

	if non_term_idx >= len(bnfGrammar[non_term_on_lhs][choices_idx])-1 {
		if non_term_on_lhs == bnfGrammar[non_term_on_lhs][choices_idx][non_term_idx].Non_term_name {
			return follow_set
		}
		follow_set = union_sets(follow_set, ComputeFollowForNonTerminal(non_term_on_lhs))
		return follow_set
	}

	first_of_following_sequence := ComputeFirstForSequence(bnfGrammar[non_term_on_lhs][choices_idx][non_term_idx+1:])
	follow_set = union_sets_wo_epsilon(follow_set, first_of_following_sequence)

	if utils.Contains(first_of_following_sequence, utils.Epsilon) {
		if non_term_on_lhs == bnfGrammar[non_term_on_lhs][choices_idx][non_term_idx].Non_term_name {
			return follow_set
		}
		follow_set = union_sets(follow_set, ComputeFollowForNonTerminal(non_term_on_lhs))
	}
	return follow_set
}

func ComputeFollowForNonTerminal(non_term string) []dfa.TokenType {

	follow_set, processed_already := FollowSets[non_term]
	if processed_already {
		return follow_set
	}

	follow_set = []dfa.TokenType{dfa.EOF}
	FollowSets[non_term] = follow_set // So that this non-term is not tried again in the recursion

	for lhs_non_term, definitions := range bnfGrammar {
		for or_idx, def := range definitions {
			idx := production_contains_non_term(def, non_term)
			if idx != -1 {
				follow_set = union_sets_wo_epsilon(follow_set, ComputeFollowFromSequence(lhs_non_term, or_idx, idx))
			}
		}
	}

	FollowSets[non_term] = follow_set
	return follow_set
}

func ComputeFollowSets(bnf_grammar map[string]([][]utils.Grammar_element)) map[string]([]dfa.TokenType) {

	bnfGrammar = bnf_grammar

	for non_term := range bnf_grammar {
		FollowSets[non_term] = ComputeFollowForNonTerminal(non_term)
	}
	return FollowSets
}
