package dfa

import "unicode"

/*
The aim of this file is to give a state machine which can detecting all whitespace execpt for newline in Unicode (UTF-8 encoding)
*/

type whitespaceDfaState int

const (
	whitespace_start whitespaceDfaState = iota
	whitespace_blanks
)

type WhitespaceDFA struct {
	state whitespaceDfaState
}

func (dfa *WhitespaceDFA) Initialize() {
	dfa.state = whitespace_start
}

func (dfa *WhitespaceDFA) Step(input rune) DfaReturn {
	if dfa.state == whitespace_start {
		if unicode.IsSpace(input) {
			dfa.state = whitespace_blanks
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	if dfa.state == whitespace_blanks {
		if unicode.IsSpace(input) {
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	dfa.state = -1
	return INVALID
}

func (dfa *WhitespaceDFA) Reset() {
	dfa.state = whitespace_start
}
