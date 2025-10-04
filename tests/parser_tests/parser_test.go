package parser_tests

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/parser"
)

// Recursively check that the AST matches the expected nesting order
func checkNodeNestingInOrder(nodes []parser.Node, expectedTypes []string) bool {
	if len(expectedTypes) == 0 {
		return true // All expected levels matched
	}
	if len(nodes) == 0 {
		return false // No nodes left to check
	}

	for _, node := range nodes {
		var nodeType string
		var children []parser.Node

		switch n := node.(type) {
		case *parser.Grammar_grammar:
			nodeType = "grammar"
			children = n.Arguments
		case *parser.Grammar_expression:
			nodeType = "expression"
			children = n.Arguments
		case *parser.Grammar_comma:
			nodeType = "comma"
			children = n.Arguments
		case *parser.Grammar_equality:
			nodeType = "equality"
			children = n.Arguments
		case *parser.Grammar_comparison:
			nodeType = "comparison"
			children = n.Arguments
		case *parser.Grammar_term:
			nodeType = "term"
			children = n.Arguments
		case *parser.Grammar_factor:
			nodeType = "factor"
			children = n.Arguments
		case *parser.Grammar_unary:
			nodeType = "unary"
			children = n.Arguments
		case *parser.Grammar_primary:
			nodeType = "primary"
			children = n.Arguments
		default:
			continue // skip other nodes like Literal
		}

		if nodeType == expectedTypes[0] {
			// Match found, move to next expected type
			if checkNodeNestingInOrder(children, expectedTypes[1:]) {
				return true
			}
		}
	}

	return false // No matching node found at this level
}

// Recursively print the AST with indentation
func logAST(nodes []parser.Node, indent string) {
	for _, node := range nodes {
		switch n := node.(type) {
		case *parser.Grammar_grammar:
			fmt.Println(indent + "Grammar_grammar")
			logAST(n.Arguments, indent+"\t")
		case *parser.Grammar_expression:
			fmt.Println(indent + "Grammar_expression")
			logAST(n.Arguments, indent+"\t")
		case *parser.Grammar_comma:
			fmt.Println(indent + "Grammar_comma")
			logAST(n.Arguments, indent+"\t")
		case *parser.Grammar_equality:
			fmt.Println(indent + "Grammar_equality")
			logAST(n.Arguments, indent+"\t")
		case *parser.Grammar_comparison:
			fmt.Println(indent + "Grammar_comparison")
			logAST(n.Arguments, indent+"\t")
		case *parser.Grammar_term:
			fmt.Println(indent + "Grammar_term")
			logAST(n.Arguments, indent+"\t")
		case *parser.Grammar_factor:
			fmt.Println(indent + "Grammar_factor")
			logAST(n.Arguments, indent+"\t")
		case *parser.Grammar_unary:
			fmt.Println(indent + "Grammar_unary")
			logAST(n.Arguments, indent+"\t")
		case *parser.Grammar_primary:
			fmt.Println(indent + "Grammar_primary")
			logAST(n.Arguments, indent+"\t")
		case *parser.Literal:
			fmt.Printf(indent+"Literal: `%s`\n", string(n.Value.Lexemme))
		default:
			fmt.Printf(indent+"Unknown node type: %T\n", n)
		}
	}
}

// func Test_ParseExpressionMinimal(t *testing.T) {
// 	const scanner_buf_cap uint32 = 10

// 	code := `1+2*3`
// 	scanner := lexer.BufferedLexer{}
// 	buf_reader := bufio.NewReader(strings.NewReader(code))
// 	scanner.Initialize(buf_reader, scanner_buf_cap)

// 	nodes, ok, err := parser.Parse_grammar(&scanner)
// 	if err != nil {
// 		t.Fatalf("Parsing failed: %v", err)
// 	}
// 	if !ok {
// 		t.Fatal("Parsing returned not ok")
// 	}
// 	if len(nodes) == 0 {
// 		t.Fatal("No nodes returned")
// 	}

// 	logAST(nodes, "")

