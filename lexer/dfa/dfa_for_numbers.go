package dfa

/*
The aim of this file is to give a state machine which can detect a number lexeme
*/

type numberDfaState int

const (
	number_start numberDfaState = iota
	number_before_decimal
	number_decimal
	number_after_decimal
)

type NumberDFA struct {
	state numberDfaState
}

func (dfa *NumberDFA) Initialize() {
	dfa.state = number_start
}

func (dfa *NumberDFA) Step(input rune) DfaReturn {
	if dfa.state == number_start {
		if isNumber(input) {
			dfa.state = number_before_decimal
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	if dfa.state == number_before_decimal {
		if input == '.' {
			dfa.state = number_decimal
			return INTERMEDIATE
		}
		if isNumber(input) {
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	if dfa.state == number_decimal {
		if isNumber(input) {
			dfa.state = number_after_decimal
			return VALID
		}
		dfa.state = -1 // Transition to an invalid state
		return INVALID
	}
	if dfa.state == number_after_decimal {
		if isNumber(input) {
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	dfa.state = -1
	return INVALID
}

func (dfa *NumberDFA) Reset() {
	dfa.state = number_start
}

func isNumber(input rune) bool {
	return input >= '0' && input <= '9'
}
