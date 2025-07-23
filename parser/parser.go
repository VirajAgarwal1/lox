package parser

import (
	"fmt"

	"github.com/VirajAgarwal1/lox/errorhandler"
	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/lexer/dfa"
)

func match(buf_scanner *lexer.BufferedLexicalAnalyzer, token dfa.TokenType) (Node, error) {
	t, err := buf_scanner.CurrentTokenWithoutConsume()
	if err != nil {
		return nil, errorhandler.RetErr("", err)
	}
	if t.TypeOfToken == token {
		buf_scanner.ConsumeOneToken()
		literal_value := Literal{
			Value: t,
		}
		return &literal_value, nil
	}
	return nil, nil
}

type Value struct {
	LoxType string
	Inner   any
}
type Node interface {
	Evaluate() *Value
}

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
type Literal struct {
	Value *lexer.Token
}

func (non_terminal *Grammar_expression) Evaluate() *Value {
	return &Value{}
}
func (non_terminal *Grammar_comma) Evaluate() *Value {
	return &Value{}
}
func (non_terminal *Grammar_equality) Evaluate() *Value {
	return &Value{}
}
func (non_terminal *Grammar_primary) Evaluate() *Value {
	return &Value{}
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

func Parse_expression(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
	output := Grammar_expression{}

	arg1, err := Parse_comma(buf)
	if err != nil {
		output.Arguments = append(output.Arguments, arg1) // TODO: This is ONLY temporarily
		return &output, err
	}
	output.Arguments = append(output.Arguments, arg1)

	return &output, nil
}
func Parse_comma(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
	output := Grammar_comma{}

	arg1, err := Parse_equality(buf)
	if err != nil {
		return &output, err
	}
	output.Arguments = append(output.Arguments, arg1)

	for {
		tok, err := match(buf, dfa.COMMA)
		if err != nil {
			return &output, err
		}
		if tok == nil {
			break
		}

		arg2, err := Parse_equality(buf)
		if err != nil {
			return &output, err
		}
		output.Arguments = append(output.Arguments, arg2)
	}

	return &output, nil
}
func Parse_equality(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
	// output := grammar_equality{}
	// return &output, nil
	return Parse_primary(buf)
}
func Parse_primary(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
	output := Grammar_primary{}

	arg1, err := match(buf, dfa.IDENTIFIER)
	if err != nil {
		return &output, err
	}
	if arg1 != nil {
		output.Arguments = append(output.Arguments, arg1)
		return &output, nil
	}

	arg2, err := match(buf, dfa.NUMBER)
	if err != nil {
		return &output, err
	}
	if arg2 != nil {
		output.Arguments = append(output.Arguments, arg2)
		return &output, nil
	}

	arg3, err := match(buf, dfa.STRING)
	if err != nil {
		return &output, err
	}
	if arg3 != nil {
		output.Arguments = append(output.Arguments, arg3)
		return &output, nil
	}

	arg4, err := match(buf, dfa.TRUE)
	if err != nil {
		return &output, err
	}
	if arg4 != nil {
		output.Arguments = append(output.Arguments, arg4)
		return &output, nil
	}

	arg5, err := match(buf, dfa.FALSE)
	if err != nil {
		return &output, err
	}
	if arg5 != nil {
		output.Arguments = append(output.Arguments, arg5)
		return &output, nil
	}

	arg6, err := match(buf, dfa.NIL)
	if err != nil {
		return &output, err
	}
	if arg6 != nil {
		output.Arguments = append(output.Arguments, arg6)
		return &output, nil
	}

	arg7, err := match(buf, dfa.LEFT_PAREN)
	if err != nil {
		return &output, err
	}
	if arg7 != nil {
		arg8, err := Parse_expression(buf)
		if err != nil {
			return &output, err
		}
		output.Arguments = append(output.Arguments, arg8)

		arg9, err := match(buf, dfa.RIGHT_PAREN)
		if err != nil {
			return &output, err
		}
		if arg9 != nil {
			return &output, nil
		}
	}

	return &output, fmt.Errorf("no match for non-terminal `primary`")
}
