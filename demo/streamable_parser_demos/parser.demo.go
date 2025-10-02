package streamable_parser_demos

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/streamable_parser"
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

func prettyPrintEmit(e *streamable_parser.EmitElem) string {
	if e == nil {
		return "<nil>"
	}

	switch e.Type {
	case streamable_parser.EmitElemType_Start:
		return fmt.Sprintf("Start: %s", e.Content)

	case streamable_parser.EmitElemType_End:
		return fmt.Sprintf("End: %s", e.Content)

	case streamable_parser.EmitElemType_Leaf:
		if e.Leaf != nil {
			return fmt.Sprintf("Leaf: %s (%s)", e.Content, string(e.Leaf.Lexemme))
		}
		return fmt.Sprintf("Leaf: %s", e.Content)

	case streamable_parser.EmitElemType_Error:
		return fmt.Sprintf("Error: %s", e.Content)

	default:
		return fmt.Sprintf("Unknown: %s", e.Content)
	}
}

func Sample_streamable_parser_demo() {
	wd, _ := os.Getwd()
	fmt.Println("Working dir:", wd)

	// Sample input: conforms to your grammar
	file_reader := strings.NewReader(`foo+42,("bar"==baz)`)
	source := bufio.NewReader(file_reader)

	// Initialize lexer and parser
	buf_scanner := lexer.BufferedLexicalAnalyzer{}
	buf_scanner.Initialize(source)

	parser := streamable_parser.StreamableParser{}
	parser.Initialize(&buf_scanner)

	// Collect parse events
	var out []*streamable_parser.EmitElem
	for i := 0; i < 100; i++ {
		evt := parser.Parse()
		fmt.Println(prettyPrintEmit(evt))
		out = append(out, evt)
	}

	// Build tree and pretty print it
	root := buildTree(out)
	prettyPrint(root, "", true) // root has no prefix and is treated as "last"
}
