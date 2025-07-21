package lexer_tests

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"

	lexer "github.com/VirajAgarwal1/lox/lexer"
	dfa "github.com/VirajAgarwal1/lox/lexer/dfa"
)

// ----------------------------
// Helper Functions
// ----------------------------

func createScannerFromString(input string) *lexer.LexicalAnalyzer {
	scanner := &lexer.LexicalAnalyzer{}
	reader := bufio.NewReader(strings.NewReader(input))
	scanner.Initialize(reader)
	return scanner
}

func scanAllTokens(t *testing.T, input string) ([]lexer.Token, []error) {
	scanner := createScannerFromString(input)
	var tokens []lexer.Token
	var errors []error

	for {
		token, err := scanner.ReadToken()
		if err != nil && err != io.EOF {
			errors = append(errors, err)
			continue
		}
		tokens = append(tokens, token)
		if err == io.EOF {
			break
		}
	}
	return tokens, errors
}

// Basic Token Testing
func TestTokenProperties(t *testing.T) {
	token := lexer.Token{}

	// Test SetTokenProperties
	token.SetTokenProperties(dfa.IDENTIFIER, 1, 5, []rune("test"))

	if token.TypeOfToken != dfa.IDENTIFIER {
		t.Errorf("Expected token type %v, got %v", dfa.IDENTIFIER, token.TypeOfToken)
	}
	if token.Line != 1 {
		t.Errorf("Expected line 1, got %d", token.Line)
	}
	if token.Offset != 5 {
		t.Errorf("Expected offset 5, got %d", token.Offset)
	}

	// Test ToString method
	tokenStr := token.ToString()
	expected := "|1|5| [IDENTIFIER]Token -> `test`"
	if tokenStr != expected {
		t.Errorf("Expected ToString() to return '%s', got '%s'", expected, tokenStr)
	}
}

// Basic Lexical Scanner Testing
func TestScannerInitialization(t *testing.T) {
	scanner := &lexer.LexicalAnalyzer{}
	reader := bufio.NewReader(strings.NewReader("test"))

	scanner.Initialize(reader)

	// Test that initialization sets proper default values
	// Note: We can't directly test private fields, but we can test behavior
	token, err := scanner.ReadToken()
	if err != nil && err != io.EOF {
		t.Errorf("Unexpected error during initialization test: %v", err)
	}

	// The token should have line 0 and offset 0 for the first token
	if token.Line != 0 {
		t.Errorf("Expected initial line to be 0, got %d", token.Line)
	}
	if token.Offset != 0 {
		t.Errorf("Expected initial offset to be 0, got %d", token.Offset)
	}
}

func TestEmptyInput(t *testing.T) {
	scanner := createScannerFromString("")

	token, err := scanner.ReadToken()
	if err != io.EOF {
		t.Errorf("Expected EOF error for empty input, got %v", err)
	}
	if token.TypeOfToken != dfa.EOF {
		t.Errorf("Expected EOF token type, got %v", token.TypeOfToken)
	}
}

// Test single character tokens
func TestSingleCharacterTokens(t *testing.T) {
	// Test various single character inputs
	testCases := []struct {
		input    string
		expected dfa.TokenType
	}{
		{"(", dfa.LEFT_PAREN},
		{")", dfa.RIGHT_PAREN},
		{"{", dfa.LEFT_BRACE},
		{"}", dfa.RIGHT_BRACE},
		{";", dfa.SEMICOLON},
		{",", dfa.COMMA},
		{".", dfa.DOT},
		{"+", dfa.PLUS},
		{"-", dfa.MINUS},
		{"*", dfa.STAR},
		{"/", dfa.SLASH},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Single char: %s", tc.input), func(t *testing.T) {
			tokens, errors := scanAllTokens(t, tc.input)

			if len(errors) != 0 {
				t.Errorf("Unexpected errors: %v", errors)
			}
			if len(tokens) != 2 { // Token + EOF
				t.Errorf("Expected 2 tokens (token + EOF), got %d", len(tokens))
			}
			if tokens[0].TypeOfToken != tc.expected {
				t.Errorf("Expected token type %v, got %v", tc.expected, tokens[0].TypeOfToken)
			}
			if tokens[1].TypeOfToken != dfa.EOF {
				t.Errorf("Expected EOF token, got %v", tokens[1].TypeOfToken)
			}
		})
	}
}

