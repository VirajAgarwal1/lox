// // -----------------------------------
// // CODE INDEPENDANT OF GRAMMAR START
// // -----------------------------------

// package parser

// import (
// 	"fmt"
// 	"io"

// 	"github.com/VirajAgarwal1/lox/errorhandler"
// 	"github.com/VirajAgarwal1/lox/lexer"
// 	"github.com/VirajAgarwal1/lox/lexer/dfa"
// )

// type Value struct {
// 	LoxType string
// 	Inner   any
// }
// type Node interface {
// 	Evaluate() *Value
// }
// type Literal struct {
// 	Value *lexer.Token
// }

// func (non_terminal *Literal) Evaluate() *Value {
// 	return &Value{
// 		LoxType: determineLoxType(non_terminal.Value),
// 		Inner:   non_terminal.Value.Lexemme,
// 	}
// }
// func match(buf_scanner *lexer.BufferedLexicalAnalyzer, token dfa.TokenType) (Node, bool, error) {
// 	t, err := buf_scanner.CurrentTokenWithoutConsume()
// 	if err != nil && err != io.EOF {
// 		return nil, false, errorhandler.RetErr("", err)
// 	}
// 	if t.TypeOfToken == token {
// 		buf_scanner.ConsumeOneToken()
// 		literal_value := Literal{
// 			Value: t,
// 		}
// 		return &literal_value, true, nil
// 	}
// 	return nil, false, nil
// }
// func determineLoxType(tok *lexer.Token) string {
// 	if tok.TypeOfToken == dfa.NUMBER {
// 		return "number"
// 	}
// 	if tok.TypeOfToken == dfa.STRING {
// 		return "string"
// 	}
// 	if tok.TypeOfToken == dfa.TRUE || tok.TypeOfToken == dfa.FALSE {
// 		return "bool"
// 	}
// 	if tok.TypeOfToken == dfa.NIL {
// 		return "nil"
// 	}
// 	return string(tok.TypeOfToken)
// }

// func parseZeroOrMore(buf *lexer.BufferedLexicalAnalyzer, part func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error)) ([]Node, bool, error) {
// 	output := []Node{}

// 	for {
// 		arg, isOk, err := part(buf)
// 		if err != nil {
// 			return output, false, err
// 		}
// 		if !isOk {
// 			return output, true, nil
// 		}
// 		output = append(output, arg...)
// 	}
// }
// func parseOneOrMore(buf *lexer.BufferedLexicalAnalyzer, part func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error)) ([]Node, bool, error) {
// 	output := []Node{}

// 	arg, isOk, err := part(buf)
// 	if err != nil || !isOk {
// 		return output, false, err
// 	}
// 	output = append(output, arg...)

// 	nextOutput, isOk, err := parseZeroOrMore(buf, part)
// 	if err != nil {
// 		return output, true, err
// 	}
// 	if !isOk {
// 		return output, true, nil
// 	}
// 	output = append(output, nextOutput...)

// 	return output, true, nil
// }

// func parseBracket(parts ...func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error)) func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {

// 	return func(bla *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
// 		output := []Node{}

// 		for _, parseFunc := range parts {
// 			args, isOk, err := parseFunc(bla)
// 			if err != nil {
// 				return output, false, err
// 			}
// 			if !isOk {
// 				return output, false, nil
// 			}
// 			output = append(output, args...)
// 		}

// 		return output, true, nil
// 	}
// }

// func parseOr(parts ...func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error)) func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
// 	return func(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
// 		for _, part := range parts {
// 			nodes, ok, err := part(buf)
// 			if err != nil {
// 				return nil, false, err
// 			}
// 			if ok {
// 				return nodes, true, nil
// 			}
// 		}
// 		return nil, false, nil
// 	}
// }

// // -----------------------------------
// // CODE INDEPENDANT OF GRAMMAR END
// // -----------------------------------

// type Grammar_expression struct {
// 	Arguments []Node
// }
// type Grammar_comma struct {
// 	Arguments []Node
// }
// type Grammar_equality struct {
// 	Arguments []Node
// }
// type Grammar_primary struct {
// 	Arguments []Node
// }
// type Grammar_comparison struct {
// 	Arguments []Node
// }

