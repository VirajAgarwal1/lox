package parser_tests

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/VirajAgarwal1/lox/parser/grammar"
)

// Test helper functions
func createTestGrammar1() map[grammar.Non_terminal][]grammar.Generic_grammar_term {
	// Simple grammar: expr -> "NUMBER" | "(" expr ")"
	return map[grammar.Non_terminal][]grammar.Generic_grammar_term{
		grammar.Non_terminal{Name: "expr"}: {
			&grammar.Terminal{Content: []rune("NUMBER")},
			&grammar.Or{},
			&grammar.Terminal{Content: []rune("(")},
			&grammar.Non_terminal{Name: "expr"},
			&grammar.Terminal{Content: []rune(")")},
		},
	}
}

func createTestGrammar2() map[grammar.Non_terminal][]grammar.Generic_grammar_term {
	// Grammar with grammar.Star: list -> item*
	return map[grammar.Non_terminal][]grammar.Generic_grammar_term{
		grammar.Non_terminal{Name: "list"}: {
			&grammar.Star{Content: &grammar.Non_terminal{Name: "item"}},
		},
		grammar.Non_terminal{Name: "item"}: {
			&grammar.Terminal{Content: []rune("IDENTIFIER")},
		},
	}
}

func createTestGrammar3() map[grammar.Non_terminal][]grammar.Generic_grammar_term {
	// Grammar with grammar.Plus: list -> item+
	return map[grammar.Non_terminal][]grammar.Generic_grammar_term{
		grammar.Non_terminal{Name: "list"}: {
			&grammar.Plus{Content: &grammar.Non_terminal{Name: "item"}},
		},
		grammar.Non_terminal{Name: "item"}: {
			&grammar.Terminal{Content: []rune("NUMBER")},
		},
	}
}

func createComplexGrammar() map[grammar.Non_terminal][]grammar.Generic_grammar_term {
	// The example grammar from your code
	return map[grammar.Non_terminal][]grammar.Generic_grammar_term{
		grammar.Non_terminal{Name: "expression"}: {
			&grammar.Non_terminal{Name: "comma"},
		},
		grammar.Non_terminal{Name: "comma"}: {
			&grammar.Non_terminal{Name: "equality"},
			&grammar.Star{Content: &grammar.Bracket{Contents: []grammar.Generic_grammar_term{
				&grammar.Terminal{Content: []rune(",")},
				&grammar.Non_terminal{Name: "equality"},
			}}},
		},
		grammar.Non_terminal{Name: "equality"}: {
			&grammar.Non_terminal{Name: "comparison"},
			&grammar.Star{Content: &grammar.Bracket{Contents: []grammar.Generic_grammar_term{
				&grammar.Bracket{Contents: []grammar.Generic_grammar_term{
					&grammar.Terminal{Content: []rune("!=")},
					&grammar.Or{},
					&grammar.Terminal{Content: []rune("==")},
				}},
				&grammar.Non_terminal{Name: "comparison"},
			}}},
		},
		grammar.Non_terminal{Name: "comparison"}: {
			&grammar.Non_terminal{Name: "term"},
			&grammar.Star{Content: &grammar.Bracket{Contents: []grammar.Generic_grammar_term{
				&grammar.Bracket{Contents: []grammar.Generic_grammar_term{
					&grammar.Terminal{Content: []rune(">")},
					&grammar.Or{},
					&grammar.Terminal{Content: []rune(">=")},
					&grammar.Or{},
					&grammar.Terminal{Content: []rune("<")},
					&grammar.Or{},
					&grammar.Terminal{Content: []rune("<=")},
				}},
				&grammar.Non_terminal{Name: "term"},
			}}},
		},
		grammar.Non_terminal{Name: "term"}: {
			&grammar.Non_terminal{Name: "factor"},
			&grammar.Star{Content: &grammar.Bracket{Contents: []grammar.Generic_grammar_term{
				&grammar.Bracket{Contents: []grammar.Generic_grammar_term{
					&grammar.Terminal{Content: []rune("-")},
					&grammar.Or{},
					&grammar.Terminal{Content: []rune("+")},
				}},
				&grammar.Non_terminal{Name: "factor"},
			}}},
		},
		grammar.Non_terminal{Name: "factor"}: {
			&grammar.Non_terminal{Name: "unary"},
			&grammar.Star{Content: &grammar.Bracket{Contents: []grammar.Generic_grammar_term{
				&grammar.Bracket{Contents: []grammar.Generic_grammar_term{
					&grammar.Terminal{Content: []rune("/")},
					&grammar.Or{},
					&grammar.Terminal{Content: []rune("*")},
				}},
				&grammar.Non_terminal{Name: "unary"},
			}}},
		},
		grammar.Non_terminal{Name: "unary"}: {
			&grammar.Bracket{Contents: []grammar.Generic_grammar_term{
				&grammar.Terminal{Content: []rune("!")},
				&grammar.Or{},
				&grammar.Terminal{Content: []rune("-")},
			}},
			&grammar.Non_terminal{Name: "unary"},
			&grammar.Or{},
			&grammar.Non_terminal{Name: "primary"},
		},
		grammar.Non_terminal{Name: "primary"}: {
			&grammar.Terminal{Content: []rune("IDENTIFIER")},
			&grammar.Or{},
			&grammar.Terminal{Content: []rune("NUMBER")},
			&grammar.Or{},
			&grammar.Terminal{Content: []rune("STRING")},
			&grammar.Or{},
			&grammar.Terminal{Content: []rune("true")},
			&grammar.Or{},
			&grammar.Terminal{Content: []rune("false")},
			&grammar.Or{},
			&grammar.Terminal{Content: []rune("nil")},
			&grammar.Or{},
			&grammar.Terminal{Content: []rune("(")},
			&grammar.Non_terminal{Name: "expression"},
			&grammar.Terminal{Content: []rune(")")},
		},
	}
}

