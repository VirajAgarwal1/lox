package grammar

import (
	"bufio"
	"os"
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

func WriteStructsForNonTerminals(writer *bufio.Writer, processedGrammar map[Non_terminal]([]Generic_grammar_term)) error {

	getStringForNonTerminal := func(symbol string) string {
		output := `type Grammar_` + symbol + ` struct {
	Arguments []Node
}
`
		return output
	}

	for nonTerminalSymbol := range processedGrammar {
		_, err := writer.WriteString(getStringForNonTerminal(nonTerminalSymbol.Name))
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

func WriteEvaluateMethodsForNonTerminals(writer *bufio.Writer, processedGrammar map[Non_terminal]([]Generic_grammar_term)) error {

	getStringForNonTerminal := func(symbol string) string {
		output := `func (non_terminal *Grammar_` + symbol + `) Evaluate() *Value {
	return nil
}
`
		return output
	}

	for nonTerminalSymbol := range processedGrammar {
		_, err := writer.WriteString(getStringForNonTerminal(nonTerminalSymbol.Name))
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

func WriteParseFunctionsForNonTerminals(writer *bufio.Writer, processedGrammar map[Non_terminal]([]Generic_grammar_term)) error {

	getStringForNonTerminal := func(nonTerminalSymbol Non_terminal, processedGrammar map[Non_terminal]([]Generic_grammar_term)) string {
		output := `func Parse_` + nonTerminalSymbol.Name + `(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	output := Grammar_` + nonTerminalSymbol.Name + `{}
	
	args, ok, err := ` + GenerateDescriptionCode(processedGrammar[nonTerminalSymbol], "(buf)") + `
	
	output.Arguments = args
	if err != nil || !ok {
		return nil, false, err
	}
	return []Node{&output}, true, nil
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

func DetectOrInDescription(description []Generic_grammar_term) []uint32 {
	out := make([]uint32, 0, len(description)/4)
	for i, term := range description {
		if term.Get_grammar_term_type() == "or" {
			out = append(out, uint32(i))
		}
	}
	return out
}

func IndentLines(input string, indentLevel int) string {
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

func GetStringGeneratorForTerm(term Generic_grammar_term, endString string) string {
	if term.Get_grammar_term_type() == "terminal" {
		return GenerateMatchCode(*term.(*Terminal), endString)
	}
	if term.Get_grammar_term_type() == "non_terminal" {
		return GenerateParseReturnCode(*term.(*Non_terminal), endString)
	}
	if term.Get_grammar_term_type() == "star" {
		return GenerateStarCode(*term.(*Star), endString)
	}
	if term.Get_grammar_term_type() == "plus" {
		return GeneratePlusCode(*term.(*Plus), endString)
	}
	if term.Get_grammar_term_type() == "bracket" {
		return GenerateBracketCode(*term.(*Bracket), endString)
	}
	return ""
}

func GenerateMatchCode(term Terminal, endString string) string {
	return "matchToken(" + TokenStringToType[string(term.Content)] + ")" + endString
}
func GenerateParseReturnCode(term Non_terminal, endString string) string {
	return "Parse_" + term.Name + endString
}
func GenerateStarCode(star_term Star, endString string) string {
	return "zeroOrMore(\n" + GenerateDescriptionCode([]Generic_grammar_term{star_term.Content}, ",\n") + ")" + endString
}
func GeneratePlusCode(plus_term Plus, endString string) string {
	return "oneOrMore(\n" + GenerateDescriptionCode([]Generic_grammar_term{plus_term.Content}, ",\n") + ")" + endString
}
func GenerateBracketCode(bracket_term Bracket, endString string) string {
	if len(bracket_term.Contents) == 0 {
		return ""
	}
	return GenerateDescriptionCode(bracket_term.Contents, endString)
}

func GenerateDescriptionCode(description []Generic_grammar_term, endString string) string {
	// This will be a recursive function
	if len(description) == 0 {
		return ""
	}

	orPositions := DetectOrInDescription(description)
	if len(orPositions) != 0 {
		output := "choice(\n"

		start := uint32(0)
		end := orPositions[0]
		output += IndentLines(
			GenerateDescriptionCode(
				description[start:end],
				",\n",
			),
			1,
		)
		for i := range orPositions {
			start = orPositions[i] + 1
			if i == len(orPositions)-1 {
				end = uint32(len(description))
			} else {
				end = orPositions[i+1]
			}
			output += IndentLines(
				GenerateDescriptionCode(
					description[start:end],
					",\n",
				),
				1,
			)
		}

		output += ")" + endString
		return output
	}

	if len(description) == 1 {
		return GetStringGeneratorForTerm(description[0], endString)
	}

	output := "sequence(\n"
	for i := 0; i < len(description); i++ {
		output += GetStringGeneratorForTerm(description[i], ",\n")
	}
	output += ")" + endString
	return output
}

// -----------------------------------------------------------------------------------
// Main Function which orchestrates the dance of all the other functions in this file
// -----------------------------------------------------------------------------------

func GenerateGrammarOutput(writer *bufio.Writer, processedGrammar map[Non_terminal]([]Generic_grammar_term)) error {

	// Writing function and strcuts which are independant of the grammar
	_, err := writer.WriteString(
		`// -----------------------------------
// CODE INDEPENDANT OF GRAMMAR START
// -----------------------------------

package parser

import (
	"io"

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
		if err != nil && err != io.EOF {
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
func oneOrMore(part func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error)) func(*lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
	return func(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error) {
		output := []Node{}
		nodes, ok, err := part(buf)
		if err != nil || !ok {
			return nil, false, err
		}
		output = append(output, nodes...)

		for {
			nodes, ok, err = part(buf)
			if err != nil || !ok {
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

`,
	)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	// Write the structs for each non-terminal symbol
	err = WriteStructsForNonTerminals(writer, processedGrammar)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	// Write the Evaluate method on struct for each non-terminal symbol
	err = WriteEvaluateMethodsForNonTerminals(writer, processedGrammar)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	// Write the Parsing functions for each non-terminal symbol
	err = WriteParseFunctionsForNonTerminals(writer, processedGrammar)
	if err != nil {
		return errorhandler.RetErr("", err)
	}

	return nil
}

// -----------------------------------------------------------------------------------
// Main Function which creates/edits the file
// -----------------------------------------------------------------------------------

func GenerateGrammarParserFile(processedGrammar map[Non_terminal]([]Generic_grammar_term), filePath string) error {
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
	err = GenerateGrammarOutput(writer, processedGrammar)
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
