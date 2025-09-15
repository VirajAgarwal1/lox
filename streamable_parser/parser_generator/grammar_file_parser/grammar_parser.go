package grammar_parser

import (
	"io"

	"github.com/VirajAgarwal1/lox/errorhandler"
	"github.com/VirajAgarwal1/lox/lexer"
	"github.com/VirajAgarwal1/lox/lexer/dfa"
)

// This interface actually is used to refer to all the types' pointers.
type Generic_grammar_term interface {
	Get_grammar_term_type() string
	// Follow_set() []dfa.TokenType
}

type Terminal struct {
	Content []rune
}
type Non_terminal struct {
	Name string
}
type Or struct {
}
type Star struct {
	Content Generic_grammar_term
}
type Plus struct {
	Content Generic_grammar_term
}
type Bracket struct {
	Contents []Generic_grammar_term
	Is_left  bool
}

func (t *Terminal) Get_grammar_term_type() string {
	return "terminal"
}
func (t *Non_terminal) Get_grammar_term_type() string {
	return "non_terminal"
}
func (t *Or) Get_grammar_term_type() string {
	return "or"
}
func (t *Star) Get_grammar_term_type() string {
	return "star"
}
func (t *Plus) Get_grammar_term_type() string {
	return "plus"
}
func (t *Bracket) Get_grammar_term_type() string {
	return "bracket"
}

type Stack_type []Generic_grammar_term

func (st *Stack_type) add(elem Generic_grammar_term) {
	*st = append(*st, elem)
}
func (st *Stack_type) pop() Generic_grammar_term {
	if len(*st) == 0 {
		return nil
	}
	out := (*st)[len(*st)-1]
	*st = (*st)[:(len(*st) - 1)]
	return out
}
func (st *Stack_type) peek() Generic_grammar_term {
	return (*st)[len(*st)-1]
}

func ProcessGrammarDefinition(scanner *lexer.LexicalAnalyzer) (map[Non_terminal]([]Generic_grammar_term), error) {

	i := -1
	current_non_terminal := Non_terminal{}
	var stack Stack_type
	var GrammarRules = make(map[Non_terminal]([]Generic_grammar_term))

	for {
		token, err := scanner.ReadToken()
		if err != nil && err != io.EOF {
			return GrammarRules, errorhandler.RetErr("Inavlid Grammar: token not recognized", err)
		}
		if err == io.EOF {
			if current_non_terminal.Name != "" {
				var new_non_terminal_def []Generic_grammar_term
				for i := 0; i < len(stack); i++ {
					new_non_terminal_def = append(new_non_terminal_def, stack[i])
				}
				GrammarRules[current_non_terminal] = new_non_terminal_def
			}
			stack = nil
			return GrammarRules, err
		}

		if token.TypeOfToken == dfa.WHITESPACE {
			continue
		}
		if token.TypeOfToken == dfa.COMMENT {
			continue
		}
		if token.TypeOfToken == dfa.NEWLINE {
			if current_non_terminal.Name != "" {
				var new_non_terminal_def []Generic_grammar_term
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
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: left expression missing", nil)
			}
			current_non_terminal.Name = string(token.Lexemme)
			continue
		case 1:
			if token.TypeOfToken != dfa.MINUS {
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: separator (->) is missing", nil)
			}
			continue
		case 2:
			if token.TypeOfToken != dfa.GREATER {
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: separator (->) is missing", nil)
			}
			continue
		}

		if token.TypeOfToken == dfa.IDENTIFIER {
			arg := Non_terminal{}
			arg.Name = string(token.Lexemme)
			stack.add(&arg)
			continue
		}
		if token.TypeOfToken == dfa.LEFT_PAREN {
			open_bracket := Bracket{}
			open_bracket.Is_left = true
			stack.add(&open_bracket)
			continue
		}
		if token.TypeOfToken == dfa.RIGHT_PAREN {
			i := len(stack) - 1
			for i >= 0 {
				if stack[i].Get_grammar_term_type() == "bracket" {
					if stack[i].(*Bracket).Is_left {
						break
					}
				}
				i--
			}
			if i < 0 {
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: no matching left bracket found for the right bracket", nil)
			}
			// Take all the elems out from the stack from this index and place them in the new bracket
			close_bracket := Bracket{}
			close_bracket.Is_left = false
			for j := i + 1; j < len(stack); j++ {
				close_bracket.Contents = append(close_bracket.Contents, stack[j])
			}
			stack = stack[:i] // Remove all the other things from the stack
			stack.add(&close_bracket)
			continue
		}
		if token.TypeOfToken == dfa.STAR {
			if len(stack) < 1 {
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: '*' needs an element before itself to function", nil)
			}
			prev_elem := stack.peek()
			if prev_elem.Get_grammar_term_type() == "bracket" && prev_elem.(*Bracket).Is_left {
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: '*' cannot have a open bracket right before itself", nil)
			}
			if prev_elem.Get_grammar_term_type() == "or" {
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: '*' cannot have the 'or` operator right before itself", nil)
			}
			prev_elem = stack.pop()
			new_star := Star{}
			new_star.Content = prev_elem
			stack.add(&new_star)
			continue
		}
		if token.TypeOfToken == dfa.PLUS {
			if len(stack) < 1 {
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: '+' needs an element before itself to function", nil)
			}
			prev_elem := stack.peek()
			if prev_elem.Get_grammar_term_type() == "bracket" && prev_elem.(*Bracket).Is_left {
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: '+' cannot have a open bracket right before itself", nil)
			}
			if prev_elem.Get_grammar_term_type() == "or" {
				return GrammarRules, errorhandler.RetErr("Inavlid Grammar: '+' cannot have the 'or` operator right before itself", nil)
			}
			prev_elem = stack.pop()
			new_plus := Plus{}
			new_plus.Content = prev_elem
			stack.add(&new_plus)
			continue
		}
		if token.TypeOfToken == dfa.OR {
			new_or := Or{}
			stack.add(&new_or)
			continue
		}
		if token.TypeOfToken == dfa.STRING {
			new_terminal := Terminal{}
			new_terminal.Content = token.Lexemme[1 : len(token.Lexemme)-1] // Excluding the apostrophies from the sides
			stack.add(&new_terminal)
			continue
		}
	}
}
