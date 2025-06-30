package dfa

/*
The aim of this file is to give a state machine which can detect identifier lexemmes
*/

type identifierDfaState int

const (
	identifier_start identifierDfaState = iota
	identifier_mid
	identifier_end
)

type IdentifierDFA struct {
	state identifierDfaState
}

func (dfa *IdentifierDFA) Initialize() {
	dfa.state = identifier_start
}

func (dfa *IdentifierDFA) Step(input rune) DfaReturn {
	if dfa.state == -1 {
		return INVALID
	}
	if dfa.state == identifier_start {
		if input == '_' {
			dfa.state = identifier_mid
			return VALID
		}
		if isAlphabet(input) {
			dfa.state = identifier_end
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	if dfa.state == identifier_mid {
		if isAlphabet(input) {
			dfa.state = identifier_end
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	if dfa.state == identifier_end {
		if isAlphaNumeric(input) || input == '_' {
			return VALID
		}
		dfa.state = -1 // Transition to an invalid state
		return INVALID
	}
	dfa.state = -1
	return INVALID
}

func isAlphabet(ip rune) bool {
	return (ip >= 'A' && ip <= 'Z') || (ip >= 'a' && ip <= 'z')
}

func isAlphaNumeric(ip rune) bool {
	return (ip >= '0' && ip <= '9') || isAlphabet(ip)
}

func (dfa *IdentifierDFA) Reset() {
	dfa.state = identifier_start
}
