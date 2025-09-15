package parser_generator

import (
	"strings"

	"github.com/VirajAgarwal1/lox/lexer/dfa"
	gfp "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/grammar_file_parser"
)

type Grammar_element struct {
	IsNonTerminal bool
	Non_term_name string
	Terminal_type dfa.TokenType
}

const Epsilon = dfa.TokenType("epsilon")

var String_to_type_string = map[string]string{
	// Literals
	"EOF":        "dfa.EOF",
	"IDENTIFIER": "dfa.IDENTIFIER",
	"STRING":     "dfa.STRING",
	"NUMBER":     "dfa.NUMBER",
	"COMMENT":    "dfa.COMMENT",

	// Single-char tokens
	" ":   "dfa.WHITESPACE",
	"\\n": "dfa.NEWLINE",
	"(":   "dfa.LEFT_PAREN",
	")":   "dfa.RIGHT_PAREN",
	"{":   "dfa.LEFT_BRACE",
	"}":   "dfa.RIGHT_BRACE",
	",":   "dfa.COMMA",
	".":   "dfa.DOT",
	"-":   "dfa.MINUS",
	"+":   "dfa.PLUS",
	";":   "dfa.SEMICOLON",
	"/":   "dfa.SLASH",
	"*":   "dfa.STAR",

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

var String_to_token = map[string]dfa.TokenType{
	// Literals
	"EOF":        dfa.EOF,
	"IDENTIFIER": dfa.IDENTIFIER,
	"STRING":     dfa.STRING,
	"NUMBER":     dfa.NUMBER,
	"COMMENT":    dfa.COMMENT,

	// Single-char tokens
	" ":   dfa.WHITESPACE,
	"\\n": dfa.NEWLINE,
	"(":   dfa.LEFT_PAREN,
	")":   dfa.RIGHT_PAREN,
	"{":   dfa.LEFT_BRACE,
	"}":   dfa.RIGHT_BRACE,
	",":   dfa.COMMA,
	".":   dfa.DOT,
	"-":   dfa.MINUS,
	"+":   dfa.PLUS,
	";":   dfa.SEMICOLON,
	"/":   dfa.SLASH,
	"*":   dfa.STAR,

	// One-or-two char tokens
	"!":  dfa.BANG,
	"!=": dfa.BANG_EQUAL,
	"=":  dfa.EQUAL,
	"==": dfa.EQUAL_EQUAL,
	">":  dfa.GREATER,
	">=": dfa.GREATER_EQUAL,
	"<":  dfa.LESS,
	"<=": dfa.LESS_EQUAL,

	// Keywords
	"and":    dfa.AND,
	"class":  dfa.CLASS,
	"else":   dfa.ELSE,
	"false":  dfa.FALSE,
	"fun":    dfa.FUN,
	"for":    dfa.FOR,
	"if":     dfa.IF,
	"nil":    dfa.NIL,
	"or":     dfa.OR,
	"print":  dfa.PRINT,
	"return": dfa.RETURN,
	"super":  dfa.SUPER,
	"this":   dfa.THIS,
	"true":   dfa.TRUE,
	"var":    dfa.VAR,
	"while":  dfa.WHILE,
}

func Detect_or_in_sequence(description []gfp.Generic_grammar_term) []uint32 {
	out := make([]uint32, 0, len(description)/4)
	for i, term := range description {
		if term.Get_grammar_term_type() == "or" {
			out = append(out, uint32(i))
		}
	}
	return out
}

func Indent_lines(input string, indentLevel int) string {
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

func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
