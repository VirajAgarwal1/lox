// -----------------------------------
// CODE INDEPENDANT OF GRAMMAR START
// -----------------------------------

package parser

import (
	"io"

	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/lexer/dfa"
)

var awiue lexer.BufferedLexer = lexer.BufferedLexer{}

type Value struct {
	LoxType string
	Inner   any
}
type Node interface {
	Evaluate() *Value
}
type Literal struct {
	Value *lexer.Token
}
type ParseFunc func(*lexer.BufferedLexer) ([]Node, bool, error)

func (non_terminal *Literal) Evaluate() *Value {
	return &Value{
		LoxType: determineLoxType(non_terminal.Value),
		Inner:   non_terminal.Value.Lexemme,
	}
}

func determineLoxType(tok *lexer.Token) string {
	if tok.TypeOfToken == dfa.NUMBER {
		return "number"
	}
	if tok.TypeOfToken == dfa.STRING {
		return "string"
	}
	if tok.TypeOfToken == dfa.TRUE || tok.TypeOfToken == dfa.FALSE {
		return "bool"
	}
	if tok.TypeOfToken == dfa.NIL {
		return "nil"
	}
	return string(tok.TypeOfToken)
}

// -------------------- COMBINATOR HELPERS --------------------

func matchToken(t dfa.TokenType) ParseFunc {
	return func(buf *lexer.BufferedLexer) ([]Node, bool, error) {
		chk := buf.MakeCheckpoint()
		tok, err := buf.ReadToken()
		if err != nil && err != io.EOF {
			buf.RollbackTo(chk)
			return nil, false, err
		}
		if tok.TypeOfToken == t {
			return []Node{&Literal{tok}}, true, nil
		}
		buf.RollbackTo(chk)
		return nil, false, nil
		// fmt.Errorf("Unexpected token '%v' found at line %d, offset %d. Expected token '%v'", string(tok.TypeOfToken), tok.Line, tok.Offset, string(t))
	}
}

func sequence(parts ...ParseFunc) ParseFunc {
	return func(buf *lexer.BufferedLexer) ([]Node, bool, error) {
		chk := buf.MakeCheckpoint()
		output := []Node{}
		for _, part := range parts {
			nodes, ok, err := part(buf)
			if err != nil || !ok {
				buf.RollbackTo(chk)
				return nil, false, err
			}
			output = append(output, nodes...)
		}
		return output, true, nil
	}
}

func choice(parts ...ParseFunc) ParseFunc {
	return func(buf *lexer.BufferedLexer) ([]Node, bool, error) {
		chk := buf.MakeCheckpoint()
		for _, part := range parts {
			nodes, ok, err := part(buf)
			if err != nil {
				buf.RollbackTo(chk)
				return nil, false, err
			}
			if ok {
				return nodes, true, nil
			}
		}
		buf.RollbackTo(chk)
		return nil, false, nil
	}
}

func zeroOrMore(part ParseFunc) ParseFunc {
	return func(buf *lexer.BufferedLexer) ([]Node, bool, error) {
		chk := buf.MakeCheckpoint()
		output := []Node{}
		for {
			nodes, ok, err := part(buf)
			if err != nil || !ok {
				buf.RollbackTo(chk)
				break
			}
			output = append(output, nodes...)
		}
		return output, true, nil
	}
}

func oneOrMore(part ParseFunc) ParseFunc {
	return func(buf *lexer.BufferedLexer) ([]Node, bool, error) {
		chk := buf.MakeCheckpoint()
		output := []Node{}
		nodes, ok, err := part(buf)
		if err != nil || !ok {
			buf.RollbackTo(chk)
			return nil, false, err
		}
		output = append(output, nodes...)

		for {
			nodes, ok, err = part(buf)
			if err != nil || !ok {
				buf.RollbackTo(chk)
				break
			}
			output = append(output, nodes...)
		}
		return output, true, nil
	}
}

// -----------------------------------
// CODE INDEPENDANT OF GRAMMAR END
// -----------------------------------