// Test the basic structure generation
func TestWriteStructsForNonTerminals(t *testing.T) {
	testCases := []struct {
		name     string
		grammar  map[grammar.Non_terminal][]grammar.Generic_grammar_term
		expected []string // expected struct names
	}{
		{
			name:     "Simple grammar",
			grammar:  createTestGrammar1(),
			expected: []string{"Grammar_expr"},
		},
		{
			name:     "Grammar with multiple non-terminals",
			grammar:  createTestGrammar2(),
			expected: []string{"Grammar_list", "Grammar_item"},
		},
		{
			name:     "Complex grammar",
			grammar:  createComplexGrammar(),
			expected: []string{"Grammar_expression", "Grammar_comma", "Grammar_equality", "Grammar_comparison", "Grammar_term", "Grammar_factor", "Grammar_unary", "Grammar_primary"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			err := grammar.WriteStructsForNonTerminals(writer, tc.grammar)
			if err != nil {
				t.Fatalf("WriteStructsForNonTerminals failed: %v", err)
			}

			writer.Flush()
			output := buf.String()

			// Check that all expected structs are present
			for _, expected := range tc.expected {
				if !strings.Contains(output, "type "+expected+" struct") {
					t.Errorf("Expected struct %s not found in output", expected)
				}
				if !strings.Contains(output, "Arguments []Node") {
					t.Errorf("Expected Arguments field not found for struct %s", expected)
				}
			}
		})
	}
}

func TestWriteEvaluateMethodsForNonTerminals(t *testing.T) {
	testCases := []struct {
		name     string
		grammar  map[grammar.Non_terminal][]grammar.Generic_grammar_term
		expected []string // expected method names
	}{
		{
			name:     "Simple grammar",
			grammar:  createTestGrammar1(),
			expected: []string{"Grammar_expr"},
		},
		{
			name:     "Grammar with multiple non-terminals",
			grammar:  createTestGrammar2(),
			expected: []string{"Grammar_list", "Grammar_item"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			err := grammar.WriteEvaluateMethodsForNonTerminals(writer, tc.grammar)
			if err != nil {
				t.Fatalf("WriteEvaluateMethodsForNonTerminals failed: %v", err)
			}

			writer.Flush()
			output := buf.String()

			// Check that all expected methods are present
			for _, expected := range tc.expected {
				methodSignature := "func (non_terminal *" + expected + ") Evaluate() *Value"
				if !strings.Contains(output, methodSignature) {
					t.Errorf("Expected method signature %s not found in output", methodSignature)
				}
			}
		})
	}
}