func TestTokenRecognition(t *testing.T) {
	testCases := []struct {
		input    string
		expected dfa.TokenType
	}{
		// Single character tokens
		{"(", dfa.LEFT_PAREN},
		{")", dfa.RIGHT_PAREN},
		{"{", dfa.LEFT_BRACE},
		{"}", dfa.RIGHT_BRACE},
		{",", dfa.COMMA},
		{".", dfa.DOT},
		{"-", dfa.MINUS},
		{"+", dfa.PLUS},
		{";", dfa.SEMICOLON},
		{"/", dfa.SLASH},
		{"*", dfa.STAR},

		// One-or-two char tokens
		{"!", dfa.BANG},
		{"!=", dfa.BANG_EQUAL},
		{"=", dfa.EQUAL},
		{"==", dfa.EQUAL_EQUAL},
		{">", dfa.GREATER},
		{">=", dfa.GREATER_EQUAL},
		{"<", dfa.LESS},
		{"<=", dfa.LESS_EQUAL},

		// Keywords
		{"and", dfa.AND},
		{"class", dfa.CLASS},
		{"else", dfa.ELSE},
		{"false", dfa.FALSE},
		{"fun", dfa.FUN},
		{"for", dfa.FOR},
		{"if", dfa.IF},
		{"nil", dfa.NIL},
		{"or", dfa.OR},
		{"print", dfa.PRINT},
		{"return", dfa.RETURN},
		{"super", dfa.SUPER},
		{"this", dfa.THIS},
		{"true", dfa.TRUE},
		{"var", dfa.VAR},
		{"while", dfa.WHILE},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Multi char: %s", tc.input), func(t *testing.T) {
			tokens, errors := scanAllTokens(t, tc.input)

			if len(errors) != 0 {
				t.Errorf("Unexpected errors: %v", errors)
			}
			if len(tokens) != 2 { // Token + EOF
				t.Errorf("Expected 2 tokens (token + EOF), got %d", len(tokens))
			}
			if tokens[0].TypeOfToken != tc.expected {
				t.Errorf("Expected token type %v, got %v", tc.expected, tokens[0].TypeOfToken)
			}
		})
	}
}

func TestTokenRecognition2(t *testing.T) {
	testCases := []struct {
		input    string
		expected dfa.TokenType
	}{
		// Comments
		{"// comment", dfa.COMMENT},
		{"// co", dfa.COMMENT},
		{"//comment", dfa.COMMENT},
		{"//", dfa.COMMENT},

		// Identifiers
		{"identifier", dfa.IDENTIFIER},
		{"myVar", dfa.IDENTIFIER},
		{"_underscore", dfa.IDENTIFIER},
		{"var123", dfa.IDENTIFIER},

		// String literals
		{"\"\"", dfa.STRING},
		{"\"string\"", dfa.STRING},
		{"\".   ssw\"", dfa.STRING},

		// Number literals
		{"123", dfa.NUMBER},
		{"0", dfa.NUMBER},
		{"456.789", dfa.NUMBER},
		{"0.5", dfa.NUMBER},
		{"42", dfa.NUMBER},

		// Whitespace and newlines
		{" ", dfa.WHITESPACE},
		{" \t\t\t\t", dfa.WHITESPACE},
		{"\n", dfa.NEWLINE},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Identifier/Keyword: %s", tc.input), func(t *testing.T) {
			tokens, errors := scanAllTokens(t, tc.input)

			if len(errors) != 0 {
				t.Errorf("Unexpected errors: %v", errors)
			}
			if len(tokens) != 2 { // Token + EOF
				t.Errorf("Expected 2 tokens (token + EOF), got %d", len(tokens))
			}
			if tokens[0].TypeOfToken != tc.expected {
				t.Errorf("%+v\n", tokens)
				t.Errorf("%+v\n", tokens[0].TypeOfToken)
				t.Errorf("%+v\n", tc.expected)
				t.Errorf("Expected token type %v, got %v", tc.expected, tokens[0].TypeOfToken)
			}
		})
	}
}

// Test number literals
func TestNumberLiterals(t *testing.T) {
	testCases := []string{
		"123",
		"123.456",
		"0",
		"0.0",
		"999.999",
		"1.0",
	}

	for _, input := range testCases {
		t.Run(fmt.Sprintf("Number: %s", input), func(t *testing.T) {
			tokens, errors := scanAllTokens(t, input)

			if len(errors) != 0 {
				t.Errorf("Unexpected errors: %v", errors)
			}
			if len(tokens) != 2 { // Token + EOF
				t.Errorf("Expected 2 tokens (token + EOF), got %d", len(tokens))
			}
			if tokens[0].TypeOfToken != dfa.NUMBER {
				t.Errorf("Expected NUMBER token, got %v", tokens[0].TypeOfToken)
			}
		})
	}
}

