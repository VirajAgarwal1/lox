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
	lexemme     []rune
	Line        uint32
	Offset      uint32
}

func (t *Token) SetTokenProperties(_type dfa.TokenType, _lineNum uint32, _offset uint32, _lexemme []rune) {
	t.TypeOfToken = _type
	t.Line = _lineNum
	t.Offset = _offset
	copy(t.lexemme, _lexemme)
}

type LexicalAnalyzer struct {
	source          *bufio.Reader
	tokenDFAa       map[dfa.TokenType]dfa.DFA
	lineNum         uint32
	lineOffset      uint32
	lexemme         []rune
	lastInputExists bool
	lastInput       rune
}

func (scanner *LexicalAnalyzer) Initialize(source *bufio.Reader) {
	scanner.source = source
	scanner.tokenDFAa = dfa.GenerateDFAs()
	scanner.lineNum = 0
	scanner.lineOffset = 0
	scanner.lexemme = nil
	scanner.lastInputExists = false
}

// This function reades one rune at a time from the source reader and returns 1 token at a time in return. It follows `maximal munching` methodlogy for settling tie between 2 valid token dfas being satisfied. And if both DFAs end up having the same token length then the token which is written later in the `dfa.TokensList` is given higher priority and is returned
func (scanner *LexicalAnalyzer) ReadToken() (Token, error) {
	scanner.lexemme = nil // Resetting the lexemme to be stored

	dfaTokenResults := make(map[dfa.TokenType]dfa.DfaReturn, len(dfa.TokensList))
	returnToken := Token{}
	var input rune
	var err error

	lastLoop_isAnyValid := false
	lastLoop_isAnyIntermediate := false
	var lastLoop_foundValidToken dfa.TokenType
	var lastLoop_foundIntermediateToken dfa.TokenType

	isAnyValid := false
	isAnyIntermediate := false
	var foundValidToken dfa.TokenType
	var foundIntermediateToken dfa.TokenType
	for i := 0; ; i++ {
		// Read one rune from the source
		if !scanner.lastInputExists {
			input, _, err = scanner.source.ReadRune()
			if err != nil && err != io.EOF {
				return returnToken, errorhandler.RetErr("", err)
			}
			if err == io.EOF {
				// This is the end of the file so just wrap up your findings and return
				if i == 0 {
					// Users should run this function to get the EOF token even if internally in this function the EOF was reached in the last run
					returnToken.SetTokenProperties(dfa.EOF, scanner.lineNum, scanner.lineOffset, []rune(string(dfa.EOF)))
					return returnToken, io.EOF
				}

				// Check if any valid token was found in the last iteration, if yes then we report it as our found token
				if lastLoop_isAnyValid {
					returnToken.SetTokenProperties(
						lastLoop_foundValidToken,
						scanner.lineNum,
						scanner.lineOffset-uint32(len(scanner.lexemme)),
						scanner.lexemme,
					)
					return returnToken, nil
				}
				// If any intermediates were there then we will use those for error reporting
				if lastLoop_isAnyIntermediate {
					return returnToken, errorhandler.RetErr(
						fmt.Sprintf(
							"TokenError: invalid token found at line %v at offset %v, most resembling token type was %v",
							scanner.lineNum,
							scanner.lineOffset-uint32(len(scanner.lexemme)),
							string(lastLoop_foundIntermediateToken),
						),
						nil,
					)
				}
				// Return error and not EOF token here, users can get the EOF token in the next run despite the error
				return returnToken, errorhandler.RetErr(
					fmt.Sprintf("TokenError: invalid token found at line %v at offset %v", scanner.lineNum, scanner.lineOffset-uint32(len(scanner.lexemme))),
					nil,
				)
			}
		} else {
			input = scanner.lastInput
			scanner.lastInputExists = false
		}
		scanner.lexemme = append(scanner.lexemme, input)

		// Execute a step in all the dfas with the current rune.
		for i := 0; i < len(dfa.TokensList); i++ { // The token written after in order will get higher priority
			token := dfa.TokensList[i]
			dfaForToken := scanner.tokenDFAa[token]
			dfa_result := dfaForToken.Step(input)
			dfaTokenResults[token] = dfa_result
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
			// Check if any valid token was found in the last iteration, if yes then we report it as our found token
			if lastLoop_isAnyValid {
				returnToken.SetTokenProperties(
					lastLoop_foundValidToken,
					scanner.lineNum,
					scanner.lineOffset-uint32(len(scanner.lexemme)),
					scanner.lexemme,
				)
				return returnToken, nil
			}
			// If any intermediates were there then we will use those for error reporting
			if lastLoop_isAnyIntermediate {
				return returnToken, errorhandler.RetErr(
					fmt.Sprintf(
						"TokenError: invalid token found at line %v at offset %v, most resembling token type was %v",
						scanner.lineNum,
						scanner.lineOffset-uint32(len(scanner.lexemme)),
						string(lastLoop_foundIntermediateToken),
					),
					nil,
				)
			}
			// The Program Counter can only get here if this is the 1st iteration and the very 1st rune did not satisfay any of the token types' dfa
			return returnToken, errorhandler.RetErr(
				fmt.Sprintf("TokenError: invalid token found at line %v at offset %v", scanner.lineNum, scanner.lineOffset-uint32(len(scanner.lexemme))),
				nil,
			)
		}

		// Record the offsets and the lineNums in the scanner
		scanner.lineOffset++
		if input == '\n' {
			scanner.lineNum++
			scanner.lineOffset = 0
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
