package lexer

import (
	"bufio"

	dfa "github.com/VirajAgarwal1/lox/lexer/dfa"
)

/*
	The order of conversion is this:
		Bytes -> Runes -> Lexemes -> Tokens
*/

type Token struct {
	TypeOfToken dfa.TokenType
	Line        int
	Offset      int
}

type LexicalAnalyzer struct {
	source     *bufio.Reader
	tokenDFAa  map[dfa.TokenType]dfa.DFA
	lineNum    int
	lineOffset int
	lexemme    []rune
}

func (scanner *LexicalAnalyzer) Initialize(source *bufio.Reader) {
	scanner.tokenDFAa = dfa.GenerateDFAs()
	scanner.source = source
}

func rec_token_parser(scanner *LexicalAnalyzer, source *bufio.Reader, state []rune, lexemmes_to_check []bool, input rune) {
	new_state := append(state, input)
	next_input, _, _ := source.ReadRune() // TODO: Handle the edge cases of reading runes here
	filtered_lexemmes_to_check := make([]bool, len(lexemmes_to_check), len(lexemmes_to_check))
	copy(filtered_lexemmes_to_check, lexemmes_to_check)

	foundValids := make([]bool, len(lexemmes_to_check), len(lexemmes_to_check))
	numToCheckCont := 0

	for i, to_check := range filtered_lexemmes_to_check {
		if !to_check {
			continue
		}
		dfaStepResult := scanner.tokenDFAa[dfa.TokensList[i]].Step(next_input)
		if dfaStepResult.IsInvalid() {
			filtered_lexemmes_to_check[i] = false
			continue
		}
		numToCheckCont++
		if dfaStepResult.IsValid() {
			foundValids[i] = true
		}
	}
	if numToCheckCont == 0 {
		return false, 
	}
	rec_foundValidToken,  := rec_token_parser(scanner, source, new_state, filtered_lexemmes_to_check, next_input)
}

func (scanner *LexicalAnalyzer) PipeTokens() Token {
	// TODO: Add the scanner logic here
	for {
		// Read one rune from the source

		// Execute a step in all the dfas (which are were marked for checking in the previous loop) with the current rune.

		// If atleast one dfa in the array of tokens' dfa returns valid/intermediate
			// If atleast one valid is there in the array then have a mark of it outside of the loop in some variable
			// If no valid but atleast one intermedite is there, then we dont mark the variable outside the loop. In case, if next runes do not satisfy we do need a way to say that last ones were intermediates for lexical sytanx error reporting
			// 

		// Feed the rune to every
	}
}


var
a
=
9