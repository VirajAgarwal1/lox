package ebnf_to_bnf

import (
	"strconv"

	gfp "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/grammar_file_parser"
	utils "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/utils"
)

/*
# NOTEs

1. Artificial non-terminals will always start with '999_'. Example, '999_1'. It starts with a number because non-terminals in grammar definition cannot start with numbers
*/

const Artificial_non_term_prefix string = "999_"

var artificial_non_terminal_counter int = 0
var bnf_grammar = map[string]([][]utils.Grammar_element){}

func process_term(term gfp.Generic_grammar_term) utils.Grammar_element {

	switch term.Get_grammar_term_type() {
	case "terminal":
		terminal_term := term.(*gfp.Terminal)
		return utils.Grammar_element{
			IsNonTerminal: false,
			Terminal_type: utils.String_to_token[string(terminal_term.Content)],
		}

	case "non_terminal":
		non_terminal_term := term.(*gfp.Non_terminal)
		return utils.Grammar_element{
			IsNonTerminal: true,
			Non_term_name: non_terminal_term.Name,
		}

	case "star":
		star_term := term.(*gfp.Star)
		new_artificial_non_term_name := Artificial_non_term_prefix + strconv.Itoa(artificial_non_terminal_counter)
		artificial_non_terminal_counter++

		// Get the production for the new artifical non-terminal
		bnf_grammar[new_artificial_non_term_name] = [][]utils.Grammar_element{
			{
				process_term(star_term.Content), {IsNonTerminal: true, Non_term_name: new_artificial_non_term_name},
			},
			{
				{IsNonTerminal: false, Terminal_type: utils.Epsilon},
			},
		}

		return utils.Grammar_element{
			IsNonTerminal: true,
			Non_term_name: new_artificial_non_term_name,
		}

	case "plus":
		plus_term := term.(*gfp.Plus)
		new_artificial_non_term_name := Artificial_non_term_prefix + strconv.Itoa(artificial_non_terminal_counter)
		artificial_non_terminal_counter++

		// Get the production for the new artifical non-terminal
		bnf_grammar[new_artificial_non_term_name] = [][]utils.Grammar_element{
			{
				process_term(plus_term.Content), {IsNonTerminal: true, Non_term_name: new_artificial_non_term_name},
			},
			{
				process_term(plus_term.Content),
			},
		}

		return utils.Grammar_element{
			IsNonTerminal: true,
			Non_term_name: new_artificial_non_term_name,
		}

	case "bracket":
		bracket_term := term.(*gfp.Bracket)
		new_artificial_non_term_name := Artificial_non_term_prefix + strconv.Itoa(artificial_non_terminal_counter)
		artificial_non_terminal_counter++

		// Get the production for the new artifical non-terminal
		bnf_grammar[new_artificial_non_term_name] = process_sequence(bracket_term.Contents)

		return utils.Grammar_element{
			IsNonTerminal: true,
			Non_term_name: new_artificial_non_term_name,
		}
	}

	return utils.Grammar_element{}
}
func process_or(choices [][]gfp.Generic_grammar_term) [][]utils.Grammar_element {
	output := make([][]utils.Grammar_element, 0, len(choices))
	for _, path := range choices {
		if len(path) < 1 {
			continue
		}
		if len(path) == 1 {
			output = append(output, []utils.Grammar_element{process_term(path[0])})
		} else {
			new_artificial_non_term_name := Artificial_non_term_prefix + strconv.Itoa(artificial_non_terminal_counter)
			artificial_non_terminal_counter++
			output = append(output, []utils.Grammar_element{{
				IsNonTerminal: true,
				Non_term_name: new_artificial_non_term_name,
			}})
			bnf_grammar[new_artificial_non_term_name] = process_sequence(path)
		}
	}
	return output
}
func process_sequence(sequence []gfp.Generic_grammar_term) [][]utils.Grammar_element {

	or_positions := utils.Detect_or_in_sequence(sequence)
	if len(or_positions) != 0 {
		choices := [][]gfp.Generic_grammar_term{}
		for i, pos := range or_positions {
			start := uint32(0)
			end := pos
			if i > 0 {
				start = or_positions[i-1] + 1
			}
			choices = append(choices, sequence[start:end])
		}
		choices = append(choices, sequence[or_positions[len(or_positions)-1]+1:])
		output := process_or(choices)
		return output
	}

	output := [][]utils.Grammar_element{{}}
	for _, term := range sequence {
		output[0] = append(output[0], process_term(term))
	}

	return output
}

func EbnfToBnfConverter(ebnf_grammar map[gfp.Non_terminal]([]gfp.Generic_grammar_term)) map[string]([][]utils.Grammar_element) {
	// The first slice is for incorporating 'or' and the internal slices for the actual definition

	bnf_grammar = map[string]([][]utils.Grammar_element){}

	for non_term, def := range ebnf_grammar {
		bnf_grammar[non_term.Name] = process_sequence(def)
	}

	return bnf_grammar
}