// func (non_terminal *Grammar_expression) Evaluate() *Value {
// 	return nil
// }
// func (non_terminal *Grammar_comma) Evaluate() *Value {
// 	return nil
// }
// func (non_terminal *Grammar_equality) Evaluate() *Value {
// 	return nil
// }
// func (non_terminal *Grammar_primary) Evaluate() *Value {
// 	return nil
// }
// func (non_terminal *Grammar_comparison) Evaluate() *Value {
// 	return nil
// }

// func Parse_expression(buf *lexer.BufferedLexicalAnalyzer) (Node, bool, error) {
// 	output := Grammar_expression{}

// 	arg, isOk, err := Parse_comma(buf)
// 	if err != nil {
// 		return &output, false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg)
// 	} else {
// 		return &output, false, nil
// 	}

// 	return &output, true, nil
// }
// func Parse_comma(buf *lexer.BufferedLexicalAnalyzer) (Node, bool, error) {
// 	output := Grammar_comma{}

// 	arg1, isOk, err := Parse_equality(buf)
// 	if err != nil {
// 		return &output, false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg1)
// 	} else {
// 		return &output, false, nil
// 	}

// 	for {
// 		arg2, isOk, err := match(buf, dfa.COMMA)
// 		if err != nil {
// 			return &output, false, err
// 		}
// 		if isOk {
// 			output.Arguments = append(output.Arguments, arg2)
// 		} else {
// 			break
// 		}

// 		arg3, isOk, err := Parse_equality(buf)
// 		if err != nil {
// 			return &output, false, err
// 		}
// 		if isOk {
// 			output.Arguments = append(output.Arguments, arg3)
// 		} else {
// 			return &output, false, nil
// 			// TODO: MAybe I can add custom error message at each !isOk for better syntax error display to user
// 		}
// 	}

// 	return &output, true, nil
// }

// func Parse_equality(buf *lexer.BufferedLexicalAnalyzer) (Node, bool, error) {
// 	// output := grammar_equality{}
// 	// return &output, nil
// 	return Parse_primary(buf)
// }
// func Parse_primary(buf *lexer.BufferedLexicalAnalyzer) (Node, bool, error) {
// 	output := Grammar_primary{}

// 	arg1, isOk, err := match(buf, dfa.IDENTIFIER)
// 	if err != nil {
// 		return &output, false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg1)
// 		return &output, true, nil
// 	} else {

// 	}

// 	arg2, isOk, err := match(buf, dfa.NUMBER)
// 	if err != nil {
// 		return &output, false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg2)
// 		return &output, true, nil
// 	} else {

// 	}

// 	arg3, isOk, err := match(buf, dfa.STRING)
// 	if err != nil {
// 		return &output, false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg3)
// 		return &output, true, nil
// 	} else {

// 	}

// 	arg4, isOk, err := match(buf, dfa.TRUE)
// 	if err != nil {
// 		return &output, false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg4)
// 		return &output, true, nil
// 	} else {

// 	}

// 	arg5, isOk, err := match(buf, dfa.FALSE)
// 	if err != nil {
// 		return &output, false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg5)
// 		return &output, true, nil
// 	} else {

// 	}

// 	arg6, isOk, err := match(buf, dfa.NIL)
// 	if err != nil {
// 		return &output, false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg6)
// 		return &output, true, nil
// 	} else {

// 	}

// 	arg7, isOk, err := match(buf, dfa.LEFT_PAREN)
// 	if err != nil {
// 		return &output, false, err
// 	}
// 	if !isOk {
// 		arg8, isOk, err := Parse_expression(buf)
// 		if err != nil {
// 			return &output, err
// 		}
// 		if !isOk {

// 		}
// 		output.Arguments = append(output.Arguments, arg8)

// 		arg9, isOk, err := match(buf, dfa.RIGHT_PAREN)
// 		if err != nil {
// 			return &output, err
// 		}
// 		if !isOk {

// 		}
// 		if arg9 != nil {
// 			return &output, nil
// 		}
// 	}

// 	return &output, fmt.Errorf("no match for non-terminal `primary`")
// }

// func Parse_comparison(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {

// 	output := Grammar_comparison{}

// 	expr := parseBracket(Parse_equality, parseZeroOrMoreWrapper(
// 		parseBracket(parseZeroOrMore, parseBracket),
// 	))

