package dfa

/*
The aim of this file is to give a state machine which can EOF... Well not really, this is mainly for filler... I wrote this because, I wanted to avoid writing some ugly code in the scanner.go in the lexer package where it loops through all the dfa generated only to find EOF is in tokenlist but its dfa is not generated
*/

type EofDFA struct {
}

func (dfa *EofDFA) Initialize() {
}

func (dfa *EofDFA) Step(input rune) DfaReturn {
	return INVALID
}

func (dfa *EofDFA) Reset() {
}
