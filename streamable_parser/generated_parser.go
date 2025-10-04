package streamable_parser

import (
	dfa "github.com/VirajAgarwal1/lox/lexer/dfa"
	utils "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/utils"
)

const StartingNonTerminal string = "expression"

var grammarRules = map[string]ProductionRule{
	"unary": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.SLASH: {},	dfa.STAR: {},	dfa.MINUS: {},	dfa.PLUS: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG: {},	dfa.MINUS: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_0"},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "primary"},
				},
			},
		},
	},	"comma": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "equality"},	{IsNonTerminal: true, Non_term_name: "999_3"},
				},
			},
		},
	},	"999_7": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG_EQUAL: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.BANG_EQUAL},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.EQUAL_EQUAL: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.EQUAL_EQUAL},
				},
			},
		},
	},	"equality": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "comparison"},	{IsNonTerminal: true, Non_term_name: "999_5"},
				},
			},
		},
	},	"999_10": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.GREATER: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.GREATER},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.GREATER_EQUAL: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.GREATER_EQUAL},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.LESS: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.LESS},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.LESS_EQUAL: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.LESS_EQUAL},
				},
			},
		},
	},	"comparison": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "term"},	{IsNonTerminal: true, Non_term_name: "999_8"},
				},
			},
		},
	},	"999_12": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.MINUS: {},	dfa.PLUS: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.MINUS: {},	dfa.PLUS: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_13"},	{IsNonTerminal: true, Non_term_name: "factor"},
				},
			},
		},
	},	"expression": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "comma"},
				},
			},
		},
	},	"999_5": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_6"},	{IsNonTerminal: true, Non_term_name: "999_5"},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					utils.Epsilon: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: utils.Epsilon},
				},
			},
		},
	},	"999_8": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_9"},	{IsNonTerminal: true, Non_term_name: "999_8"},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					utils.Epsilon: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: utils.Epsilon},
				},
			},
		},
	},	"999_13": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.MINUS: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.MINUS},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.PLUS: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.PLUS},
				},
			},
		},
	},	"999_11": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.MINUS: {},	dfa.PLUS: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_12"},	{IsNonTerminal: true, Non_term_name: "999_11"},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					utils.Epsilon: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: utils.Epsilon},
				},
			},
		},
	},	"term": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "factor"},	{IsNonTerminal: true, Non_term_name: "999_11"},
				},
			},
		},
	},	"999_15": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.SLASH: {},	dfa.STAR: {},	dfa.MINUS: {},	dfa.PLUS: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.SLASH: {},	dfa.STAR: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_16"},	{IsNonTerminal: true, Non_term_name: "unary"},
				},
			},
		},
	},	"999_1": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.BANG},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.MINUS: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.MINUS},
				},
			},
		},
	},	"999_0": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.SLASH: {},	dfa.STAR: {},	dfa.MINUS: {},	dfa.PLUS: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG: {},	dfa.MINUS: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_1"},	{IsNonTerminal: true, Non_term_name: "unary"},
				},
			},
		},
	},	"999_2": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.SLASH: {},	dfa.STAR: {},	dfa.MINUS: {},	dfa.PLUS: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.LEFT_PAREN: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.LEFT_PAREN},	{IsNonTerminal: true, Non_term_name: "expression"},	{IsNonTerminal: false, Terminal_type: dfa.RIGHT_PAREN},
				},
			},
		},
	},	"999_6": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_7"},	{IsNonTerminal: true, Non_term_name: "comparison"},
				},
			},
		},
	},	"999_14": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.MINUS: {},	dfa.PLUS: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.SLASH: {},	dfa.STAR: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_15"},	{IsNonTerminal: true, Non_term_name: "999_14"},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					utils.Epsilon: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: utils.Epsilon},
				},
			},
		},
	},	"primary": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.SLASH: {},	dfa.STAR: {},	dfa.MINUS: {},	dfa.PLUS: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.IDENTIFIER: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.IDENTIFIER},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.NUMBER: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.NUMBER},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.STRING: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.STRING},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.TRUE: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.TRUE},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.FALSE: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.FALSE},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.NIL: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.NIL},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.LEFT_PAREN: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_2"},
				},
			},
		},
	},	"999_4": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.COMMA: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.COMMA},	{IsNonTerminal: true, Non_term_name: "equality"},
				},
			},
		},
	},	"999_3": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.COMMA: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_4"},	{IsNonTerminal: true, Non_term_name: "999_3"},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					utils.Epsilon: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: utils.Epsilon},
				},
			},
		},
	},	"999_9": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "999_10"},	{IsNonTerminal: true, Non_term_name: "term"},
				},
			},
		},
	},	"999_16": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.SLASH: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.SLASH},
				},
			},	{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.STAR: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: false, Terminal_type: dfa.STAR},
				},
			},
		},
	},	"factor": {
		FollowSet: map[dfa.TokenType]struct{}{
			dfa.EOF: {},	dfa.MINUS: {},	dfa.PLUS: {},	dfa.GREATER: {},	dfa.GREATER_EQUAL: {},	dfa.LESS: {},	dfa.LESS_EQUAL: {},	dfa.BANG_EQUAL: {},	dfa.EQUAL_EQUAL: {},	dfa.COMMA: {},	dfa.RIGHT_PAREN: {},
		},
		Sequences: []GrammarSequence{
			{
				FirstSet: map[dfa.TokenType]struct{}{
					dfa.BANG: {},	dfa.MINUS: {},	dfa.IDENTIFIER: {},	dfa.NUMBER: {},	dfa.STRING: {},	dfa.TRUE: {},	dfa.FALSE: {},	dfa.NIL: {},	dfa.LEFT_PAREN: {},
				},
				Elements: []utils.Grammar_element{
					{IsNonTerminal: true, Non_term_name: "unary"},	{IsNonTerminal: true, Non_term_name: "999_14"},
				},
			},
		},
	},
}
