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
*/

type Token struct {
	TypeOfToken dfa.TokenType
	Lexemme     []rune
	Line        uint32
	Offset      uint32
}

func (t *Token) SetTokenProperties(_type dfa.TokenType, _lineNum uint32, _offset uint32, _lexemme []rune) {
	t.TypeOfToken = _type
	t.Line = _lineNum
	t.Offset = _offset
	t.Lexemme = _lexemme
}

func (t *Token) ToString() string {
	return fmt.Sprintf("|%d|%d| [%s]Token -> `%s`", t.Line, t.Offset, string(t.TypeOfToken), string(t.Lexemme))
}

// TODO: Setup a new struct here which will have the lastInput properties along with the dfa tokens

type LexicalAnalyzer struct {
	source              *bufio.Reader
	tokenDFAa           map[dfa.TokenType]dfa.DFA
	lineNum             uint32
	lineOffset          uint32
	prevLineOffset      uint32 // When a '\n' rune ends a token then it also ends up resetting offset to 0 (and lineNum++) which makes this the only way one can the offset for the token which got ended by the `\n`
	lexemme             []rune
	lastInputExists     bool
	lastInput           rune
	lastInputLineNum    uint32
	lastInputLineOffset uint32
}

// *NOTE: We need this lastInput fields because when the dfas detect a valid token, the scanner needs to consume one more rune, which will make all the dfas return invalid. And then the token is returned. But, this process consumes a rune whiich didn't belong to the previous token, and so we save it and reuse it in the next time the token reader function is called. Yes, this does mean that the procesing of a starting rune of each token is done TWICE. The fix is easy by adding the processing info in the struct as well... But, if I add it right now that can mess up the code a lot.
// ^ This might also potentially need changing some things in the dfa package, because I have not exposed a way to start each DFA on some particular state or even a way to save all the dfas state in an easy to do manner.
// So, the twice execution problem is here to stay for now, But I'll still put a todo. ->
// TODO: Fix the problem inefficiency of the 1st rune of a tokken going through all the lexemmes DFAs twice.

func (scanner *LexicalAnalyzer) Initialize(source *bufio.Reader) {
	scanner.source = source
	scanner.tokenDFAa = dfa.GenerateDFAs()
	scanner.lineNum = 0
	scanner.lineOffset = 0
	scanner.lexemme = nil
	scanner.lastInputExists = false
}

func (scanner *LexicalAnalyzer) resetDFAs() {
	for _, token := range dfa.TokensList {
		scanner.tokenDFAa[token].Reset()
	}
}

