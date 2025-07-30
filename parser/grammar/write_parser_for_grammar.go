package grammar

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/VirajAgarwal1/lox/errorhandler"
)

var TokenStringToType = map[string]string{
	// Literals
	"EOF":        "dfa.EOF",
	"IDENTIFIER": "dfa.IDENTIFIER",
	"STRING":     "dfa.STRING",
	"NUMBER":     "dfa.NUMBER",
	"COMMENT":    "dfa.COMMENT",

	// Single-char tokens
	" ":  "dfa.WHITESPACE",
	"\n": "dfa.NEWLINE",
	"(":  "dfa.LEFT_PAREN",
	")":  "dfa.RIGHT_PAREN",
	"{":  "dfa.LEFT_BRACE",
	"}":  "dfa.RIGHT_BRACE",
	",":  "dfa.COMMA",
	".":  "dfa.DOT",
	"-":  "dfa.MINUS",
	"+":  "dfa.PLUS",
	";":  "dfa.SEMICOLON",
	"/":  "dfa.SLASH",
	"*":  "dfa.STAR",

	// One-or-two char tokens
	"!":  "dfa.BANG",
	"!=": "dfa.BANG_EQUAL",
	"=":  "dfa.EQUAL",
	"==": "dfa.EQUAL_EQUAL",
	">":  "dfa.GREATER",
	">=": "dfa.GREATER_EQUAL",
	"<":  "dfa.LESS",
	"<=": "dfa.LESS_EQUAL",

	// Keywords
	"and":    "dfa.AND",
	"class":  "dfa.CLASS",
	"else":   "dfa.ELSE",
	"false":  "dfa.FALSE",
	"fun":    "dfa.FUN",
	"for":    "dfa.FOR",
	"if":     "dfa.IF",
	"nil":    "dfa.NIL",
	"or":     "dfa.OR",
	"print":  "dfa.PRINT",
	"return": "dfa.RETURN",
	"super":  "dfa.SUPER",
	"this":   "dfa.THIS",
	"true":   "dfa.TRUE",
	"var":    "dfa.VAR",
	"while":  "dfa.WHILE",
}

func writeStructsForNonTerminals(writer *bufio.Writer, processedGrammar map[non_terminal]([]generic_grammar_term)) error {

	getStringForNonTerminal := func(symbol string) string {
		output := `type Grammar_` + symbol + ` struct {
	Arguments []Node
}
`
		return output
	}

	for nonTerminalSymbol := range processedGrammar {
		_, err := writer.WriteString(getStringForNonTerminal(nonTerminalSymbol.name))
		if err != nil {
			return errorhandler.RetErr("", err)
		}
	}

	// Add a newline for cleanliness
	_, err := writer.WriteString("\n")
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	return nil
}

func writeEvaluateMethodsForNonTerminals(writer *bufio.Writer, processedGrammar map[non_terminal]([]generic_grammar_term)) error {

	getStringForNonTerminal := func(symbol string) string {
		output := `func (non_terminal *Grammar_` + symbol + `) Evaluate() *Value {
	return nil
}
`
		return output
	}

	for nonTerminalSymbol := range processedGrammar {
		_, err := writer.WriteString(getStringForNonTerminal(nonTerminalSymbol.name))
		if err != nil {
			return errorhandler.RetErr("", err)
		}
	}

	// Add a newline for cleanliness
	_, err := writer.WriteString("\n")
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	return nil
}

