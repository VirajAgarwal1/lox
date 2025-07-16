package parser

import (
	lexer "github.com/VirajAgarwal1/lox/lexer"
)

type Expression struct {
	arg1 *Comma
}
type Comma struct {
	arg1 *Equality
	arg2 *Comparision
}
type Equality struct {
}
type Comparision struct {
}

func (nonTerminal *Expression) parse(scanner *lexer.LexicalAnalyzer) bool {
	if !nonTerminal.arg1.parse(scanner) {
		return false
	}
	return true
}
func (nonTerminal *Comma) parse(scanner *lexer.LexicalAnalyzer) bool {
	if !nonTerminal.arg1.parse(scanner) {
		return false
	}
	return true
}
func (nonTerminal *Equality) parse(scanner *lexer.LexicalAnalyzer) bool {
	return true
}
func (nonTerminal *Comparision) parse(scanner *lexer.LexicalAnalyzer) bool {
	return true
}
