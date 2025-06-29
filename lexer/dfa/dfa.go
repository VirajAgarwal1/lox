package dfa

type DFA interface {
	Step(input rune) DfaReturn
	Reset()
}

const (
	INVALID DfaReturn = iota
	INTERMEDIATE
	VALID
)

type DfaReturn byte

func (s *DfaReturn) ToString() string {
	if int(*s) == 0 {
		return "INVALID"
	}
	if int(*s) == 1 {
		return "INTERMEDIATE"
	}
	if int(*s) == 2 {
		return "VALID"
	}
	return ""
}

func (s *DfaReturn) IsValid() bool {
	return *s == VALID
}

func (s *DfaReturn) IsIntermediate() bool {
	return *s == INTERMEDIATE
}

func (s *DfaReturn) IsInvalid() bool {
	return *s == INVALID
}

type TokenType string

// These are all lexemes which will identified by the scanner.
const (
	EOF TokenType = "EOF"

	// Literals
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"
	COMMENT    TokenType = "COMMENT"

	// Single-char tokens
	WHITESPACE  TokenType = " "
	NEWLINE     TokenType = "\n"
	LEFT_PAREN  TokenType = "("
	RIGHT_PAREN TokenType = ")"
	LEFT_BRACE  TokenType = "{"
	RIGHT_BRACE TokenType = "}"
	COMMA       TokenType = ","
	DOT         TokenType = "."
	MINUS       TokenType = "-"
	PLUS        TokenType = "+"
	SEMICOLON   TokenType = ";"
	SLASH       TokenType = "/"
	STAR        TokenType = "*"

	// One-or-two char tokens
	BANG          TokenType = "!"
	BANG_EQUAL    TokenType = "!="
	EQUAL         TokenType = "="
	EQUAL_EQUAL   TokenType = "=="
	GREATER       TokenType = ">"
	GREATER_EQUAL TokenType = ">="
	LESS          TokenType = "<"
	LESS_EQUAL    TokenType = "<="

	// Keywords.
	AND    TokenType = "and"
	CLASS  TokenType = "class"
	ELSE   TokenType = "else"
	FALSE  TokenType = "false"
	FUN    TokenType = "fun"
	FOR    TokenType = "for"
	IF     TokenType = "if"
	NIL    TokenType = "nil"
	OR     TokenType = "or"
	PRINT  TokenType = "print"
	RETURN TokenType = "return"
	SUPER  TokenType = "super"
	THIS   TokenType = "this"
	TRUE   TokenType = "true"
	VAR    TokenType = "var"
	WHILE  TokenType = "while"
)

var TokensList = []TokenType{
	LEFT_PAREN,
	RIGHT_PAREN,
	LEFT_BRACE,
	RIGHT_BRACE,
	COMMA,
	DOT,
	MINUS,
	PLUS,
	SEMICOLON,
	SLASH,
	STAR,
	BANG,
	BANG_EQUAL,
	EQUAL,
	EQUAL_EQUAL,
	GREATER,
	GREATER_EQUAL,
	LESS,
	LESS_EQUAL,
	AND,
	CLASS,
	ELSE,
	FALSE,
	FUN,
	FOR,
	IF,
	NIL,
	OR,
	PRINT,
	RETURN,
	SUPER,
	THIS,
	TRUE,
	VAR,
	WHILE,
	EOF,
	IDENTIFIER,
	STRING,
	NUMBER,
}

func GenerateDFAs() map[TokenType]DFA {

	output := make(map[TokenType]DFA)
	for _, k := range TokensList {
		if k == EOF {
			continue // EOF is not a string, so we skip it
		}
		if k == IDENTIFIER {
			dfa := &IdentifierDFA{}
			dfa.Initialize()
			output[k] = dfa
			continue
		}
		if k == STRING {
			dfa := &StringDFA{}
			dfa.Initialize()
			output[k] = dfa
			continue
		}
		if k == COMMENT {
			dfa := &CommentDFA{}
			dfa.Initialize()
			output[k] = dfa
			continue
		}
		if k == NUMBER {
			dfa := &NumberDFA{}
			dfa.Initialize()
			output[k] = dfa
			continue
		}

		dfa := &InputStringDFA{}
		dfa.Initialize(string(k))
		output[k] = dfa
	}

	return output
}
