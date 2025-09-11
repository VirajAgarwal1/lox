package parser_generator

import "strconv"

/*
# NOTEs

1. Artificial non-terminals will always start with '999_'. Example, '999_1'. It starts with a number because non-terminals in grammar definition cannot start with numbers
*/

const artificial_non_term_prefix string = "999_"

var artificial_non_terminal_counter int = 0
var sanitized_grammar = map[string]([][]Grammar_element){}

func process_term(term Generic_grammar_term) Grammar_element {

	switch term.Get_grammar_term_type() {
	case "terminal":
		terminal_term := term.(*Terminal)
		return Grammar_element{
			IsNonTerminal: false,
			Terminal_type: string_to_token_type[string(terminal_term.Content)],
		}

	case "non_terminal":
		non_terminal_term := term.(*Non_terminal)
		return Grammar_element{
			IsNonTerminal: true,
			Non_term_name: non_terminal_term.Name,
		}

	case "star":
		star_term := term.(*Star)
		new_artificial_non_term_name := artificial_non_term_prefix + strconv.Itoa(artificial_non_terminal_counter)
		artificial_non_terminal_counter++

		// Get the production for the new artifical non-terminal
		sanitized_grammar[new_artificial_non_term_name] = [][]Grammar_element{
			{
				process_term(star_term.Content), {IsNonTerminal: true, Non_term_name: new_artificial_non_term_name},
			},
			{
				{IsNonTerminal: false, Terminal_type: Epsilon},
			},
		}

		return Grammar_element{
			IsNonTerminal: true,
			Non_term_name: new_artificial_non_term_name,
		}

	case "plus":
		plus_term := term.(*Plus)
		new_artificial_non_term_name := artificial_non_term_prefix + strconv.Itoa(artificial_non_terminal_counter)
		artificial_non_terminal_counter++

		// Get the production for the new artifical non-terminal
		sanitized_grammar[new_artificial_non_term_name] = [][]Grammar_element{
			{
				process_term(plus_term.Content), {IsNonTerminal: true, Non_term_name: new_artificial_non_term_name},
			},
			{
				process_term(plus_term.Content),
			},
		}

		return Grammar_element{
			IsNonTerminal: true,
			Non_term_name: new_artificial_non_term_name,
		}

	case "bracket":
		bracket_term := term.(*Bracket)
		new_artificial_non_term_name := artificial_non_term_prefix + strconv.Itoa(artificial_non_terminal_counter)
		artificial_non_terminal_counter++

		// Get the production for the new artifical non-terminal
		sanitized_grammar[new_artificial_non_term_name] = process_sequence(bracket_term.Contents)

		return Grammar_element{
			IsNonTerminal: true,
			Non_term_name: new_artificial_non_term_name,
		}
	}

	return Grammar_element{}
}
func process_or(choices [][]Generic_grammar_term) [][]Grammar_element {
	output := make([][]Grammar_element, 0, len(choices))
	for _, path := range choices {
		if len(path) < 1 {
			continue
		}
		if len(path) == 1 {
			output = append(output, []Grammar_element{process_term(path[0])})
		} else {
			new_artificial_non_term_name := artificial_non_term_prefix + strconv.Itoa(artificial_non_terminal_counter)
			artificial_non_terminal_counter++
			output = append(output, []Grammar_element{{
				IsNonTerminal: true,
				Non_term_name: new_artificial_non_term_name,
			}})
			sanitized_grammar[new_artificial_non_term_name] = process_sequence(path)
		}
	}
	return output
}
func process_sequence(sequence []Generic_grammar_term) [][]Grammar_element {

	or_positions := detect_or_in_sequence(sequence)
	if len(or_positions) != 0 {
		choices := [][]Generic_grammar_term{}
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

	output := [][]Grammar_element{{}}
	for _, term := range sequence {
		output[0] = append(output[0], process_term(term))
	}

	return output
}

func EbnfToBnfConverter(grammar_rules map[Non_terminal]([]Generic_grammar_term)) map[string]([][]Grammar_element) {
	// The first slice is for incorporating 'or' and the internal slices for the actual definition

	sanitized_grammar = map[string]([][]Grammar_element){}

	for non_term, def := range grammar_rules {
		sanitized_grammar[non_term.Name] = process_sequence(def)
	}

	return sanitized_grammar
}
