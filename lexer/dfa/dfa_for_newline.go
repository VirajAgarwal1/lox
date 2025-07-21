package dfa

/*
The aim of this file is to give a state machine which can detect newline runes
*/

type newlineDfaState int

const (
	newline_start newlineDfaState = iota
	newline_newline
)

type NewlineDFA struct {
	state newlineDfaState
}

func (dfa *NewlineDFA) Initialize() {
	dfa.state = newline_start
}

func (dfa *NewlineDFA) Step(input rune) DfaResult {
	if dfa.state == -1 {
		return INVALID
	}
	if dfa.state == newline_start {
		if input == '\n' {
			dfa.state = newline_newline
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	if dfa.state == newline_newline {
		if input == '\n' {
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	dfa.state = -1
	return INVALID
}

func (dfa *NewlineDFA) Reset() {
	dfa.state = newline_start
}
