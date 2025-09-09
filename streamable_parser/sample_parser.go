package streamableparser

import (
	"fmt"

	"github.com/VirajAgarwal1/lox/errorhandler"
	"github.com/VirajAgarwal1/lox/lexer/dfa"
	dfa "github.com/VirajAgarwal1/lox/lexer/dfa"
)

/*
# IDEAS

	1. The type in the stack could be a struct with the first byte being what is it (terminal, non-terminal or an operation) and the other bytes would be the actual data

	2. 1 operation could be for denoting the end of a production rule.

	3. All the operations ("star", "plus", "or" and "bracket") will be treated as seperate non-terminals without emit events. This parsing them and also keeping them out of the AST Tree.

	3. ALL the artifical non-terminals will start with the runes '69'
*/
///////////////////////////////////
///////////////////////////////////
///////////////////////////////////
/*
# PSEUDOCODE For SAX-tyle LL(1) Predictive Parser

function Parse( lexer ) {
	st <- []
	PushStartSymbol( st )

	while IsNotEmpty( stack ) {

		top <- Peek( st )
		Switch( top ) {
			case IsNonTerminal( top ) {
				if LookAhead( lexer ) In First( top ) {
					Pop( st )
					EmitEvent( 'startNode(`top`)' )
					Push( st, Reverse( ProductionRule( top, LookAhead( lexer ) ) ) )
				}
				else {
					while LookAhead( lexer ) Not In Follow( top ) {
						ConsumeToken( lexer )
					}
					EmitEvent( 'Error' )
				}
			}
			case IsTerminal( top ) {
				if LookAhead( lexer ) == top {
					Pop( st )
					EmitEvent( 'leafNode(`top`)' )
					ConsumeToken( lexer )
				}
				else {
					EmitEvent( 'Error' )
				}
			}
			case IsOperation( top ) {
				OperationHandler <- GetOperationHandler( top )
				OperationHandler( lexer )
			}
		}
	}
	if LookAhead( lexer ) == EOF {
		EmitEvent( 'endNode(`Start`)' )
	}
	else {
		EmitEvent( 'Error' )
	}
}
*/

var First = map[string][]dfa.TokenType{
	"expression": {dfa.TRUE, dfa.FALSE, dfa.IDENTIFIER, dfa.STRING},
	// ...
}

var Follow = map[string][]dfa.TokenType{
	"expression": {dfa.TRUE, dfa.FALSE, dfa.IDENTIFIER, dfa.STRING},
	// ...
}

const epsilon = dfa.TokenType("epsilon") // Temporary epsilon, made it cause I need it for Production Rules

type data_type_for_stack struct {
	isNonTerminal bool
	non_term_name string
	terminal_type dfa.TokenType
}

type token_set map[dfa.TokenType]struct{}
type syntaxRule struct {
	set token_set
	def []data_type_for_stack
}

var grammarRuleSet = map[string]([]syntaxRule){
	"expression": {
		syntaxRule{
			set: token_set{
				dfa.TRUE:       {},
				dfa.FALSE:      {},
				dfa.IDENTIFIER: {},
				dfa.STRING:     {},
			},
			def: []data_type_for_stack{
				{isNonTerminal: true, non_term_name: "term"},
				{isNonTerminal: true, non_term_name: "expression_end"},
			},
		},
	},
	// ...
}

func ProductionRule(non_terminal_name string, look_ahead_token dfa.TokenType) (error, []data_type_for_stack) {
	slice_of_rules, ok := grammarRuleSet[non_terminal_name]
	if !ok {
		return errorhandler.RetErr(fmt.Sprintf("non-terminal {%s} not found", non_terminal_name), nil), nil
	}
	for _, rule := range slice_of_rules {
		if _, ok := rule.set[look_ahead_token]; ok {
			return nil, rule.def
		}
	}
	return errorhandler.RetErr(fmt.Sprintf("no production rule found for non-terminal {%s} with look-ahead token {%s}", non_terminal_name, string(look_ahead_token)), nil), nil
}

/*


star_start
bracket_start
bracket_start
"!="
or
"=="
bracket_end
comparison
bracket_end
star_end
equality_end
star_start
bracket_start
"comma"
equality
bracket_end
star_end
comma_end
expression_end

*/
