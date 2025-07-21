package dfa

/*
The aim of this file is to give a state machine which can detect comments
*/

type commentDfaState int

const (
	comment_start commentDfaState = iota
	comment_first_slash
	comment_second_slash
	comment_newline
)

type CommentDFA struct {
	state commentDfaState
}

func (dfa *CommentDFA) Initialize() {
	dfa.state = comment_start
}

func (dfa *CommentDFA) Step(input rune) DfaResult {
	if dfa.state == -1 {
		return INVALID
	}
	if dfa.state == comment_start {
		if input == '/' {
			dfa.state = comment_first_slash
			return INTERMEDIATE
		}
		dfa.state = -1
		return INVALID
	}
	if dfa.state == comment_first_slash {
		if input == '/' {
			dfa.state = comment_second_slash
			return VALID
		}
		dfa.state = -1
		return INVALID
	}
	if dfa.state == comment_second_slash {
		if input == '\n' {
			dfa.state = comment_newline
			return INVALID
		}
		return VALID
	}
	dfa.state = -1
	return INVALID
}

func (dfa *CommentDFA) Reset() {
	dfa.state = comment_start
}
