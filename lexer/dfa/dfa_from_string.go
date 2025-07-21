package dfa

/*
The purpose of this file is to provide a function which can take a string and give back a state machine which can be in 3 states:
	1. VALID -> if the slice of runes up until now, match the string
	2. INVALID -> if the slice of runes up until now, has a foreign rune in it and thus cannot possible match the string
	3. INTERMEDIATE -> if the slice of runes up until now, doens't match the string but if the right sequence of runes follow in the future then there is possibility of getting to "VALID" state
*/

type InputStringDFA struct {
	str   []rune
	state int
}

func (dfa *InputStringDFA) Initialize(_str string) {
	dfa.str = []rune(_str)
	dfa.state = 0
}

func (dfa *InputStringDFA) Step(input rune) DfaResult {
	if dfa.state == -1 {
		return INVALID
	}
	if input == dfa.str[dfa.state] {
		if dfa.state == len(dfa.str)-1 {
			dfa.state = -1
			return VALID
		}
		dfa.state++
		return INTERMEDIATE
	}
	dfa.state = -1
	return INVALID
}

func (dfa *InputStringDFA) Reset() {
	dfa.state = 0
}
