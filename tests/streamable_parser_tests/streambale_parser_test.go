package streamable_parser_tests

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/streamable_parser"
)

// helper: run parser fully and collect events
func collectEvents(input string) []*streamable_parser.EmitElem {
	source := strings.NewReader(input)
	bufSource := bufio.NewReader(source)
	scanner := lexer.BufferedLexicalAnalyzer{}
	scanner.Initialize(bufSource)

	var sp streamable_parser.StreamableParser
	sp.Initialize(&scanner)

	var events []*streamable_parser.EmitElem
	for {
		ev := sp.Parse()
		if ev.Type == streamable_parser.EmitElemType_Error && ev.Content == io.EOF.Error() {
			break
		}
		events = append(events, ev)
	}
	return events
}

func TestParserIntegration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []streamable_parser.EmitElemType
	}{
		{
			name:  "simple number literal",
			input: "42",
			expected: []streamable_parser.EmitElemType{
				// descent into grammar
				streamable_parser.EmitElemType_Start, // expression
				streamable_parser.EmitElemType_Start, // comma
				streamable_parser.EmitElemType_Start, // equality
				streamable_parser.EmitElemType_Start, // comparison
				streamable_parser.EmitElemType_Start, // term
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "42"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				streamable_parser.EmitElemType_End,   // factor
				streamable_parser.EmitElemType_End,   // term
				streamable_parser.EmitElemType_End,   // comparison
				streamable_parser.EmitElemType_End,   // equality
				streamable_parser.EmitElemType_End,   // comma
				streamable_parser.EmitElemType_End,   // expression
			},
		},
		{
			name:  "parenthesized expr",
			input: "(123)",
			expected: []streamable_parser.EmitElemType{
				streamable_parser.EmitElemType_Start, // expression
				streamable_parser.EmitElemType_Start, // comma
				streamable_parser.EmitElemType_Start, // equality
				streamable_parser.EmitElemType_Start, // comparison
				streamable_parser.EmitElemType_Start, // term
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "("
				// nested expression
				streamable_parser.EmitElemType_Start, // expression
				// ... inner nesting identical down to NUMBER
				streamable_parser.EmitElemType_Start, // comma
				streamable_parser.EmitElemType_Start, // equality
				streamable_parser.EmitElemType_Start, // comparison
				streamable_parser.EmitElemType_Start, // term
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "123"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				streamable_parser.EmitElemType_End,   // factor
				streamable_parser.EmitElemType_End,   // term
				streamable_parser.EmitElemType_End,   // comparison
				streamable_parser.EmitElemType_End,   // equality
				streamable_parser.EmitElemType_End,   // comma
				streamable_parser.EmitElemType_End,   // expression
				// back to outer primary
				streamable_parser.EmitElemType_Leaf, // ")"
				streamable_parser.EmitElemType_End,  // primary
				streamable_parser.EmitElemType_End,  // unary
				streamable_parser.EmitElemType_End,  // factor
				streamable_parser.EmitElemType_End,  // term
				streamable_parser.EmitElemType_End,  // comparison
				streamable_parser.EmitElemType_End,  // equality
				streamable_parser.EmitElemType_End,  // comma
				streamable_parser.EmitElemType_End,  // expression
			},
		},
		{
			name:  "binary expr",
			input: "1+2",
			expected: []streamable_parser.EmitElemType{
				streamable_parser.EmitElemType_Start, // expression
				streamable_parser.EmitElemType_Start, // comma
				streamable_parser.EmitElemType_Start, // equality
				streamable_parser.EmitElemType_Start, // comparison
				streamable_parser.EmitElemType_Start, // term
				// left side
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "1"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				streamable_parser.EmitElemType_End,   // factor
				// "+"
				streamable_parser.EmitElemType_Leaf,
				// right side
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "2"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				streamable_parser.EmitElemType_End,   // factor
				streamable_parser.EmitElemType_End,   // term
				streamable_parser.EmitElemType_End,   // comparison
				streamable_parser.EmitElemType_End,   // equality
				streamable_parser.EmitElemType_End,   // comma
				streamable_parser.EmitElemType_End,   // expression
			},
		},
		{
			name:  "parse error recovery",
			input: "1+*2",
			expected: []streamable_parser.EmitElemType{
				streamable_parser.EmitElemType_Start, // expression
				streamable_parser.EmitElemType_Start, // comma
				streamable_parser.EmitElemType_Start, // equality
				streamable_parser.EmitElemType_Start, // comparison
				streamable_parser.EmitElemType_Start, // term
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "1"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				streamable_parser.EmitElemType_End,   // factor
				streamable_parser.EmitElemType_Leaf,  // "+"
				streamable_parser.EmitElemType_Error, // error at "*"
				streamable_parser.EmitElemType_End,   // term
				streamable_parser.EmitElemType_End,   // comparison
				streamable_parser.EmitElemType_End,   // equality
				streamable_parser.EmitElemType_End,   // comma
				streamable_parser.EmitElemType_End,   // expression
			},
		},
		{
			name:  "complex mixed binary expr",
			input: `1+1+"hello"+2*3*3`,
			expected: []streamable_parser.EmitElemType{
				streamable_parser.EmitElemType_Start, // expression
				streamable_parser.EmitElemType_Start, // comma
				streamable_parser.EmitElemType_Start, // equality
				streamable_parser.EmitElemType_Start, // comparison
				streamable_parser.EmitElemType_Start, // term

				// leftmost "1"
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "1"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				streamable_parser.EmitElemType_End,   // factor

				// "+"
				streamable_parser.EmitElemType_Leaf,

				// next "1"
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "1"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				streamable_parser.EmitElemType_End,   // factor

				// "+"
				streamable_parser.EmitElemType_Leaf,

				// next "hello"
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "hello"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				streamable_parser.EmitElemType_End,   // factor

				// "+"
				streamable_parser.EmitElemType_Leaf,

				// right side: 2*3*3
				streamable_parser.EmitElemType_Start, // factor
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "2"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				// "*"
				streamable_parser.EmitElemType_Leaf,
				// next 3
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "3"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary
				// "*"
				streamable_parser.EmitElemType_Leaf,
				// next 3
				streamable_parser.EmitElemType_Start, // unary
				streamable_parser.EmitElemType_Start, // primary
				streamable_parser.EmitElemType_Leaf,  // "3"
				streamable_parser.EmitElemType_End,   // primary
				streamable_parser.EmitElemType_End,   // unary

				streamable_parser.EmitElemType_End, // factor (2*3*3)
				streamable_parser.EmitElemType_End, // term
				streamable_parser.EmitElemType_End, // comparison
				streamable_parser.EmitElemType_End, // equality
				streamable_parser.EmitElemType_End, // comma
				streamable_parser.EmitElemType_End, // expression
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := collectEvents(tt.input)

			if len(events) != len(tt.expected) {
				t.Fatalf("len mismatch: got %v, expected %v\nEvents: %#v", len(events), len(tt.expected), events)
			}
			for i := range events {
				if events[i].Type != tt.expected[i] {
					t.Errorf("mismatch at %d: got %v, expected %v", i, events[i].Type, tt.expected[i])
				}
			}
		})
	}
}
