package lexer

import (
	"bufio"
	"fmt"
	"io"

	errorhandler "github.com/VirajAgarwal1/lox/errorhandler"
	dfa "github.com/VirajAgarwal1/lox/lexer/dfa"
)

/*
	The order of conversion is this:
		Bytes -> Runes -> Lexemes -> Tokens

	*NOTE*:
	We need this lastInput fields (in the LexicalAnalyzer type) because when the dfas detect a valid token, the scanner needs to consume one more rune, which will make all the dfas return invalid. And then the token is returned. But, this process consumes a rune whiich didn't belong to the previous token, and so we save it and reuse it in the next time the token reader function is called. Yes, this does mean that the procesing of a starting rune of each token is done TWICE. The fix is easy by adding the processing info in the struct as well... But, if I add it right now that can mess up the code a lot.

	^ This might also potentially need changing some things in the dfa package, because I have not exposed a way to start each DFA on some particular state or even a way to save all the dfas state in an easy to do manner.

	So, the twice execution problem is here to stay for now, But I'll still put a todo. ->
	TODO: Fix the problem inefficiency of the 1st rune of a tokken going through all the lexemmes DFAs twice.

	TODO: Setup a new struct which will have the lastInput properties along with the dfa tokens
*/

// type lastReadTokenRun
type inputRunePosition struct {
	lineNum        uint32
	lineOffset     uint32
	prevLineOffset uint32
}
type Token struct {
	TypeOfToken dfa.TokenType
	Lexemme     []rune
	Line        uint32
	Offset      uint32
}
type LexicalAnalyzer struct {
	source              *bufio.Reader
	stateManger         *dfa.DFAStatesManager
	currentInput        rune
	currentPos          *inputRunePosition
	lexemme             []rune
	lastInputExists     bool
	sustainCurrentInput bool
}

//	func (tokenPos *inputRunePosition) stepBack() {
//		if tokenPos.lineOffset == 0 {
//			tokenPos.lineNum--
//			tokenPos.lineOffset = tokenPos.prevLineOffset
//			return
//		}
//		tokenPos.lineOffset--
//	}
//
//	func (tokenPos *inputRunePosition) getPos() (lineNum uint32, lineOffset uint32) {
//		lineNum = tokenPos.lineNum
//		lineOffset = tokenPos.lineOffset
//		return
//	}
func (tokenPos *inputRunePosition) initialize() {
	tokenPos.lineNum = 0
	tokenPos.lineOffset = 0
	tokenPos.prevLineOffset = 0
}
func (tokenPos *inputRunePosition) step(input rune) {
	if input == '\n' {
		tokenPos.lineNum++
		tokenPos.prevLineOffset = tokenPos.lineOffset
		tokenPos.lineOffset = 0
		return
	}
	tokenPos.lineNum++
	tokenPos.lineOffset++
}
func (tokenPos *inputRunePosition) reset() {
	tokenPos.initialize()
}

func (tok *Token) SetTokenProperties(_type dfa.TokenType, _lineNum uint32, _offset uint32, _lexemme []rune) {
	tok.TypeOfToken = _type
	tok.Line = _lineNum
	tok.Offset = _offset
	tok.Lexemme = _lexemme
}
func (tok *Token) ToString() string {
	return fmt.Sprintf("|%d|%d| [%s]Token -> `%s`", tok.Line, tok.Offset, string(tok.TypeOfToken), string(tok.Lexemme))
}

func (scanner *LexicalAnalyzer) Initialize(source *bufio.Reader) {
	scanner.stateManger = &dfa.DFAStatesManager{}
	scanner.currentPos = &inputRunePosition{}

	scanner.source = source
	scanner.stateManger.Initialize()
	scanner.currentPos.initialize()
	scanner.lexemme = nil
	scanner.lastInputExists = false
	scanner.sustainCurrentInput = true
}
func (scanner *LexicalAnalyzer) Reset() {
	scanner.source = nil
	scanner.stateManger.Reset()
	scanner.currentPos.reset()
	scanner.lexemme = nil
	scanner.lastInputExists = false
	scanner.sustainCurrentInput = true
}

