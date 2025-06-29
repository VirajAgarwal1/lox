package dfa

/*
The aim of this file is to give a state machine which can detect a string lexeme
*/

type stringDfaState int

const (
	string_start stringDfaState = iota
	string_left_apos
	string_content
	string_right_apos
)

type StringDFA struct {
	state stringDfaState
}

func (dfa *StringDFA) Initialize() {
	dfa.state = string_start
}

func (dfa *StringDFA) Step(input rune) DfaReturn {
	if dfa.state == string_start {
		if input == '"' {
			dfa.state = string_left_apos
			return INTERMEDIATE
		}
		dfa.state = -1
		return INVALID
	}
	if dfa.state == string_left_apos {
		if input == '"' {
			dfa.state = string_right_apos
			return VALID
		}
		dfa.state = string_content
		return INTERMEDIATE
	}
	if dfa.state == string_content {
		if input == '"' {
			dfa.state = string_right_apos
			return VALID
		}
		return INTERMEDIATE
	}
	dfa.state = -1
	return INVALID
}

func (dfa *StringDFA) Reset() {
	dfa.state = string_start
}