type Grammar_primary struct {
	Arguments []Node
}
type Grammar_equality struct {
	Arguments []Node
}
type Grammar_term struct {
	Arguments []Node
}
type Grammar_unary struct {
	Arguments []Node
}
type Grammar_grammar struct {
	Arguments []Node
}
type Grammar_expression struct {
	Arguments []Node
}
type Grammar_comma struct {
	Arguments []Node
}
type Grammar_comparison struct {
	Arguments []Node
}
type Grammar_factor struct {
	Arguments []Node
}

func (non_terminal *Grammar_equality) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_term) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_unary) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_grammar) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_expression) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_comma) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_comparison) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_factor) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_primary) Evaluate() *Value {
	return nil
}

func Parse_primary(buf *lexer.BufferedLexer) ([]Node, bool, error) {
	output := Grammar_primary{}

	args, ok, err := choice(
		matchToken(dfa.IDENTIFIER),
		matchToken(dfa.NUMBER),
		matchToken(dfa.STRING),
		matchToken(dfa.TRUE),
		matchToken(dfa.FALSE),
		matchToken(dfa.NIL),
		sequence(
			matchToken(dfa.LEFT_PAREN),
			Parse_expression,
			matchToken(dfa.RIGHT_PAREN),
		),
	)(buf)

	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
}
func Parse_equality(buf *lexer.BufferedLexer) ([]Node, bool, error) {
	output := Grammar_equality{}

	args, ok, err := sequence(
		Parse_comparison,
		zeroOrMore(
			sequence(
				choice(
					matchToken(dfa.BANG_EQUAL),
					matchToken(dfa.EQUAL_EQUAL),
				),
				Parse_comparison,
			),
		),
	)(buf)

	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
}
func Parse_term(buf *lexer.BufferedLexer) ([]Node, bool, error) {
	output := Grammar_term{}

	args, ok, err := sequence(
		Parse_factor,
		zeroOrMore(
			sequence(
				choice(
					matchToken(dfa.MINUS),
					matchToken(dfa.PLUS),
				),
				Parse_factor,
			),
		),
	)(buf)

	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
}
func Parse_unary(buf *lexer.BufferedLexer) ([]Node, bool, error) {
	output := Grammar_unary{}

	args, ok, err := choice(
		sequence(
			choice(
				matchToken(dfa.BANG),
				matchToken(dfa.MINUS),
			),
			Parse_unary,
		),
		Parse_primary,
	)(buf)

	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
}
func Parse_grammar(buf *lexer.BufferedLexer) ([]Node, bool, error) {
	output := Grammar_grammar{}

	args, ok, err := sequence(
		zeroOrMore(
			sequence(
				zeroOrMore(
					matchToken(dfa.NEWLINE),
				),
				Parse_expression,
				zeroOrMore(
					matchToken(dfa.NEWLINE),
				),
			),
		),
		matchToken(dfa.EOF),
	)(buf)

	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
}
func Parse_expression(buf *lexer.BufferedLexer) ([]Node, bool, error) {
	output := Grammar_expression{}

	args, ok, err := Parse_comma(buf)

	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
}
func Parse_comma(buf *lexer.BufferedLexer) ([]Node, bool, error) {
	output := Grammar_comma{}

	args, ok, err := sequence(
		Parse_equality,
		zeroOrMore(
			sequence(
				matchToken(dfa.COMMA),
				Parse_equality,
			),
		),
	)(buf)

	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
}
func Parse_comparison(buf *lexer.BufferedLexer) ([]Node, bool, error) {
	output := Grammar_comparison{}

	args, ok, err := sequence(
		Parse_term,
		zeroOrMore(
			sequence(
				choice(
					matchToken(dfa.GREATER),
					matchToken(dfa.GREATER_EQUAL),
					matchToken(dfa.LESS),
					matchToken(dfa.LESS_EQUAL),
				),
				Parse_term,
			),
		),
	)(buf)

	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
}
func Parse_factor(buf *lexer.BufferedLexer) ([]Node, bool, error) {
	output := Grammar_factor{}

	args, ok, err := sequence(
		Parse_unary,
		zeroOrMore(
			sequence(
				choice(
					matchToken(dfa.SLASH),
					matchToken(dfa.STAR),
				),
				Parse_unary,
			),
		),
	)(buf)

	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
}
