package lexer_demos

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/lexer/dfa"
)

// ANSI color codes for pretty output
const (
	Purple = "\033[35m"
	Gray   = "\033[90m"
)

// TokenTypeName returns a human-readable name for a TokenType
func TokenTypeName(t dfa.TokenType) string {
	switch t {
	case dfa.IDENTIFIER:
		return "IDENTIFIER"
	case dfa.NUMBER:
		return "NUMBER"
	case dfa.STRING:
		return "STRING"
	case dfa.IF:
		return "IF"
	case dfa.VAR:
		return "VAR"
	case dfa.PRINT:
		return "PRINT"
	case dfa.EQUAL:
		return "EQUAL"
	case dfa.EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case dfa.LEFT_PAREN:
		return "LEFT_PAREN"
	case dfa.RIGHT_PAREN:
		return "RIGHT_PAREN"
	case dfa.LEFT_BRACE:
		return "LEFT_BRACE"
	case dfa.RIGHT_BRACE:
		return "RIGHT_BRACE"
	case dfa.SEMICOLON:
		return "SEMICOLON"
	case dfa.EOF:
		return "EOF"
	// ... add all other token types
	default:
		return fmt.Sprintf("UNKNOWN(%v)", t)
	}
}

func ScannerDemo() {
	// Section 1: Header
	fmt.Printf("\n%s%s📜 Full Lexical Analysis Demo %s\n", Bold+Cyan, "", Reset)
	fmt.Printf("%sWelcome to the scanner showcase! This demo will tokenize source code using your full LexicalAnalyzer implementation.%s\n", "", Reset)
	fmt.Printf("%sYou’ll see how source code is split into tokens with types, positions, and lexemes.%s\n", "", Reset)

	// Section 2: Input Code
	sourceCode := `
var count = 0;
while (count < 3) {
  if (count == 1) {
    print "middle";
  } else {
    print count;
  }
  count = count + 1;
}
`

	fmt.Printf("\n%s%s🧾 Source Code:%s\n", Bold+Yellow, "", Reset)

	lines := strings.Split(sourceCode, "\n")
	for i, line := range lines {
		fmt.Printf("%s%2d | %s%s\n", Gray, i, line, Reset)
	}

	// Section 3: Initialize scanner
	reader := bufio.NewReader(strings.NewReader(sourceCode))
	scanner := lexer.LexicalAnalyzer{}
	scanner.Initialize(reader)

	fmt.Printf("\n%s%s🔬 Beginning Tokenization...%s\n\n", Bold+Green, "", Reset)
	time.Sleep(400 * time.Millisecond)

	// Section 4: Read tokens one by one
	for {
		token, err := scanner.ReadToken()

		if err != nil && err != io.EOF {
			fmt.Printf("%s❗ Error: %s%s\n", Red, err.Error(), Reset)
			break
		}

		// Skip whitespace-like unknowns (quick workaround)
		if (token.TypeOfToken == dfa.WHITESPACE) || (token.TypeOfToken == dfa.NEWLINE) {
			continue
		}

		printStyledToken(token)

		if err != nil {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}

	// Section 5: Done
	fmt.Printf("\n%s%s✅ Done! All tokens emitted.%s\n", Bold+Green, "", Reset)
	fmt.Printf("%sTry modifying the source input to explore other token types and error conditions.%s\n", "", Reset)
}

// Prints each token nicely with token type, lexeme, and position
func printStyledToken(token lexer.Token) {
	color := tokenColor(token.TypeOfToken)

	fmt.Printf("%s🧩 Token:%s  %s%-14s%s  at Line %2d, Offset %2d  →  `%s`\n",
		Cyan, Reset,
		color, TokenTypeName(token.TypeOfToken), Reset,
		token.Line, token.Offset,
		visualLexeme(token.Lexemme),
	)
}

// Converts invisible lexeme characters to visible representation
func visualLexeme(runes []rune) string {
	escaped := ""
	for _, r := range runes {
		switch r {
		case ' ':
			escaped += "␣"
		case '\n':
			escaped += "\\n"
		case '\r':
			escaped += "\\r"
		case '\t':
			escaped += "\\t"
		default:
			if r < 32 || r == 127 {
				escaped += fmt.Sprintf("\\x%02x", r)
			} else {
				escaped += string(r)
			}
		}
	}
	// Wrap in quotes
	return strconv.Quote(escaped)
}

// Assign colors based on token type categories (customize more if needed)
func tokenColor(tokenType dfa.TokenType) string {
	switch tokenType {
	case dfa.IDENTIFIER:
		return Yellow
	case dfa.NUMBER:
		return Blue
	case dfa.STRING:
		return Purple
	case dfa.IF, dfa.WHILE, dfa.FOR, dfa.VAR, dfa.FUN, dfa.RETURN, dfa.PRINT:
		return Cyan
	case dfa.EQUAL, dfa.EQUAL_EQUAL, dfa.LEFT_PAREN, dfa.RIGHT_PAREN, dfa.LEFT_BRACE, dfa.RIGHT_BRACE, dfa.SEMICOLON:
		return Gray
	case dfa.EOF:
		return Green
	default:
		return Red
	}
}
