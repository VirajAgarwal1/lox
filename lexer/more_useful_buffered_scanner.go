package lexer

import (
	"bufio"

	"github.com/VirajAgarwal1/lox/errorhandler"
)

type Checkpoint int
type lexerResult struct {
	tok *Token
	err error
}

type indexManager struct {
	currentIndex int
}

func (im *indexManager) initialize() {
	im.currentIndex = -1
}
func (im *indexManager) makeCheckpoint() Checkpoint {
	new_checkpoint := Checkpoint(im.currentIndex)
	return new_checkpoint
}
func (im *indexManager) rollbackTo(chk Checkpoint) {
	im.currentIndex = int(chk)
}

type BufferedLexer struct {
	scanner *LexicalAnalyzer
	buffer  []lexerResult
	index   indexManager
}

func (buf_lex *BufferedLexer) Initialize(source *bufio.Reader, max_bufer_capacity uint32) {
	scanner := LexicalAnalyzer{}
	scanner.Initialize(source)
	buf_lex.scanner = &scanner

	buf_lex.buffer = make([]lexerResult, 0, max_bufer_capacity)

	buf_lex.index.initialize()
}
func (buf_lex *BufferedLexer) currentIndexPointsToLegitToken() bool {
	return buf_lex.index.currentIndex >= 0 && buf_lex.index.currentIndex < len(buf_lex.buffer)
}
func (buf_lex *BufferedLexer) MakeCheckpoint() Checkpoint {
	return buf_lex.index.makeCheckpoint()
}
func (buf_lex *BufferedLexer) RollbackTo(chk Checkpoint) {
	buf_lex.index.rollbackTo(chk)
}
func (buf_lex *BufferedLexer) ReadToken() (*Token, error) {
	buf_lex.index.currentIndex += 1
	// If the current index points to a place in buffer which actually correspond to a token, If the it does then you return it
	if buf_lex.currentIndexPointsToLegitToken() {
		return buf_lex.buffer[buf_lex.index.currentIndex].tok, buf_lex.buffer[buf_lex.index.currentIndex].err
	}
	// , then you read a new token from the lexer, but only if the buffer has space to accomodate the new token
	if len(buf_lex.buffer) == cap(buf_lex.buffer) {
		return nil, errorhandler.RetErr("Lexer Error: Input Buffer Overflow", nil)
	}
	tok, err := buf_lex.scanner.ReadToken()
	buf_lex.buffer = append(buf_lex.buffer, lexerResult{tok, err})
	return tok, err
}
func (buf_lex *BufferedLexer) ClearBuffer() {
	buf_lex.index.currentIndex = -1
	buf_lex.buffer = buf_lex.buffer[:0]
}