// 	// Validate presence of expected non-terminals
// 	expectedOrder := []string{
// 		"grammar",
// 		"expression",
// 		"comma",
// 		"equality",
// 		"comparison",
// 		"term",
// 		"factor",
// 		"unary",
// 		"primary",
// 	}
// 	if !checkNodeNestingInOrder(nodes, expectedOrder) {
// 		t.Errorf("AST does not match expected nesting order: %v", expectedOrder)
// 	}

// }

// TODO: Fix the error of tokens being consumed even when a non-terminal ends up not matching completely

func Test_ParserCases(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectSuccess bool
		expectedOrder []string
	}{
		// {
		// 	name:          "Simple precedence test",
		// 	code:          "1+2*3",
		// 	expectSuccess: true,
		// 	expectedOrder: []string{"grammar", "expression", "comma", "equality", "comparison", "term", "factor", "unary", "primary"},
		// },
		// {
		// 	name:          "Parentheses override precedence",
		// 	code:          "(1+2)*3",
		// 	expectSuccess: true,
		// 	expectedOrder: []string{"grammar", "expression", "comma", "equality", "comparison", "term", "factor", "unary", "primary"},
		// },
		// {
		// 	name:          "Multiple unary",
		// 	code:          "--1",
		// 	expectSuccess: true,
		// 	expectedOrder: []string{"grammar", "expression", "comma", "equality", "comparison", "term", "factor", "unary", "unary", "unary", "primary"},
		// },
		// {
		// 	name:          "Comma-separated expressions",
		// 	code:          "a,b,c",
		// 	expectSuccess: true,
		// 	expectedOrder: []string{"grammar", "expression", "comma", "equality", "comparison", "term", "factor", "unary", "primary"},
		// },
		// {
		// 	name:          "Equality and comparison",
		// 	code:          "1==2<=3",
		// 	expectSuccess: true,
		// 	expectedOrder: []string{"grammar", "expression", "comma", "equality", "comparison", "term", "factor", "unary", "primary"},
		// },
		// {
		// 	name:          "Booleans and nil",
		// 	code:          "true!=false==nil",
		// 	expectSuccess: true,
		// 	expectedOrder: []string{"grammar", "expression", "comma", "equality", "comparison", "term", "factor", "unary", "primary"},
		// },
		// {
		// 	name:          "Strings",
		// 	code:          `"hello"+"world"`,
		// 	expectSuccess: true,
		// 	expectedOrder: []string{"grammar", "expression", "comma", "equality", "comparison", "term", "factor", "unary", "primary"},
		// },
		{
			name:          "Unmatched parentheses (syntax error)",
			code:          "(1+2",
			expectSuccess: false,
		},
		// {
		// 	name:          "Trailing garbage (partial parse)",
		// 	code:          "1+2*3extra",
		// 	expectSuccess: false,
		// },
		// {
		// 	name:          "Empty input",
		// 	code:          "",
		// 	expectSuccess: false,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			const scanner_buf_cap uint32 = 10

			scanner := lexer.BufferedLexer{}
			buf_reader := bufio.NewReader(strings.NewReader(test.code))
			scanner.Initialize(buf_reader, scanner_buf_cap)

			nodes, ok, err := parser.Parse_grammar(&scanner)
			if err != nil && test.expectSuccess {
				t.Fatalf("Unexpected parse error: %v", err)
			}
			if !ok && test.expectSuccess {
				t.Fatal("Parse returned not ok")
			}
			if ok && !test.expectSuccess {
				t.Fatal("Expected failure but got success")
			}
			if !test.expectSuccess {
				return // no need to check nesting
			}

			if len(nodes) == 0 {
				t.Fatal("No nodes returned")
			}

			// Print for visual debug (optional)
			logAST(nodes, "")

			if test.expectedOrder != nil {
				if !checkNodeNestingInOrder(nodes, test.expectedOrder) {
					t.Errorf("AST does not match expected nesting order: %v", test.expectedOrder)
				}
			}
		})
	}
}
