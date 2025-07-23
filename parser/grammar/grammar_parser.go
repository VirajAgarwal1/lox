/*
This file is responsible for generating the code for the Lox grammar. This is needed because manually writing code for the whole grammar wil become tedious
*/

package grammar

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/VirajAgarwal1/lox/errorhandler"
	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/lexer/dfa"
)

var stack stack_type
var GrammarRules = make(map[non_terminal]([]generic_grammar_term))

// This interface actually is used to refer to all the types' pointers.
type generic_grammar_term interface {
	get_grammar_term_type() string
}
type terminal struct {
	content []rune
}
type non_terminal struct {
	name string
}
type or struct {
}
type star struct {
	content generic_grammar_term
}
type plus struct {
	content generic_grammar_term
}
type bracket struct {
	contents []generic_grammar_term
	is_left  bool
}
type stack_type []generic_grammar_term

func (t *terminal) get_grammar_term_type() string {
	return "terminal"
}
func (t *non_terminal) get_grammar_term_type() string {
	return "non_terminal"
}
func (t *or) get_grammar_term_type() string {
	return "or"
}
func (t *star) get_grammar_term_type() string {
	return "star"
}
func (t *plus) get_grammar_term_type() string {
	return "plus"
}
func (t *bracket) get_grammar_term_type() string {
	return "bracket"
}
func (st *stack_type) add(elem generic_grammar_term) {
	*st = append(*st, elem)
}
func (st *stack_type) pop() generic_grammar_term {
	if len(*st) == 0 {
		return nil
	}
	out := (*st)[len(*st)-1]
	*st = (*st)[:(len(*st) - 1)]
	return out
}
func (st *stack_type) peek() generic_grammar_term {
	return (*st)[len(*st)-1]
}

func processGrammarDefinition(scanner lexer.LexicalAnalyzer) {

	i := -1
	current_non_terminal := non_terminal{}

	for {
		token, err := scanner.ReadToken()
		if err != nil && err != io.EOF {
			fmt.Print(errorhandler.RetErr("Inavlid Grammar: token not recognized", nil))
			return
		}
		if err == io.EOF {
			if current_non_terminal.name != "" {
				var new_non_terminal_def []generic_grammar_term
				for i := 0; i < len(stack); i++ {
					new_non_terminal_def = append(new_non_terminal_def, stack[i])
				}
				GrammarRules[current_non_terminal] = new_non_terminal_def
			}
			stack = nil
			return
		}

		fmt.Println(token)

		if token.TypeOfToken == dfa.WHITESPACE {
			continue
		}
		if token.TypeOfToken == dfa.COMMENT {
			continue
		}
		if token.TypeOfToken == dfa.NEWLINE {
			if current_non_terminal.name != "" {
				var new_non_terminal_def []generic_grammar_term
				for i := 0; i < len(stack); i++ {
					new_non_terminal_def = append(new_non_terminal_def, stack[i])
				}
				GrammarRules[current_non_terminal] = new_non_terminal_def
			}
			stack = stack[:0] // Clear out the stack
			i = -1
			continue
		}

		i++

		switch i {
		case 0:
			// We are the starting of a new line
			if token.TypeOfToken != dfa.IDENTIFIER {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: left expression missing", nil))
				return
			}
			current_non_terminal.name = string(token.Lexemme)
			continue
		case 1:
			if token.TypeOfToken != dfa.MINUS {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: separator (->) is missing", nil))
				return
			}
			continue
		case 2:
			if token.TypeOfToken != dfa.GREATER {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: separator (->) is missing", nil))
				return
			}
			continue
		}

		if token.TypeOfToken == dfa.IDENTIFIER {
			arg := non_terminal{}
			arg.name = string(token.Lexemme)
			stack.add(&arg)
			continue
		}
		if token.TypeOfToken == dfa.LEFT_PAREN {
			open_bracket := bracket{}
			open_bracket.is_left = true
			stack.add(&open_bracket)
			continue
		}
		if token.TypeOfToken == dfa.RIGHT_PAREN {
			i := len(stack) - 1
			for i >= 0 {
				if stack[i].get_grammar_term_type() == "bracket" {
					if stack[i].(*bracket).is_left {
						break
					}
				}
				i--
			}
			if i < 0 {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: no matching left bracket found for the right bracket", nil))
				return
			}
			// Take all the elems out from the stack from this index and place them in the new bracket
			close_bracket := bracket{}
			close_bracket.is_left = false
			for j := i + 1; j < len(stack); j++ {
				close_bracket.contents = append(close_bracket.contents, stack[j])
			}
			stack = stack[:i] // Remove all the other things from the stack
			stack.add(&close_bracket)
			continue
		}
		if token.TypeOfToken == dfa.STAR {
			if len(stack) < 1 {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: '*' needs an element before itself to function", nil))
				return
			}
			prev_elem := stack.peek()
			if prev_elem.get_grammar_term_type() == "bracket" && prev_elem.(*bracket).is_left {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: '*' cannot have a open bracket right before itself", nil))
				return
			}
			if prev_elem.get_grammar_term_type() == "or" {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: '*' cannot have the 'or` operator right before itself", nil))
				return
			}
			prev_elem = stack.pop()
			new_star := star{}
			new_star.content = prev_elem
			stack.add(&new_star)
			continue
		}
		if token.TypeOfToken == dfa.PLUS {
			if len(stack) < 1 {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: '+' needs an element before itself to function", nil))
				return
			}
			prev_elem := stack.peek()
			if prev_elem.get_grammar_term_type() == "bracket" && prev_elem.(*bracket).is_left {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: '+' cannot have a open bracket right before itself", nil))
				return
			}
			if prev_elem.get_grammar_term_type() == "or" {
				fmt.Print(errorhandler.RetErr("Inavlid Grammar: '+' cannot have the 'or` operator right before itself", nil))
				return
			}
			prev_elem = stack.pop()
			new_plus := plus{}
			new_plus.content = prev_elem
			stack.add(&new_plus)
			continue
		}
		if token.TypeOfToken == dfa.OR {
			new_or := or{}
			stack.add(&new_or)
			continue
		}
		if token.TypeOfToken == dfa.STRING {
			new_terminal := terminal{}
			new_terminal.content = token.Lexemme[1 : len(token.Lexemme)-1] // Excluding the apostrophies from the sides
			stack.add(&new_terminal)
			continue
		}
	}
}

func Grammar_parser() {
	file_reader, err := os.Open("parser/lox.grammar")
	if err != nil {
		fmt.Print(err)
		return
	}
	buf_file_reader := bufio.NewReader(file_reader)
	scanner := lexer.LexicalAnalyzer{}
	scanner.Initialize(buf_file_reader)
	processGrammarDefinition(scanner)
}
