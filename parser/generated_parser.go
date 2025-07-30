// // -----------------------------------
// // CODE INDEPENDANT OF GRAMMAR START
// // -----------------------------------

// package parser

// import (
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
// func match(buf_scanner *lexer.BufferedLexicalAnalyzer, token dfa.TokenType) (Node, error) {
// 	t, err := buf_scanner.CurrentTokenWithoutConsume()
// 	if err != nil && err != io.EOF {
// 		return nil, errorhandler.RetErr("", err)
// 	}
// 	if t.TypeOfToken == token {
// 		buf_scanner.ConsumeOneToken()
// 		literal_value := Literal{
// 			Value: t,
// 		}
// 		return &literal_value, nil
// 	}
// 	return nil, nil
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

// // -----------------------------------
// // CODE INDEPENDANT OF GRAMMAR END
// // -----------------------------------

// type Grammar_equality struct {
// 	Arguments []Node
// }
// type Grammar_comparison struct {
// 	Arguments []Node
// }
// type Grammar_term struct {
// 	Arguments []Node
// }
// type Grammar_factor struct {
// 	Arguments []Node
// }
// type Grammar_unary struct {
// 	Arguments []Node
// }
// type Grammar_primary struct {
// 	Arguments []Node
// }
// type Grammar_expression struct {
// 	Arguments []Node
// }
// type Grammar_comma struct {
// 	Arguments []Node
// }

// func (non_terminal *Grammar_term) Evaluate() *Value {
// 	return nil
// }
// func (non_terminal *Grammar_factor) Evaluate() *Value {
// 	return nil
// }
// func (non_terminal *Grammar_unary) Evaluate() *Value {
// 	return nil
// }
// func (non_terminal *Grammar_primary) Evaluate() *Value {
// 	return nil
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
// func (non_terminal *Grammar_comparison) Evaluate() *Value {
// 	return nil
// }

// func Parse_expression(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
// 	output := Grammar_expression{}
// 	return &output, nil
// }
// func Parse_comma(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
// 	output := Grammar_comma{}
// 	return &output, nil
// }
// func Parse_equality(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
// 	output := Grammar_equality{}
// 	return &output, nil
// }
// func Parse_comparison(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
// 	output := Grammar_comparison{}
// 	return &output, nil
// }
// func Parse_term(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
// 	output := Grammar_term{}
// 	return &output, nil
// }
// func Parse_factor(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
// 	output := Grammar_factor{}
// 	return &output, nil
// }
// func Parse_unary(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
// 	output := Grammar_unary{}
// 	return &output, nil
// }
// func Parse_primary(buf *lexer.BufferedLexicalAnalyzer) (Node, error) {
// 	output := Grammar_primary{}
// 	return &output, nil
// }