// This function reades one rune at a time from the source reader and returns 1 token at a time in return. It follows `maximal munching` methodlogy for settling tie between 2 valid token dfas being satisfied. And if both DFAs end up having the same token length then the token which is written later in the `dfa.TokensList` is given higher priority and is returned
func (scanner *LexicalAnalyzer) ReadToken() (Token, error) {
	scanner.lexemme = nil // Resetting the lexemme to be stored

	returnToken := Token{}
	var input rune
	var err error
	var tokenLineNum uint32
	var tokenLineOffset uint32
	if scanner.lastInputExists {
		tokenLineNum = scanner.lastInputLineNum
		tokenLineOffset = scanner.lastInputLineOffset
	} else {
		tokenLineNum = scanner.lineNum
		tokenLineOffset = scanner.lineOffset
	}

	isAnyValid := false // For checking if any vlaid token was found in this run of this function's for loop.
	isAnyIntermediate := false
	var foundValidToken dfa.TokenType
	var foundIntermediateToken dfa.TokenType

	lastLoop_isAnyValid := false
	lastLoop_isAnyIntermediate := false
	var lastLoop_foundValidToken dfa.TokenType
	var lastLoop_foundIntermediateToken dfa.TokenType

	for i := 0; ; i++ {
		// Read one rune from the source
		if scanner.lastInputExists {
			input = scanner.lastInput
			scanner.lastInputExists = false
		} else {
			input, _, err = scanner.source.ReadRune()
			if err != nil && err != io.EOF {
				scanner.resetDFAs()
				return returnToken, errorhandler.RetErr("", err)
			}
			if err == io.EOF {
				scanner.resetDFAs()
				// This is the end of the file so just wrap up your findings and return
				if i == 0 {
					// Users should run this function to get the EOF token even if internally in this function the EOF was reached in the last run
					returnToken.SetTokenProperties(dfa.EOF, tokenLineNum, tokenLineOffset, []rune(string(dfa.EOF)))
					// TODO: See if this lineNum and lineOffset is okay for not returning the eof
					return returnToken, io.EOF
				}

				// Check if any valid token was found in the last iteration, if yes then we report it as our found token
				if lastLoop_isAnyValid {
					returnToken.SetTokenProperties(
						lastLoop_foundValidToken,
						tokenLineNum,
						tokenLineOffset,
						scanner.lexemme,
					)
					return returnToken, nil
				}
				// If any intermediates were there then we will use those for error reporting
				if lastLoop_isAnyIntermediate {
					return returnToken, errorhandler.RetErr(
						fmt.Sprintf(
							"TokenError: invalid token found at line %v at offset %v, most resembling token type was %v",
							tokenLineNum,
							tokenLineOffset,
							string(lastLoop_foundIntermediateToken),
						),
						nil,
					)
				}
				// I do not expect the code to reach this line. Because, to reach here the last iteration of this function would have to have all Invalid tokens, and still decided to continue parsing. Which should'nt happen.
				// Return error and not EOF token here, users can get the EOF token in the next run despite the error
				return returnToken, errorhandler.RetErr(
					fmt.Sprintf("TokenError: invalid token found at line %v at offset %v", scanner.lineNum, scanner.lineOffset-uint32(len(scanner.lexemme))),
					nil,
				)
			}
			// Record the offsets and the lineNums in the scanner
			scanner.lineOffset++
			if input == '\n' {
				scanner.lineNum++
				scanner.prevLineOffset = scanner.lineOffset - 1
				scanner.lineOffset = 0
			}
		}
		scanner.lexemme = append(scanner.lexemme, input)

		// TODO: This loop can be done concurrently, and that should speed up the processing speed by a lot.
		// Execute a step in all the dfas with the current rune.
		for i := 0; i < len(dfa.TokensList); i++ { // The token written after in the order of dfa.TokensList will get higher priority
			token := dfa.TokensList[i]
			dfaForToken := scanner.tokenDFAa[token]
			dfa_result := dfaForToken.Step(input)
			if dfa_result.IsValid() {
				isAnyValid = true
				foundValidToken = token
			}
			if dfa_result.IsIntermediate() {
				isAnyIntermediate = true
				foundIntermediateToken = token
			}
		}

		// Stop iterating if all the DFAs are yielding INVALID
		if !isAnyValid && !isAnyIntermediate {
			scanner.lastInputExists = true
			scanner.lastInput = input
			if scanner.lineOffset == 0 {
				scanner.lastInputLineNum = scanner.lineNum - 1
				scanner.lastInputLineOffset = scanner.prevLineOffset
			} else {
				scanner.lastInputLineNum = scanner.lineNum
				scanner.lastInputLineOffset = scanner.lineOffset - 1
			}
			scanner.resetDFAs()
			// Check if any valid token was found in the last iteration, if yes then we report it as our found token
			if lastLoop_isAnyValid {
				returnToken.SetTokenProperties(
					lastLoop_foundValidToken,
					tokenLineNum,
					tokenLineOffset,
					scanner.lexemme[:len(scanner.lexemme)-1], // We remove the last rune which was not valid
				)
				return returnToken, nil
			}
			// If any intermediates were there then we will use those for error reporting
			if lastLoop_isAnyIntermediate {
				return returnToken, errorhandler.RetErr(
					fmt.Sprintf(
						"TokenError: invalid token found at line %v at offset %v, most resembling token type was %v",
						tokenLineNum,
						tokenLineOffset,
						string(lastLoop_foundIntermediateToken),
					),
					nil,
				)
			}
			// The Program Counter can only get here if this is the 1st iteration and the very 1st rune did not satisfay any of the token types' dfa
			scanner.lastInputExists = false
			return returnToken, errorhandler.RetErr(
				fmt.Sprintf("TokenError: invalid token found at line %v at offset %v", tokenLineNum, tokenLineOffset),
				nil,
			)
		}

		// Set this loop important values in the lastLoop variables
		lastLoop_isAnyValid = isAnyValid
		lastLoop_isAnyIntermediate = isAnyIntermediate
		lastLoop_foundValidToken = foundValidToken
		lastLoop_foundIntermediateToken = foundIntermediateToken

		// Reset variables (which need it) outside of the loop
		isAnyValid = false
		isAnyIntermediate = false
	}
}
