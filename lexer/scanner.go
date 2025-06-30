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
	Line        uint32
	Offset      uint32
}

func (t *Token) SetTokenProperties(_type dfa.TokenType, _lineNum uint32, _offset uint32) {
	t.TypeOfToken = _type
	t.Line = _lineNum
	t.Offset = _offset
}

type LexicalAnalyzer struct {
	source     *bufio.Reader
	tokenDFAa  map[dfa.TokenType]dfa.DFA
	lineNum    uint32
	lineOffset uint32
	lexemme    []rune
}

func (scanner *LexicalAnalyzer) Initialize(source *bufio.Reader) {
	scanner.tokenDFAa = dfa.GenerateDFAs()
	scanner.source = source
}

func (scanner *LexicalAnalyzer) ReadToken() (Token, error) {
	// TODO: Add some way of storing the input rune (from which eveyrthing went INVALID) in state and then always start checking from that rune
	dfaTokenResults := make(map[dfa.TokenType]dfa.DfaReturn, len(dfa.TokensList))
	returnToken := Token{}

	lastLoop_isAnyValid := false
	lastLoop_isAnyIntermediate := false
	var lastLoop_foundValidToken dfa.TokenType
	var lastLoop_foundIntermediateToken dfa.TokenType

	isAnyValid := false
	isAnyIntermediate := false
	var foundValidToken dfa.TokenType
	var foundIntermediateToken dfa.TokenType
	for {
		// Read one rune from the source
		input, _, err := scanner.source.ReadRune()
		if err != nil && err != io.EOF {
			return returnToken, errorhandler.RetErr("", err)
		}
		if err == io.EOF {
			returnToken.SetTokenProperties(dfa.EOF, scanner.lineNum, scanner.lineOffset)
			return returnToken, io.EOF
		}

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

		if !isAnyValid && !isAnyIntermediate {
			// Check if any valid token was found in the last iteration
			if lastLoop_isAnyValid {
				returnToken.SetTokenProperties(
					lastLoop_foundValidToken,
					scanner.lineNum,
					scanner.lineOffset-uint32(len(string(lastLoop_foundValidToken))))
				return returnToken, nil
			}

			// If any intermediates were there then we will use those for error reporting
			if lastLoop_isAnyIntermediate {
				return returnToken, errorhandler.RetErr(
					fmt.Sprintf(
						"TokenError: invalid token found at line %v at offset %v, most resembling token type was %v",
						scanner.lineNum,
						scanner.lineOffset,
						string(lastLoop_foundIntermediateToken),
					),
					nil,
				)
			}

			// The Program Counter can only get here if this is the 1st iteration and the very 1st rune did not satisfay any of the token types' dfa
			return returnToken, errorhandler.RetErr(
				fmt.Sprintf("TokenError: invalid token found at line %v at offset %v", scanner.lineNum, scanner.lineOffset),
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
