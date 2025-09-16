package first_follow

import (
	"github.com/VirajAgarwal1/lox/lexer/dfa"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/utils"
)

var bnfGrammar = map[string]([][]utils.Grammar_element){}

// A U (B - {E})
func union_first_sets_wo_epsilon(A []dfa.TokenType, B []dfa.TokenType) []dfa.TokenType {
	for _, tok := range B {
		if tok != utils.Epsilon && !utils.Contains(A, tok) {
			A = append(A, tok)
		}
	}
	return A
}

// A U B
func union_first_sets(A []dfa.TokenType, B []dfa.TokenType) []dfa.TokenType {
	for _, tok := range B {
		if !utils.Contains(A, tok) {
			A = append(A, tok)
		}
	}
	return A
}
