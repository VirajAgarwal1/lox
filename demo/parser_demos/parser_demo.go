package parser_demos

// import (
// 	"bufio"
// 	"fmt"
// 	"strings"

// 	"github.com/VirajAgarwal1/lox/lexer"
// 	"github.com/VirajAgarwal1/lox/parser"
// )

// // Recursively print AST
// func printAST(node parser.Node, depth int) {
// 	indent := func(n int) string {
// 		return strings.Repeat("  ", n) // 2 spaces per level
// 	}

// 	switch n := node.(type) {
// 	case *parser.Grammar_expression:
// 		fmt.Println(indent(depth), "Expr:")
// 		for _, arg := range n.Arguments {
// 			printAST(arg, depth+1)
// 		}
// 	case *parser.Grammar_comma:
// 		fmt.Println(indent(depth), "Comma:")
// 		for _, arg := range n.Arguments {
// 			printAST(arg, depth+1)
// 		}
// 	case *parser.Grammar_primary:
// 		fmt.Println(indent(depth), "Primary:")
// 		for _, arg := range n.Arguments {
// 			printAST(arg, depth+1)
// 		}
// 	case *parser.Literal:
// 		fmt.Printf(indent(depth)+"Literal: %v\n", n.Value)
// 	default:
// 		fmt.Println(indent(depth), "Unknown Node")
// 	}
// }

// func ParserDemo() {
// 	// Sample input:  42, "hello", true
// 	sample_input := bufio.NewReader(strings.NewReader("42,\"hello\",true,identifier,false,2.89"))
// 	buf_scanner := lexer.BufferedLexicalAnalyzer{}
// 	buf_scanner.Initialize(sample_input)

// 	// Run the parser
// 	ast, err := parser.Parse_expression(&buf_scanner)
// 	if err != nil {
// 		fmt.Println("Parser error:", err)
// 		fmt.Println()
// 		fmt.Println()
// 	}

// 	// Print the AST
// 	printAST(ast, 0)
// }
