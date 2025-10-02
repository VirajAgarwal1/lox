package streamable_parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/VirajAgarwal1/lox/lexer"
	dfa "github.com/VirajAgarwal1/lox/lexer/dfa"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/ebnf_to_bnf"
	utils "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/utils"
)

type StackElemType byte

type EmitElemType byte

// Represents one sequence of grammar elements and its First set
type GrammarSequence struct {
	Elements []utils.Grammar_element    // clearer than "definition"
	FirstSet map[dfa.TokenType]struct{} // more descriptive than "first"
}

// Represents all production rules for a non-terminal
type ProductionRule struct {
	Sequences []GrammarSequence          // clearer than "definitions"
	FollowSet map[dfa.TokenType]struct{} // more descriptive than "follow"
}

// StackElem represents an event emitted by the parser during parsing. It can be either a non-terminal expansion or a leaf (token) emission. It is also the type of the object in the stack
type StackElem struct {
	Type         StackElemType // kind of emit: start, end or leaf
	NonTermName  string
	TerminalType dfa.TokenType
}

type EmitElem struct {
	Type    EmitElemType // kind of emit: start, end, leaf or error
	Content string       // name of the non-terminal (valid for start/end) or error message (for error event)
	Leaf    *lexer.Token
}

// StreamableParser represents the LL(1) parser state machine. It maintains a parsing stack and a lexical scanner to consume tokens.
type StreamableParser struct {
	stack   []StackElem                    // the parser’s working stack (terminals & non-terminals)
	scanner *lexer.BufferedLexicalAnalyzer // the input token stream
}

const (
	StackElemType_Start StackElemType = iota
	StackElemType_End
	StackElemType_Leaf
)

const (
	EmitElemType_Start EmitElemType = iota
	EmitElemType_End
	EmitElemType_Leaf
	EmitElemType_Error
)

func in_first_of_non_term(tok *lexer.Token, non_term string) int {
	seqs := grammarRules[non_term].Sequences
	for i := range seqs {
		if utils.InFirstSet(tok, seqs[i].FirstSet) {
			return i
		}
	}
	return -1
}
func in_follow_of_non_term(tok *lexer.Token, non_term string) bool {
	folllow_set := grammarRules[non_term].FollowSet
	for el := range folllow_set {
		if el == tok.TypeOfToken {
			return true
		}
	}
	return false
}

func (sp *StreamableParser) stack_pop() StackElem {
	top := sp.stack_peek()
	sp.stack = sp.stack[:len(sp.stack)-1]
	return top
}
func (sp *StreamableParser) stack_push(el *StackElem) {
	sp.stack = append(sp.stack, *el)
}
func (sp *StreamableParser) stack_peek() StackElem {
	return sp.stack[len(sp.stack)-1]
}