// Test string literals
func TestStringLiterals(t *testing.T) {
	testCases := []string{
		`"hello"`,
		`"world"`,
		`""`,
		`"hello world"`,
		`"test\n"`,
	}

	for _, input := range testCases {
		t.Run(fmt.Sprintf("String: %s", input), func(t *testing.T) {
			tokens, errors := scanAllTokens(t, input)

			if len(errors) != 0 {
				t.Errorf("Unexpected errors: %v", errors)
			}
			if len(tokens) != 2 { // Token + EOF
				t.Errorf("Expected 2 tokens (token + EOF), got %d", len(tokens))
			}
			if tokens[0].TypeOfToken != dfa.STRING {
				t.Errorf("Expected STRING token, got %v", tokens[0].TypeOfToken)
			}
		})
	}
}

// Test whitespace handling
func TestWhitespaceHandling(t *testing.T) {
	testCases := []struct {
		input    string
		expected int // number of non-EOF tokens expected
	}{
		{"   ", 1},   // Only whitespace
		{"\t\t", 1},  // Only whitespace
		{"\n\n", 1},  // 1 newlines
		{"  a  ", 3}, // Whitespace around token
		{"a b", 3},   // Whitespace between tokens
		{"a\nb", 3},  // Newline between tokens
		{"a\tb", 3},  // Tab between tokens
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Whitespace: %q", tc.input), func(t *testing.T) {
			tokens, errors := scanAllTokens(t, tc.input)

			if len(errors) != 0 {
				t.Errorf("Unexpected errors: %v", errors)
			}

			// Count non-EOF tokens
			nonEOFCount := 0
			for _, token := range tokens {
				if token.TypeOfToken != dfa.EOF {
					nonEOFCount++
				}
			}

			if nonEOFCount != tc.expected {
				t.Errorf("Expected %d non-EOF tokens, got %d", tc.expected, nonEOFCount)
			}
		})
	}
}

// Test line and offset tracking
func TestLineAndOffsetTracking(t *testing.T) {
	input := "a\nb\nc"
	tokens, errors := scanAllTokens(t, input)

	if len(errors) != 0 {
		t.Errorf("Unexpected errors: %v", errors)
	}

	// Should have 3 identifier tokens + 2 newline tokens + 1 EOF token
	if len(tokens) != 6 {
		t.Errorf("Expected 4 tokens, got %d", len(tokens))
	}

	// Check line numbers
	expectedLines := []uint32{0, 0, 1, 1, 2, 2} // a, \n, b, \n, c, EOF
	for i, expectedLine := range expectedLines {
		if tokens[i].Line != expectedLine {
			t.Errorf("Token %d: expected line %d, got %d", i, expectedLine, tokens[i].Line)
		}
	}

	// Check offsets
	expectedOffsets := []uint32{0, 1, 0, 1, 0, 1} // a, b, c, EOF
	for i, expectedOffset := range expectedOffsets {
		if tokens[i].Offset != expectedOffset {
			t.Errorf("Token %d: expected offset %d, got %d", i, expectedOffset, tokens[i].Offset)
		}
	}
}