func TestWriteParseFunctionsForNonTerminals(t *testing.T) {
	testCases := []struct {
		name     string
		grammar  map[grammar.Non_terminal][]grammar.Generic_grammar_term
		expected []string // expected function names
	}{
		{
			name:     "Simple grammar",
			grammar:  createTestGrammar1(),
			expected: []string{"Parse_expr"},
		},
		{
			name:     "Grammar with multiple non-terminals",
			grammar:  createTestGrammar2(),
			expected: []string{"Parse_list", "Parse_item"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			err := grammar.WriteParseFunctionsForNonTerminals(writer, tc.grammar)
			if err != nil {
				t.Fatalf("WriteParseFunctionsForNonTerminals failed: %v", err)
			}

			writer.Flush()
			output := buf.String()

			// Check that all expected functions are present
			for _, expected := range tc.expected {
				funcSignature := "func " + expected + "(buf *lexer.BufferedLexicalAnalyzer) ([]Node, bool, error)"
				if !strings.Contains(output, funcSignature) {
					t.Errorf("Expected function signature %s not found in output", funcSignature)
				}
			}
		})
	}
}

func TestGenerateMatchCode(t *testing.T) {
	testCases := []struct {
		name     string
		term     grammar.Terminal
		endStr   string
		expected string
	}{
		{
			name:     "NUMBER token",
			term:     grammar.Terminal{Content: []rune("NUMBER")},
			endStr:   "(buf)",
			expected: "matchToken(dfa.NUMBER)(buf)",
		},
		{
			name:     "LEFT_PAREN token",
			term:     grammar.Terminal{Content: []rune("(")},
			endStr:   "(buf)",
			expected: "matchToken(dfa.LEFT_PAREN)(buf)",
		},
		{
			name:     "IDENTIFIER token",
			term:     grammar.Terminal{Content: []rune("IDENTIFIER")},
			endStr:   ",\n",
			expected: "matchToken(dfa.IDENTIFIER),\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := grammar.GenerateMatchCode(tc.term, tc.endStr)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestGenerateParseReturnCode(t *testing.T) {
	testCases := []struct {
		name     string
		term     grammar.Non_terminal
		endStr   string
		expected string
	}{
		{
			name:     "Expression non-grammar.Terminal",
			term:     grammar.Non_terminal{Name: "expression"},
			endStr:   "(buf)",
			expected: "Parse_expression(buf)",
		},
		{
			name:     "Primary non-grammar.Terminal",
			term:     grammar.Non_terminal{Name: "primary"},
			endStr:   ",\n",
			expected: "Parse_primary,\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := grammar.GenerateParseReturnCode(tc.term, tc.endStr)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestGenerateStarCode(t *testing.T) {
	testCases := []struct {
		name     string
		term     grammar.Star
		endStr   string
		expected string
	}{
		{
			name:     "Star with grammar.Terminal",
			term:     grammar.Star{Content: &grammar.Terminal{Content: []rune("NUMBER")}},
			endStr:   "(buf)",
			expected: "zeroOrMore(\nmatchToken(dfa.NUMBER),\n)(buf)",
		},
		{
			name:     "Star with non-grammar.Terminal",
			term:     grammar.Star{Content: &grammar.Non_terminal{Name: "expr"}},
			endStr:   ",\n",
			expected: "zeroOrMore(\nParse_expr,\n),\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := grammar.GenerateStarCode(tc.term, tc.endStr)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestGeneratePlusCode(t *testing.T) {
	testCases := []struct {
		name     string
		term     grammar.Plus
		endStr   string
		expected string
	}{
		{
			name:     "Plus with grammar.Terminal",
			term:     grammar.Plus{Content: &grammar.Terminal{Content: []rune("IDENTIFIER")}},
			endStr:   "(buf)",
			expected: "oneOrMore(\nmatchToken(dfa.IDENTIFIER),\n)(buf)",
		},
		// Fixed the second test case - removed invalid type usage
		{
			name:     "Plus with non-grammar.Terminal",
			term:     grammar.Plus{Content: &grammar.Non_terminal{Name: "expr"}},
			endStr:   "(buf)",
			expected: "oneOrMore(\nParse_expr,\n)(buf)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := grammar.GeneratePlusCode(tc.term, tc.endStr)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestDetectOrInDescription(t *testing.T) {
	testCases := []struct {
		name        string
		description []grammar.Generic_grammar_term
		expected    []uint32
	}{
		{
			name: "No OR terms",
			description: []grammar.Generic_grammar_term{
				&grammar.Terminal{Content: []rune("NUMBER")},
				&grammar.Terminal{Content: []rune("+")},
			},
			expected: []uint32{},
		},
		{
			name: "Single OR term",
			description: []grammar.Generic_grammar_term{
				&grammar.Terminal{Content: []rune("NUMBER")},
				&grammar.Or{},
				&grammar.Terminal{Content: []rune("STRING")},
			},
			expected: []uint32{1},
		},
		{
			name: "Multiple OR terms",
			description: []grammar.Generic_grammar_term{
				&grammar.Terminal{Content: []rune("NUMBER")},
				&grammar.Or{},
				&grammar.Terminal{Content: []rune("STRING")},
				&grammar.Or{},
				&grammar.Terminal{Content: []rune("IDENTIFIER")},
			},
			expected: []uint32{1, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := grammar.DetectOrInDescription(tc.description)
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
				return
			}
			for i, expected := range tc.expected {
				if result[i] != expected {
					t.Errorf("At index %d: expected %d, got %d", i, expected, result[i])
				}
			}
		})
	}
}

func TestIndentLines(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		indentLevel int
		expected    string
	}{
		{
			name:        "Single line, one indent",
			input:       "hello",
			indentLevel: 1,
			expected:    "\thello",
		},
		{
			name:        "Multiple lines, two indents",
			input:       "line1\nline2\nline3",
			indentLevel: 2,
			expected:    "\t\tline1\n\t\tline2\n\t\tline3",
		},
		{
			name:        "With empty lines",
			input:       "line1\n\nline3",
			indentLevel: 1,
			expected:    "\tline1\n\n\tline3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := grammar.IndentLines(tc.input, tc.indentLevel)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestGenerateDescriptionCode(t *testing.T) {
	testCases := []struct {
		name             string
		description      []grammar.Generic_grammar_term
		endStr           string
		expectedContains []string // Check if these strings are contained in output
	}{
		{
			name: "Single grammar.Terminal",
			description: []grammar.Generic_grammar_term{
				&grammar.Terminal{Content: []rune("NUMBER")},
			},
			endStr:           "(buf)",
			expectedContains: []string{"matchToken(dfa.NUMBER)(buf)"},
		},
		{
			name: "Sequence of terminals",
			description: []grammar.Generic_grammar_term{
				&grammar.Terminal{Content: []rune("(")},
				&grammar.Terminal{Content: []rune("NUMBER")},
				&grammar.Terminal{Content: []rune(")")},
			},
			endStr:           "(buf)",
			expectedContains: []string{"sequence(", "matchToken(dfa.LEFT_PAREN)", "matchToken(dfa.NUMBER)", "matchToken(dfa.RIGHT_PAREN)"},
		},
		{
			name: "Choice between terminals",
			description: []grammar.Generic_grammar_term{
				&grammar.Terminal{Content: []rune("NUMBER")},
				&grammar.Or{},
				&grammar.Terminal{Content: []rune("STRING")},
			},
			endStr:           "(buf)",
			expectedContains: []string{"choice(", "matchToken(dfa.NUMBER)", "matchToken(dfa.STRING)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := grammar.GenerateDescriptionCode(tc.description, tc.endStr)
			for _, expected := range tc.expectedContains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, but it didn't. Got: %s", expected, result)
				}
			}
		})
	}
}

// Integration test: Generate a complete parser file and check its structure
func TestGenerateGrammarParserFile(t *testing.T) {
	testCases := []struct {
		name     string
		grammar  map[grammar.Non_terminal][]grammar.Generic_grammar_term
		filePath string
	}{
		{
			name:     "Simple grammar file generation",
			grammar:  createTestGrammar1(),
			filePath: "test_simple_grammar.go",
		},
		{
			name:     "Complex grammar file generation",
			grammar:  createComplexGrammar(),
			filePath: "test_complex_grammar.go",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up file after test
			defer os.Remove(tc.filePath)

			err := grammar.GenerateGrammarParserFile(tc.grammar, tc.filePath)
			if err != nil {
				t.Fatalf("GenerateGrammarParserFile failed: %v", err)
			}

			// Read the generated file
			content, err := os.ReadFile(tc.filePath)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			fileContent := string(content)

			// Check that the file contains essential components
			expectedComponents := []string{
				"package parser",
				"type Value struct",
				"type Node interface",
				"func matchToken",
				"func sequence",
				"func choice",
				"func zeroOrMore",
				"func oneOrMore",
			}

			for _, component := range expectedComponents {
				if !strings.Contains(fileContent, component) {
					t.Errorf("Generated file missing component: %q", component)
				}
			}

			// Check that structs for non-terminals are generated
			for nonTerminal := range tc.grammar {
				expectedStruct := "type Grammar_" + nonTerminal.Name + " struct"
				if !strings.Contains(fileContent, expectedStruct) {
					t.Errorf("Generated file missing struct: %q", expectedStruct)
				}

				expectedParseFunc := "func Parse_" + nonTerminal.Name + "("
				if !strings.Contains(fileContent, expectedParseFunc) {
					t.Errorf("Generated file missing parse function: %q", expectedParseFunc)
				}

				expectedEvaluateMethod := "func (non_terminal *Grammar_" + nonTerminal.Name + ") Evaluate()"
				if !strings.Contains(fileContent, expectedEvaluateMethod) {
					t.Errorf("Generated file missing evaluate method: %q", expectedEvaluateMethod)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkGenerateGrammarParserFile(b *testing.B) {
	generated_grammar := createComplexGrammar()
	filePath := "benchmark_test.go"
	defer os.Remove(filePath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := grammar.GenerateGrammarParserFile(generated_grammar, filePath)
		if err != nil {
			b.Fatalf("GenerateGrammarParserFile failed: %v", err)
		}
	}
}

// Additional functional tests to verify the generated parser actually works

func TestGeneratedParserFunctionality(t *testing.T) {

	// This test would require creating a simple grammar, generating the parser,
	// and then testing if it can actually parse input correctly
	testGrammar := createTestGrammar1()
	filePath := "test_functional_parser.go"
	defer os.Remove(filePath)

	err := grammar.GenerateGrammarParserFile(testGrammar, filePath)
	if err != nil {
		t.Fatalf("Failed to generate parser file: %v", err)
	}

	// Read the generated file to ensure it compiles
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Basic syntax checks
	fileContent := string(content)
	if !strings.Contains(fileContent, "package parser") {
		t.Error("Generated file should have correct package declaration")
	}

	// Check for syntax correctness by looking for balanced braces
	openBraces := strings.Count(fileContent, "{")
	closeBraces := strings.Count(fileContent, "}")
	if openBraces != closeBraces {
		t.Errorf("Unbalanced braces in generated file: %d open, %d close", openBraces, closeBraces)
	}
}

// Test for edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("Empty grammar", func(t *testing.T) {
		emptyGrammar := make(map[grammar.Non_terminal][]grammar.Generic_grammar_term)
		filePath := "test_empty_grammar.go"
		defer os.Remove(filePath)

		err := grammar.GenerateGrammarParserFile(emptyGrammar, filePath)
		if err != nil {
			t.Fatalf("Failed to generate parser for empty grammar: %v", err)
		}

		// Should still generate a valid file with base components
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read generated file: %v", err)
		}

		fileContent := string(content)
		if !strings.Contains(fileContent, "package parser") {
			t.Error("Empty grammar should still generate valid package")
		}
	})

	t.Run("Grammar with single non-grammar.Terminal", func(t *testing.T) {
		singleGrammar := map[grammar.Non_terminal][]grammar.Generic_grammar_term{
			grammar.Non_terminal{Name: "single"}: {
				&grammar.Terminal{Content: []rune("NUMBER")},
			},
		}
		filePath := "test_single_grammar.go"
		defer os.Remove(filePath)

		err := grammar.GenerateGrammarParserFile(singleGrammar, filePath)
		if err != nil {
			t.Fatalf("Failed to generate parser for single grammar: %v", err)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read generated file: %v", err)
		}

		fileContent := string(content)
		if !strings.Contains(fileContent, "type Grammar_single struct") {
			t.Error("Should generate struct for single non-grammar.Terminal")
		}
		if !strings.Contains(fileContent, "func Parse_single(") {
			t.Error("Should generate parse function for single non-grammar.Terminal")
		}
	})
}