func (scanner *LexicalAnalyzer) prepareForNextToken() {
	if scanner.sustainCurrentInput {
		scanner.stateManger.Reset()
		scanner.lastInputExists = true
		scanner.lexemme = scanner.lexemme[:1]
		scanner.lexemme[0] = scanner.currentInput
		return
	}
	scanner.stateManger.Reset()
	scanner.lastInputExists = false
	scanner.lexemme = nil
}
func (scanner *LexicalAnalyzer) ReadToken() (*Token, error) {
	/*
		This function reades one rune at a time from the source reader and returns 1 token at a time in return. It follows `maximal munching` methodlogy for settling tie between 2 valid token dfas being satisfied. And if both DFAs end up having the same token length then the token which is written later in the `dfa.TokensList` is given higher priority and is returned
	*/
	returnToken := Token{}
	var err error
	tokenStartingLine := scanner.currentPos.lineNum
	tokenStartingLineOffset := scanner.currentPos.lineOffset

	defer scanner.prepareForNextToken()

	for i := 0; ; i++ {
		// Read one rune from the source
		if !scanner.lastInputExists {
			scanner.currentInput, _, err = scanner.source.ReadRune()
			if err != nil && err != io.EOF {
				scanner.sustainCurrentInput = false
				return &returnToken, errorhandler.RetErr("", err)
			}
			if err == io.EOF {
				scanner.sustainCurrentInput = false
				// This is the end of the file so just wrap up your findings and return
				if i == 0 {
					// Users should run this function to get the EOF token even if internally in this function the EOF was reached in the last run
					returnToken.SetTokenProperties(
						dfa.EOF,
						tokenStartingLine,
						tokenStartingLineOffset,
						[]rune(string(dfa.EOF)),
					)
					// TODO: See if this lineNum and lineOffset is okay for not returning the eof
					return &returnToken, io.EOF
				}
				// Check if any valid token was found in the last iteration, if yes then we report it as our found token
				if scanner.stateManger.PreviousLoopDfaResults.IsAnyValidToken {
					returnToken.SetTokenProperties(
						scanner.stateManger.PreviousLoopDfaResults.ValidToken,
						tokenStartingLine,
						tokenStartingLineOffset,
						scanner.lexemme,
					)
					return &returnToken, nil
				}
				// If any intermediates were there then we will use those for error reporting
				if scanner.stateManger.PreviousLoopDfaResults.IsAnyIntermediateToken {
					return &returnToken, errorhandler.RetErr(
						fmt.Sprintf(
							"TokenError: invalid token found at line %v at offset %v, most resembling token type was %v",
							tokenStartingLine,
							tokenStartingLineOffset,
							string(scanner.stateManger.PreviousLoopDfaResults.IntermediateToken),
						),
						nil,
					)
				}
				// I do not expect the code to reach this line. Because, to reach here the last iteration of this function's loop would have to have all Invalid tokens, and still decide to continue parsing. Which should'nt happen.
				// Return error and not EOF token here, users can get the EOF token in the next run despite the error
				return &returnToken, errorhandler.RetErr(
					fmt.Sprintf("TokenError: invalid token found at line %v at offset %v", scanner.currentPos.lineNum, scanner.currentPos.lineOffset-uint32(len(scanner.lexemme))),
					nil,
				)
			}
			// Record the offsets and the lineNums in the scanner
			scanner.currentPos.step(scanner.currentInput)
		}
		scanner.lexemme = append(scanner.lexemme, scanner.currentInput)

		// # The part where the actual stepping in the DFAs is taking place
		scanner.stateManger.Step(scanner.currentInput)

		// Stop iterating if all the DFAs are yielding INVALID
		if scanner.stateManger.CurrentLoopDfaResults.AreAllInvalid() {
			scanner.sustainCurrentInput = true
			// Check if any valid token was found in the last iteration, if yes then we report it as our found token
			if scanner.stateManger.PreviousLoopDfaResults.IsAnyValidToken {
				returnToken.SetTokenProperties(
					scanner.stateManger.PreviousLoopDfaResults.ValidToken,
					tokenStartingLine,
					tokenStartingLineOffset,
					scanner.lexemme[:len(scanner.lexemme)-1], // We remove the last rune which was not valid
				)
				return &returnToken, nil
			}
			// If any intermediates were there then we will use those for error reporting
			if scanner.stateManger.PreviousLoopDfaResults.IsAnyIntermediateToken {
				return &returnToken, errorhandler.RetErr(
					fmt.Sprintf(
						"TokenError: invalid token found at line %v at offset %v, most resembling token type was %v",
						tokenStartingLine,
						tokenStartingLineOffset,
						string(scanner.stateManger.PreviousLoopDfaResults.IntermediateToken),
					),
					nil,
				)
			}
			// The Program Counter can only get here if this is the 1st iteration and the very 1st rune did not satisfay any of the token types' dfa
			scanner.lastInputExists = false
			return &returnToken, errorhandler.RetErr(
				fmt.Sprintf("TokenError: invalid token found at line %v at offset %v", tokenStartingLine, tokenStartingLineOffset),
				nil,
			)
		}
	}
}
