package lexer

import (
	"bufio"
	"fmt"
)

type BufferedLexicalAnalyzer struct {
	buffer [3]struct {
		t   *Token
		err error
	}
	prev_i int8 // These 3 indices will move around to mark what is
	cur_i  int8 // the current, previuous and the next token (along
	next_i int8 // with its error)

	scanner *LexicalAnalyzer
}

func (b *BufferedLexicalAnalyzer) Initialize(source *bufio.Reader) {
	scanner := LexicalAnalyzer{}
	scanner.Initialize(source)
	b.scanner = &scanner

	b.prev_i = 0
	b.cur_i = 1
	b.next_i = 2

	b.buffer[b.prev_i].t = nil
	b.buffer[b.prev_i].err = fmt.Errorf("no previous token exists")

	t, err := b.scanner.ReadToken()
	b.buffer[b.cur_i].t = t
	b.buffer[b.cur_i].err = err

	t, err = b.scanner.ReadToken()
	b.buffer[b.next_i].t = t
	b.buffer[b.next_i].err = err
}
func (b *BufferedLexicalAnalyzer) ConsumeOneToken() (*Token, error) {

	out_token := b.buffer[b.cur_i].t
	out_error := b.buffer[b.cur_i].err

	t, err := b.scanner.ReadToken()
	b.buffer[b.prev_i].t = t
	b.buffer[b.prev_i].err = err

	old_prev := b.prev_i
	b.prev_i = b.cur_i
	b.cur_i = b.next_i
	b.next_i = old_prev

	return out_token, out_error
}
func (b *BufferedLexicalAnalyzer) NextTokenWithoutConsume() (*Token, error) {
	return b.buffer[b.next_i].t, b.buffer[b.next_i].err
}
func (b *BufferedLexicalAnalyzer) CurrentTokenWithoutConsume() (*Token, error) {
	return b.buffer[b.cur_i].t, b.buffer[b.cur_i].err
}
func (b *BufferedLexicalAnalyzer) PreviousTokenWithoutConsume() (*Token, error) {
	return b.buffer[b.prev_i].t, b.buffer[b.prev_i].err
}