func writeParseFunctionsForNonTerminals(writer *bufio.Writer, processedGrammar map[non_terminal]([]generic_grammar_term)) error {

	getStringForNonTerminal := func(nonTerminalSymbol non_terminal, processedGrammar map[non_terminal]([]generic_grammar_term)) string {
		output := `func Parse_` + nonTerminalSymbol.name + `(buf *lexer.BufferedLexicalAnalyzer) (Node, bool, error) {
	output := Grammar_` + nonTerminalSymbol.name + `{}

` + indentLines(generateDescriptionCode(processedGrammar[nonTerminalSymbol]), 1) + `

	return &output, true, nil
}
`
		return output
	}

	for nonTerminalSymbol := range processedGrammar {
		_, err := writer.WriteString(getStringForNonTerminal(nonTerminalSymbol, processedGrammar))
		if err != nil {
			return errorhandler.RetErr("", err)
		}
	}

	// Add a newline for cleanliness
	_, err := writer.WriteString("\n")
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	return nil
}

// -----------------------------------------------------------------------------------
// Functions which help in writing code in recursive fashion for the grammar
// -----------------------------------------------------------------------------------

func detectOrInDescription(description []generic_grammar_term) []uint32 {
	out := make([]uint32, len(description)/4)
	for i, term := range description {
		if term.get_grammar_term_type() == "or" {
			out = append(out, uint32(i))
		}
	}
	return out
}

func indentLines(input string, indentLevel int) string {
	indent := ""
	for i := 0; i < indentLevel; i++ {
		indent += "\t"
	}

	lines := strings.Split(input, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}

	return strings.Join(lines, "\n")
}

func generateMatchCode(term terminal, isNotOkayIfCodeSegment string, isOkayIfCodeSegment string, outputStartIndex *uint32) string {

	parsedTermNum := strconv.FormatUint(uint64(*outputStartIndex), 10)
	(*outputStartIndex)++

	output := ""

	output += "parsedTerm" + parsedTermNum + ", isOk, err := match(buf, " + TokenStringToType[string(term.content)] + ")"

	output += `
if err != nil {
	return &output, err
}
if !isOk {
	` + isNotOkayIfCodeSegment + `
}
output.Arguments = append(output.Arguments, parsedTerm` + strconv.FormatUint(uint64(*outputStartIndex), 10) + `)
` + isOkayIfCodeSegment + `
`
	return output
}

func generateParseReturnCode(term non_terminal, outputStartIndex *uint32, errorCodeSegment string, matchCodeSegment string) string {

	parsedTermNum := strconv.FormatUint(uint64(*outputStartIndex), 10)
	(*outputStartIndex)++

	output := "parsedTerm" + parsedTermNum + ", err := Parse_" + term.name + "(buf)"
	output += `
if err != nil {
	` + errorCodeSegment + `
} else {
	output.Arguments = append(output.Arguments, parsedTerm` + parsedTermNum + `)
	` + matchCodeSegment + `
}
`

	return output
}

func generateStarCode(star_term star, outputStartIndex *uint32) string {
	output := "for {\n"

	if star_term.content.get_grammar_term_type() == "terminal" {
		output += indentLines(
			generateMatchCode(
				*star_term.content.(*terminal), "break", "", outputStartIndex,
			),
			1,
		)
	}
	if star_term.content.get_grammar_term_type() == "non_terminal" {
		output += indentLines(
			generateParseReturnCode(
				*star_term.content.(*non_terminal), outputStartIndex, "break",
			),
			1,
		)
	}
	if star_term.content.get_grammar_term_type() == "star" {
		output = generateStarCode(*star_term.content.(*star), outputStartIndex)
		return output
	}
	if star_term.content.get_grammar_term_type() == "plus" {
		bypassPlusOperator := star{}
		bypassPlusOperator.content = star_term.content.(*plus).content
		output = generateStarCode(bypassPlusOperator, outputStartIndex)
		return output
	}
	if star_term.content.get_grammar_term_type() == "bracket" {
		output += indentLines(
			generateBracketCode(*star_term.content.(*bracket), outputStartIndex, true),
			1,
		)
	}

	output += "}\n"
	return output
}