// 	wrappedExpr := func(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
// 		return parseBracket(buf, expr)
// 	}

// 	return []Node{&output}, false, nil
// }

// func Parse_term(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
// 	return []Node{}, true, nil
// }

// func comparision_term2(buf *lexer.BufferedLexicalAnalyzer, output *Grammar_comparison) (bool, error) {

// 	isOk, err := comparision_term2_term1(buf) // Only made when term type is `star`, `plus` or `bracket`
// 	if err != nil {
// 		return false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg1)
// 	}

// 	arg1, isOk, err := Parse_primary(buf)
// 	if err != nil {
// 		return false, err
// 	}
// 	if isOk {
// 		output.Arguments = append(output.Arguments, arg1)
// 	} else {
// 		return false, nil
// 	}

// 	return false, nil
// }

package parser

import (
	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/lexer/dfa"
)

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

func matchToken(t dfa.TokenType) func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	return func(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
		tok, err := buf.CurrentTokenWithoutConsume()
		if err != nil {
			return nil, false, err
		}
		if tok.TypeOfToken == t {
			buf.ConsumeOneToken()
			return []Node{&Literal{tok}}, true, nil
		}
		return nil, false, nil
		// fmt.Errorf("Unexpected token '%v' found at line %d, offset %d. Expected token '%v'", string(tok.TypeOfToken), tok.Line, tok.Offset, string(t))
	}
}

func sequence(parts ...func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error)) func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	return func(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
		output := []Node{}
		for _, part := range parts {
			nodes, ok, err := part(buf)
			if err != nil || !ok {
				return nil, false, err
			}
			output = append(output, nodes...)
		}
		return output, true, nil
	}
}

func choice(parts ...func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error)) func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	return func(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
		for _, part := range parts {
			nodes, ok, err := part(buf)
			if err != nil {
				return nil, false, err
			}
			if ok {
				return nodes, true, nil
			}
		}
		return nil, false, nil
	}
}

func zeroOrMore(part func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error)) func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	return func(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
		output := []Node{}
		for {
			nodes, ok, err := part(buf)
			if err != nil || !ok {
				break
			}
			output = append(output, nodes...)
		}
		return output, true, nil
	}
}

// -------------------- PARSE RULES --------------------

func Parse_expression(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	output := Grammar_expression{}

	args, ok, err := Parse_comma(buf)

	output.Arguments = args
	if err != nil || !ok {
		return []Node{&output}, false, err
	}
	return []Node{&output}, false, nil
}

func Parse_comma(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
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
		return []Node{&output}, false, err
	}
	return []Node{&output}, false, nil
}

func Parse_equality(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
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
		return []Node{&output}, false, err
	}
	return []Node{&output}, false, nil
}

func Parse_comparison(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
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
		return []Node{&output}, false, err
	}
	return []Node{&output}, false, nil
}

func Parse_term(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	output := Grammar_term{}

	return sequence(
		Parse_factor,
		zeroOrMore(sequence(
			choice(matchToken(dfa.MINUS), matchToken(dfa.PLUS)),
			Parse_factor,
		)),
	)(buf)

	output.Arguments = args
	if err != nil || !ok {
		return []Node{&output}, false, err
	}
	return []Node{&output}, false, nil
}

func Parse_factor(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	return sequence(
		Parse_unary,
		zeroOrMore(sequence(
			choice(matchToken(dfa.SLASH), matchToken(dfa.STAR)),
			Parse_unary,
		)),
	)(buf)
}

func Parse_unary(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	return choice(
		sequence(choice(matchToken(dfa.BANG), matchToken(dfa.MINUS)), Parse_unary),
		Parse_primary,
	)(buf)
}

func Parse_primary(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	return choice(
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
}

//

type Grammar_expression struct {
	Arguments []Node
}
type Grammar_comma struct {
	Arguments []Node
}
type Grammar_equality struct {
	Arguments []Node
}
type Grammar_primary struct {
	Arguments []Node
}
type Grammar_comparison struct {
	Arguments []Node
}

func (non_terminal *Grammar_expression) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_comma) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_equality) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_primary) Evaluate() *Value {
	return nil
}
func (non_terminal *Grammar_comparison) Evaluate() *Value {
	return nil
}
