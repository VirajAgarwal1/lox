package streamable_parser_demos

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/streamable_parser"
	ebnf_to_bnf "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/ebnf_to_bnf"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/first_follow"
	grammar_file_parser "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/grammar_file_parser"
	"github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/parser_writer"
)

type temp_node struct {
	content  string
	children []*temp_node
}

// BuildTree consumes parser events and returns a root node.
func buildTree(events []*streamable_parser.EmitElem) *temp_node {
	var stack []*temp_node
	var root *temp_node

	for _, e := range events {
		if e == nil {
			continue
		}

		switch e.Type {
		case streamable_parser.EmitElemType_Start:
			node := &temp_node{content: "⟨" + e.Content + "⟩"}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, node)
			} else {
				root = node
			}
			stack = append(stack, node)

		case streamable_parser.EmitElemType_End:
			if len(stack) > 0 {
				node := &temp_node{content: "⟨/" + e.Content + "⟩"}
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, node)
				stack = stack[:len(stack)-1]
			}

		case streamable_parser.EmitElemType_Leaf:
			content := e.Content
			if e.Leaf != nil {
				content = fmt.Sprintf("Leaf(%s)", string(e.Leaf.Lexemme))
			}
			node := &temp_node{content: content}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, node)
			}

		case streamable_parser.EmitElemType_Error:
			node := &temp_node{content: "ERR: " + e.Content}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.children = append(parent.children, node)
			} else if root == nil {
				root = node
			} else {
				root.children = append(root.children, node)
			}
		}
	}
	return root
}

// PrettyPrint prints the parse tree with nice ASCII symbols.
func prettyPrint(node *temp_node, prefix string, isLast bool) {
	if node == nil {
		return
	}

	// Choose connector └── for last child, ├── for others
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	// Print this node
	fmt.Println(prefix + connector + node.content)

	// Compute new prefix for children
	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}

	// Print children
	for i, child := range node.children {
		prettyPrint(child, newPrefix, i == len(node.children)-1)
	}
}

// Sample_generate_streamable_parser generates a StreamableParser from a grammar file.
// IMPORTANT: After running this, you must run your program AGAIN to use the generated parser.
// Go compiles code at build-time, so the newly generated code won't be available until the next run.
func Sample_generate_streamable_parser() {
	fmt.Println("=== Step 1: Generate Parser Code ===")
	fmt.Println()

	// Read and parse the grammar file
	grammarFile, err := os.Open("parser/lox.grammar")
	if err != nil {
		panic(err)
	}
	defer grammarFile.Close()

	scanner := lexer.LexicalAnalyzer{}
	scanner.Initialize(bufio.NewReader(grammarFile))
	ebnfGrammar, err := grammar_file_parser.ProcessGrammarDefinition(&scanner)
	if err != nil && err != io.EOF {
		panic(err)
	}

	// Convert EBNF to BNF and compute parsing tables
	bnfGrammar := ebnf_to_bnf.EbnfToBnfConverter(ebnfGrammar)
	firsts := first_follow.ComputeFirstSets(bnfGrammar)
	follows := first_follow.ComputeFollowSets(bnfGrammar)

	// Generate the parser code
	currentDir, _ := os.Getwd()
	outputPath := currentDir + "/streamable_parser/generated_parser.go"
	err = parser_writer.WriteParser(outputPath, bnfGrammar, "expression", firsts, follows)
	if err != nil {
		panic(err)
	}

	fmt.Printf("✓ Parser generated successfully!\n")
	fmt.Printf("  Output: %s\n", outputPath)
	fmt.Println()
	fmt.Println("NEXT: Change main.go to call Sample_streamable_parser_demo() and run again.")
}

// Sample_streamable_parser_demo uses the generated parser to parse an expression.
// IMPORTANT: You must run Sample_generate_streamable_parser() FIRST (in a separate execution)
// to generate the parser code before this demo will work.
func Sample_streamable_parser_demo() {
	fmt.Println("=== Step 2: Use the Generated Parser ===")
	fmt.Println()

	// Parse a sample expression
	testInput := `1+2*3-4/"hello"+true`
	fmt.Printf("Parsing: %s\n\n", testInput)

	// Initialize lexer
	bufferedScanner := lexer.BufferedLexicalAnalyzer{}
	bufferedScanner.Initialize(bufio.NewReader(strings.NewReader(testInput)))

	// Initialize parser
	parser := streamable_parser.StreamableParser{}
	parser.Initialize(&bufferedScanner)

	// Collect parse events
	var events []*streamable_parser.EmitElem
	for range 100 {
		evt := parser.Parse()
		if evt.Type == streamable_parser.EmitElemType_Error {
			if evt.Content == io.EOF.Error() {
				break
			}
			fmt.Printf("Parse error: %s\n", evt.Content)
			break
		}
		events = append(events, evt)
	}

	// Display parse tree
	fmt.Println("Parse tree:")
	root := buildTree(events)
	prettyPrint(root, "", true)

	fmt.Println("\n✓ Parsing complete!")
}