func generateBracketCode(bracket_term bracket, outputStartIndex *uint32, isPartOfStarOrPLus bool) string {

	// Detect if there is an `or` in the description
	orPositions := detectOrInDescription(bracket_term.contents)
	if len(orPositions) == 0 {
		// All the terms in the description will concatenate

	}

	output := ""

	for i := 0; i < len(bracket_term.contents); i++ {
		if i == 0 && isPartOfStarOrPLus {
			if bracket_term.contents[i].get_grammar_term_type() == "terminal" {
				output += indentLines(
					generateMatchCode(*bracket_term.contents[i].(*terminal), false, "break", outputStartIndex),
					1,
				)
			}
		}
	}
	return output
}

func generateDescriptionCode(description []generic_grammar_term, outputStartIndex *uint32, isPartOfBiggerOr bool, returnCodeBlock string) string {
	// This will be a recursive function
	if len(description) == 0 {
		return ""
	}
	output := ""

	// Detect if there is an `or` in the description
	orPositions := detectOrInDescription(description)
	if len(orPositions) != 0 {
		// All the terms in the description will concatenate
		for i := 0; i < len(orPositions); i++ {
			if i == 0 {
				output += generateDescriptionCode(
					description[0:orPositions[i]],
					outputStartIndex,
					true,
					returnCodeBlock,
				)
				continue
			}
			output += generateDescriptionCode(
				description[orPositions[i-1]:orPositions[i]],
				outputStartIndex,
				true,
				returnCodeBlock,
			)
		}

		return output
	}

	for i := 0; i < len(description); i++ {
		term := description[i]

		if term.get_grammar_term_type() == "terminal" {
			output += generateMatchCode(*term.(*terminal), returnCodeBlock, "", outputStartIndex)
		}
	}

	return output
}

// -----------------------------------------------------------------------------------
// Main Function which orchestrates the dance of all the other functions in this file
// -----------------------------------------------------------------------------------

func generateGrammarOutput(writer *bufio.Writer, processedGrammar map[non_terminal]([]generic_grammar_term)) error {

	// Writing function and strcuts which are independant of the grammar
	_, err := writer.WriteString(
		`// -----------------------------------
// CODE INDEPENDANT OF GRAMMAR START
// -----------------------------------

package parser

import (
	"fmt"
	"io"

	"github.com/VirajAgarwal1/lox/errorhandler"
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
func match(buf_scanner *lexer.BufferedLexicalAnalyzer, token dfa.TokenType) (Node, bool, error) {
	t, err := buf_scanner.CurrentTokenWithoutConsume()
	if err != nil && err != io.EOF {
		return nil, false, errorhandler.RetErr("", err)
	}
	if t.TypeOfToken == token {
		buf_scanner.ConsumeOneToken()
		literal_value := Literal{
			Value: t,
		}
		return &literal_value, true, nil
	}
	return nil, false, nil
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

// -----------------------------------
// CODE INDEPENDANT OF GRAMMAR END
// -----------------------------------

`,
	)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	// Write the structs for each non-terminal symbol
	err = writeStructsForNonTerminals(writer, processedGrammar)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	// Write the Evaluate method on struct for each non-terminal symbol
	err = writeEvaluateMethodsForNonTerminals(writer, processedGrammar)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	// Write the Parsing functions for each non-terminal symbol
	err = writeParseFunctionsForNonTerminals(writer, processedGrammar)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	return nil
}

// -----------------------------------------------------------------------------------
// Main Function which creates/edits the file
// -----------------------------------------------------------------------------------

func GenerateGrammarParserFile(processedGrammar map[non_terminal]([]generic_grammar_term), filePath string) error {
	// Open the file with:
	// - O_CREATE: create if not exists
	// - O_TRUNC: clear the file first
	// - O_WRONLY: open in write-only mode
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return errorhandler.RetErr("", err)
	}
	defer file.Close()

	// Wrap in a buffered writer
	writer := bufio.NewWriter(file)

	// Take the processed grammar and output its parser to the file
	err = generateGrammarOutput(writer, processedGrammar)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	// Flush buffered data to disk
	err = writer.Flush()
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	return nil
}