// Test maximal munching (longest match)
func TestMaximalMunching(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected dfa.TokenType
	}{
		// Multi-character operators
		{"EQUAL_EQUAL not two EQUAL", "==", dfa.EQUAL_EQUAL},
		{"BANG_EQUAL not BANG + EQUAL", "!=", dfa.BANG_EQUAL},
		{"LESS_EQUAL not LESS + EQUAL", "<=", dfa.LESS_EQUAL},
		{"GREATER_EQUAL not GREATER + EQUAL", ">=", dfa.GREATER_EQUAL},

		// Keywords vs identifiers - keywords should be recognized as keywords
		{"var keyword not identifier", "var", dfa.VAR},
		{"if keyword not identifier", "if", dfa.IF},
		{"else keyword not identifier", "else", dfa.ELSE},
		{"while keyword not identifier", "while", dfa.WHILE},
		{"for keyword not identifier", "for", dfa.FOR},
		{"fun keyword not identifier", "fun", dfa.FUN},
		{"return keyword not identifier", "return", dfa.RETURN},
		{"true keyword not identifier", "true", dfa.TRUE},
		{"false keyword not identifier", "false", dfa.FALSE},
		{"nil keyword not identifier", "nil", dfa.NIL},
		{"and keyword not identifier", "and", dfa.AND},
		{"or keyword not identifier", "or", dfa.OR},
		{"class keyword not identifier", "class", dfa.CLASS},
		{"print keyword not identifier", "print", dfa.PRINT},
		{"super keyword not identifier", "super", dfa.SUPER},
		{"this keyword not identifier", "this", dfa.THIS},

		// Similar looking identifiers that are NOT keywords
		{"variable identifier not keyword", "variable", dfa.IDENTIFIER},
		{"ifStatement identifier not keyword", "ifStatement", dfa.IDENTIFIER},
		{"elseBranch identifier not keyword", "elseBranch", dfa.IDENTIFIER},
		{"whileLoop identifier not keyword", "whileLoop", dfa.IDENTIFIER},
		{"forEach identifier not keyword", "forEach", dfa.IDENTIFIER},
		{"function identifier not keyword", "function", dfa.IDENTIFIER},
		{"returned identifier not keyword", "returned", dfa.IDENTIFIER},
		{"truly identifier not keyword", "truly", dfa.IDENTIFIER},
		{"falsely identifier not keyword", "falsely", dfa.IDENTIFIER},
		{"nill identifier not keyword", "nill", dfa.IDENTIFIER},
		{"andOr identifier not keyword", "andOr", dfa.IDENTIFIER},
		{"orElse identifier not keyword", "orElse", dfa.IDENTIFIER},
		{"className identifier not keyword", "className", dfa.IDENTIFIER},
		{"printer identifier not keyword", "printer", dfa.IDENTIFIER},
		{"superClass identifier not keyword", "superClass", dfa.IDENTIFIER},
		{"thisRef identifier not keyword", "thisRef", dfa.IDENTIFIER},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, errors := scanAllTokens(t, tc.input)

			if len(errors) != 0 {
				t.Errorf("Unexpected errors: %v", errors)
			}

			if len(tokens) != 2 { // Token + EOF
				t.Errorf("Expected 2 tokens, got %d", len(tokens))
			}

			if tokens[0].TypeOfToken != tc.expected {
				t.Errorf("Expected %v token, got %v", tc.expected, tokens[0].TypeOfToken)
			}
		})
	}
}

// Test invalid token handling
func TestInvalidTokenHandling(t *testing.T) {
	testCases := []string{
		"@", // Invalid character
		"#", // Invalid character
		"$", // Invalid character
		`"unclosed string`,
	}

	for _, input := range testCases {
		t.Run(fmt.Sprintf("Invalid: %s", input), func(t *testing.T) {
			tokens, errors := scanAllTokens(t, input)

			if len(errors) == 0 {
				t.Errorf("Expected error for invalid input '%s', but got none", input)
			}

			// Should still get EOF token
			if len(tokens) == 0 {
				t.Errorf("Expected at least EOF token, got none")
			}
		})
	}
}

