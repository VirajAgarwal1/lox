package lexer_demos

import (
	"fmt"
	"time"

	dfa "github.com/VirajAgarwal1/lox/lexer/dfa"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

func colorByState(state dfa.DfaResult) string {
	switch state {
	case dfa.INVALID:
		return Red + "✗"
	case dfa.INTERMEDIATE:
		return Yellow + "⚡"
	case dfa.VALID:
		return Green + "✓"
	}
	return "?"
}

func demonstrateInput(input string, tokenTypes []dfa.TokenType) {
	fmt.Printf("\n%s%s🔍 Testing: \"%s\"%s\n", Bold+Cyan, "", input, Reset)
	fmt.Printf("%sWe'll now feed the input into each DFA one character at a time.%s\n", "", Reset)
	fmt.Printf("%sYou’ll see how each DFA changes state as it reads each character.%s\n", "", Reset)

	statesManager := dfa.DFAStatesManager{}
	statesManager.Initialize()
	dfas := statesManager.DfaForToken

	// Process each character
	for i, char := range input {
		fmt.Printf("   Step %d: '%c' → ", i+1, char)

		for i, token := range dfa.TokensList {
			j := -1
			for _, wantedToken := range tokenTypes {
				if token == wantedToken {
					j = i
				}
			}
			if j == -1 {
				continue
			}
			result := dfas[i].Step(char)
			fmt.Printf("%s%s %s%s",
				colorByState(result),
				string(token),
				Reset,
				func() string {
					if j < len(tokenTypes)-1 {
						return " | "
					} else {
						return ""
					}
				}())
		}
		fmt.Println()
		time.Sleep(400 * time.Millisecond)
	}

	fmt.Printf("%sAnalysis complete for \"%s\". Final DFA states reflect how each token type interpreted the input.%s\n", Yellow, input, Reset)
}

func DfaDemo() {
	// Title
	fmt.Printf("\n%s%s🚀 DFA Magic Demo 🚀%s\n", Bold+Cyan, "", Reset)
	fmt.Printf("%sWelcome! This demo shows how Deterministic Finite Automata (DFAs) recognize tokens in a programming language, step by step.%s\n", "", Reset)
	fmt.Printf("%sEach DFA represents a possible token (like 'while', '==', or a number).%s\n", "", Reset)
	fmt.Printf("%sAs we process each character, you’ll see how each DFA responds:✓ (Valid), ⚡ (Intermediate), ✗ (Invalid)%s\n", "", Reset)

	// Demo 1: Keyword vs Identifier
	fmt.Printf("\n%s%s═══ Magic #1: Keywords vs Identifiers ═══%s", Bold+Blue, "", Reset)
	fmt.Printf("\n%sIn this example, we compare the keyword 'while' with the identifier 'whiles'.%s\n", "", Reset)
	fmt.Printf("%sBoth start the same, but only 'while' exactly matches the keyword DFA. 'whiles' becomes an identifier instead.%s\n", "", Reset)
	demonstrateInput("while", []dfa.TokenType{dfa.WHILE, dfa.IDENTIFIER})
	demonstrateInput("whiles", []dfa.TokenType{dfa.WHILE, dfa.IDENTIFIER})

	// Demo 2: Maximal Munching
	fmt.Printf("\n%s%s═══ Magic #2: Maximal Munching ═══%s", Bold+Blue, "", Reset)
	fmt.Printf("\n%sThis demonstrates the principle of 'maximal munching': we always consume as many characters as possible to form the longest valid token.%s\n", "", Reset)
	fmt.Printf("%sWatch how '=' is valid for both '=' and '==', but when we add another '=', the DFA for '==' succeeds and '=' fails.%s\n", "", Reset)
	demonstrateInput("=", []dfa.TokenType{dfa.EQUAL, dfa.EQUAL_EQUAL})
	demonstrateInput("==", []dfa.TokenType{dfa.EQUAL, dfa.EQUAL_EQUAL})

	// Demo 3: Complex Tokens
	fmt.Printf("\n%s%s═══ Magic #3: Complex Patterns ═══%s", Bold+Blue, "", Reset)
	fmt.Printf("\n%sHere, we test more complex token patterns like numbers and strings.%s\n", "", Reset)
	fmt.Printf("%sNotice how DFAs carefully validate character-by-character — numbers handle digits and dots, and strings expect surrounding quotes.%s\n", "", Reset)
	demonstrateInput("123.45", []dfa.TokenType{dfa.NUMBER})
	demonstrateInput("\"hello\"", []dfa.TokenType{dfa.STRING})

	// Finale
	fmt.Printf("\n%s%s🎉 That's the magic of DFAs! 🎉%s\n", Bold+Green, "", Reset)
	fmt.Printf("%sYou just watched multiple DFAs process the same input in parallel!%s\n", "", Reset)
	fmt.Printf("%s✓ = Valid token, ⚡ = Could become valid (intermediate), ✗ = Invalid so far.%s\n", "", Reset)
	fmt.Printf("%sTry your own inputs and token types to see how the lexer behaves.%s\n", "", Reset)
}