func (sp *StreamableParser) Initialize(scanner *lexer.BufferedLexicalAnalyzer) {
	sp.stack = make([]StackElem, 0, 30)
	sp.scanner = scanner

	sp.stack = append(sp.stack, StackElem{Type: StackElemType_Start, NonTermName: StartingNonTerminal})
}
func (sp *StreamableParser) EmitEvent(err error, el *StackElem, tok *lexer.Token) *EmitElem {
	if err != nil {
		return &EmitElem{
			Type:    EmitElemType_Error,
			Content: err.Error(),
		}
	}

	switch el.Type {
	case StackElemType_Start:
		if strings.HasPrefix(el.NonTermName, ebnf_to_bnf.Artificial_non_term_prefix) {
			return nil
		}
		return &EmitElem{
			Type:    EmitElemType_Start,
			Content: el.NonTermName,
		}

	case StackElemType_End:
		if strings.HasPrefix(el.NonTermName, ebnf_to_bnf.Artificial_non_term_prefix) {
			return nil
		}
		return &EmitElem{
			Type:    EmitElemType_End,
			Content: el.NonTermName,
		}

	case StackElemType_Leaf:
		return &EmitElem{
			Type:    EmitElemType_Leaf,
			Content: string(tok.Lexemme),
			Leaf:    tok,
		}
	}

	return nil
}
func (sp *StreamableParser) Parse() *EmitElem {

	if len(sp.stack) < 1 {
		next_tok, err := sp.scanner.ReadToken()
		if err != nil {
			return &EmitElem{
				Type:    EmitElemType_Error,
				Content: err.Error(),
			}
		}
		return &EmitElem{
			Type:    EmitElemType_Error,
			Content: fmt.Sprintf("Expected \"EOF\" but got \"%s\"", next_tok.ToString()),
		}
	}

	for {

		top := sp.stack_peek()
		lookahead_token, err := sp.scanner.Peek()
		if err != nil && err != io.EOF {
			return &EmitElem{
				Type:    EmitElemType_Error,
				Content: err.Error(),
			}
		}

		switch top.Type {
		case StackElemType_Leaf:
			if top.TerminalType == utils.Epsilon {
				sp.stack_pop()
				continue
			}
			if lookahead_token.TypeOfToken == top.TerminalType {
				sp.stack_pop()
				sp.scanner.ReadToken()
				return sp.EmitEvent(nil, &top, lookahead_token)
			}
			// Else consume tokens to put into the error message until, we get to a token, which is equal to the top
			err_start := strconv.FormatUint(uint64(lookahead_token.Line), 10) + "," + strconv.FormatUint(uint64(lookahead_token.Offset), 10)
			for lookahead_token.TypeOfToken != top.TerminalType {
				sp.scanner.ReadToken()
				lookahead_token, err = sp.scanner.Peek()
				// TODO: Idealy, an error which occurs while doing error recovery should also be reported... But for simplicity's sake, I will ignore this error for now
			}
			err_end := strconv.FormatUint(uint64(lookahead_token.Line), 10) + "," + strconv.FormatUint(uint64(lookahead_token.Offset), 10)
			return sp.EmitEvent(
				fmt.Errorf("parse error from %s to %s, expected \"%v\"", err_start, err_end, top.TerminalType),
				nil,
				lookahead_token,
			)

		case StackElemType_End:
			sp.stack_pop()
			output := sp.EmitEvent(nil, &top, lookahead_token)
			if output == nil {
				continue
			}
			return output

		case StackElemType_Start:
			prod_rule := in_first_of_non_term(lookahead_token, top.NonTermName)
			sp.stack_pop()
			if prod_rule != -1 {
				sp.stack_push(&StackElem{
					Type:        StackElemType_End,
					NonTermName: top.NonTermName,
				})
				for i := len(grammarRules[top.NonTermName].Sequences[prod_rule].Elements) - 1; i > -1; i-- {
					if grammarRules[top.NonTermName].Sequences[prod_rule].Elements[i].IsNonTerminal {
						sp.stack_push(&StackElem{
							Type:        StackElemType_Start,
							NonTermName: grammarRules[top.NonTermName].Sequences[prod_rule].Elements[i].Non_term_name,
						})
					} else {
						sp.stack_push(&StackElem{
							Type:         StackElemType_Leaf,
							TerminalType: grammarRules[top.NonTermName].Sequences[prod_rule].Elements[i].Terminal_type,
						})
					}
				}
				output := sp.EmitEvent(nil, &top, lookahead_token)
				if output == nil {
					continue
				}
				return output
			}
			err_start := strconv.FormatUint(uint64(lookahead_token.Line), 10) + "," + strconv.FormatUint(uint64(lookahead_token.Offset), 10)
			for !in_follow_of_non_term(lookahead_token, top.NonTermName) {
				sp.scanner.ReadToken()
				lookahead_token, err = sp.scanner.Peek()
				// TODO: Idealy, an error which occurs while doing error recovery should also be reported... But for simplicity's sake, I will ignore this error for now
			}
			err_end := strconv.FormatUint(uint64(lookahead_token.Line), 10) + "," + strconv.FormatUint(uint64(lookahead_token.Offset), 10)
			return sp.EmitEvent(
				fmt.Errorf("parse error from %s to %s", err_start, err_end),
				nil,
				lookahead_token,
			)
		}
	}
}
