package code_snippets

import (
	"github.com/VirajAgarwal1/lox/lexer/dfa"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/first_follow"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/utils"
)

func code_tokens_set(tokens []dfa.TokenType) string {
	start := "map[dfa.TokenType]struct{}{"
	middle := "\n"
	for _, term := range tokens {
		middle += utils.Indent_lines(
			utils.String_to_type_string[string(term)]+": {},",
			1,
		)
	}
	middle += "\n"
	end := "}"

	if len(tokens) < 1 {
		return start + end
	}
	return start + middle + end
}

func code_a_Grammar_elem(el utils.Grammar_element) string {
	if el.IsNonTerminal {
		return "{IsNonTerminal: true, Non_term_name: \"" + el.Non_term_name + "\"},"
	}
	return "{IsNonTerminal: false, Terminal_type: " + utils.String_to_type_string[string(el.Terminal_type)] + "},"
}

func code_Grammar_elements(elems []utils.Grammar_element) string {
	start := "[]utils.Grammar_element{"
	middle := "\n"
	for _, el := range elems {
		middle += utils.Indent_lines(
			code_a_Grammar_elem(el),
			1,
		)
	}
	middle += "\n"
	end := "}"

	if len(elems) < 1 {
		return start + end
	}
	return start + middle + end
}

func code_GrammarSequence(first_set []dfa.TokenType, def []utils.Grammar_element) string {

	start := "{\n"
	middle := utils.Indent_lines(
		"FirstSet: "+code_tokens_set(first_set)+",\nElements: "+code_Grammar_elements(def)+",",
		1,
	)
	end := "\n}"

	return start + middle + end
}

func code_slice_of_GrammarSequence(first_sets [][]dfa.TokenType, definitions [][]utils.Grammar_element) string {

	if len(first_sets) != len(definitions) {
		panic("thier length should be the same")
	}

	start := "[]GrammarSequence{"
	middle := "\n"
	for i := range len(first_sets) {
		middle += utils.Indent_lines(
			code_GrammarSequence(first_sets[i], definitions[i])+",",
			1,
		)
	}
	middle += "\n"
	end := "}"

	if len(first_sets) < 1 {
		return start + end
	}
	return start + middle + end
}

func code_non_terminal(non_term string, follow_set []dfa.TokenType, first_sets [][]dfa.TokenType, definitions [][]utils.Grammar_element) string {
	if len(non_term) < 1 || len(follow_set) < 1 {
		return ""
	}

	start := "\"" + non_term + "\": {"
	middle := "\n"
	middle += utils.Indent_lines(
		"FollowSet: "+code_tokens_set(follow_set)+",",
		1,
	)
	middle += "\n"
	middle += utils.Indent_lines(
		"Sequences: "+code_slice_of_GrammarSequence(first_sets, definitions)+",",
		1,
	)
	middle += "\n"
	end := "}"

	return start + middle + end
}

func GrammarRules_code(bnf_grammar map[string]([][]utils.Grammar_element), firstSet map[string]first_follow.FirstSetInfo, followSet map[string]([]dfa.TokenType)) string {

	if len(bnf_grammar) != len(firstSet) && len(firstSet) != len(followSet) {
		panic("they should have same number of non terminals")
	}

	start := "var grammarRules = map[string]ProductionRule{"
	middle := "\n"
	for non_term := range bnf_grammar {
		middle += utils.Indent_lines(
			code_non_terminal(non_term, followSet[non_term], firstSet[non_term].FirstForDefinitions, bnf_grammar[non_term])+",",
			1,
		)
	}
	middle += "\n"
	end := "}"

	if len(bnf_grammar) < 1 {
		return start + end
	}
	return start + middle + end
}