// Test mixed token types
func TestMixedTokenTypes(t *testing.T) {
	input := `var x = 123.45;
	if (x > 0) {
		print "positive";
	}`

	tokens, errors := scanAllTokens(t, input)

	if len(errors) != 0 {
		t.Errorf("Unexpected errors: %v", errors)
	}

	// Expected token types in order
	expectedTypes := []dfa.TokenType{
		dfa.VAR,
		dfa.WHITESPACE,
		dfa.IDENTIFIER,
		dfa.WHITESPACE,
		dfa.EQUAL,
		dfa.WHITESPACE,
		dfa.NUMBER,
		dfa.SEMICOLON,
		dfa.NEWLINE,
		dfa.WHITESPACE,
		dfa.IF,
		dfa.WHITESPACE,
		dfa.LEFT_PAREN,
		dfa.IDENTIFIER,
		dfa.WHITESPACE,
		dfa.GREATER,
		dfa.WHITESPACE,
		dfa.NUMBER,
		dfa.RIGHT_PAREN,
		dfa.WHITESPACE,
		dfa.LEFT_BRACE,
		dfa.NEWLINE,
		dfa.WHITESPACE,
		dfa.PRINT,
		dfa.WHITESPACE,
		dfa.STRING,
		dfa.SEMICOLON,
		dfa.NEWLINE,
		dfa.WHITESPACE,
		dfa.RIGHT_BRACE,
		dfa.EOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Errorf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expectedType := range expectedTypes {
		if i < len(tokens) && tokens[i].TypeOfToken != expectedType {
			t.Errorf("Token %d: expected type %v, got %v", i, expectedType, tokens[i].TypeOfToken)
		}
	}
}

// Test error recovery
func TestErrorRecovery(t *testing.T) {
	// Mix valid and invalid tokens
	input := "var @ x = 123"

	tokens, errors := scanAllTokens(t, input)

	if len(errors) == 0 {
		t.Errorf("Expected error for invalid character '@'")
	}

	// Should recover and continue parsing after error
	validTokenFound := false
	for _, token := range tokens {
		if token.TypeOfToken == dfa.VAR || token.TypeOfToken == dfa.IDENTIFIER || token.TypeOfToken == dfa.NUMBER {
			validTokenFound = true
			break
		}
	}

	if !validTokenFound {
		t.Errorf("Expected to find valid tokens after error recovery")
	}
}

// Test edge cases with last input handling
func TestLastInputHandling(t *testing.T) {
	// Test case where a token is followed by an invalid character
	input := "var@"

	tokens, errors := scanAllTokens(t, input)

	// Should get VAR token and an error for '@'
	if len(errors) == 0 {
		t.Errorf("Expected error for invalid character")
	}

	// Should still get VAR token before the error
	foundVar := false
	for _, token := range tokens {
		if token.TypeOfToken == dfa.VAR {
			foundVar = true
			break
		}
	}

	if !foundVar {
		t.Errorf("Expected VAR token before error")
	}
}

// Test consecutive tokens without whitespace
func TestConsecutiveTokens(t *testing.T) {
	input := "(){};,.-+*/"

	tokens, errors := scanAllTokens(t, input)

	if len(errors) != 0 {
		t.Errorf("Unexpected errors: %v", errors)
	}

	expectedTypes := []dfa.TokenType{
		dfa.LEFT_PAREN,
		dfa.RIGHT_PAREN,
		dfa.LEFT_BRACE,
		dfa.RIGHT_BRACE,
		dfa.SEMICOLON,
		dfa.COMMA,
		dfa.DOT,
		dfa.MINUS,
		dfa.PLUS,
		dfa.STAR,
		dfa.SLASH,
		dfa.EOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Errorf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expectedType := range expectedTypes {
		if i < len(tokens) && tokens[i].TypeOfToken != expectedType {
			t.Errorf("Token %d: expected type %v, got %v", i, expectedType, tokens[i].TypeOfToken)
		}
	}
}

// Test complex expression
func TestComplexExpression(t *testing.T) {
	input := `(x + y) * 2.5 == result`

	tokens, errors := scanAllTokens(t, input)

	if len(errors) != 0 {
		t.Errorf("Unexpected errors: %v", errors)
	}

	expectedTypes := []dfa.TokenType{
		dfa.LEFT_PAREN,
		dfa.IDENTIFIER,
		dfa.WHITESPACE,
		dfa.PLUS,
		dfa.WHITESPACE,
		dfa.IDENTIFIER,
		dfa.RIGHT_PAREN,
		dfa.WHITESPACE,
		dfa.STAR,
		dfa.WHITESPACE,
		dfa.NUMBER,
		dfa.WHITESPACE,
		dfa.EQUAL_EQUAL,
		dfa.WHITESPACE,
		dfa.IDENTIFIER,
		dfa.EOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Errorf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expectedType := range expectedTypes {
		if i < len(tokens) && tokens[i].TypeOfToken != expectedType {
			t.Errorf("Token %d: expected type %v, got %v", i, expectedType, tokens[i].TypeOfToken)
		}
	}
}

// Benchmark test for performance
func BenchmarkScanner(b *testing.B) {
	input := strings.Repeat(`var x = 123.45;
	for (var i = 0; i < 10; i = i + 1) {
		if (x > i) {
			print "x is greater than " + i;
		}
	}
	`, 1000)

	b.ReportAllocs()

	scanner := &lexer.LexicalAnalyzer{}
	reader := bufio.NewReader(strings.NewReader(input))

	for b.Loop() {
		scanner.Initialize(reader)
		for {
			_, err := scanner.ReadToken()
			if err == io.EOF {
				break
			}
		}
		scanner.Reset()
	}
}

// TODO: Add a test case for checking if the scanner works even if the ReadToken function is called after the EOF is reached
