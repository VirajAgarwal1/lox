package dfa

const (
	INVALID DfaResult = iota
	INTERMEDIATE
	VALID

	// These are all lexemes which will identified by the scanner.
	EOF TokenType = "EOF"
	// Literals
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"
	COMMENT    TokenType = "COMMENT"
	// Single-char tokens
	WHITESPACE  TokenType = " "
	NEWLINE     TokenType = "\n"
	LEFT_PAREN  TokenType = "("
	RIGHT_PAREN TokenType = ")"
	LEFT_BRACE  TokenType = "{"
	RIGHT_BRACE TokenType = "}"
	COMMA       TokenType = ","
	DOT         TokenType = "."
	MINUS       TokenType = "-"
	PLUS        TokenType = "+"
	SEMICOLON   TokenType = ";"
	SLASH       TokenType = "/"
	STAR        TokenType = "*"
	// One-or-two char tokens
	BANG          TokenType = "!"
	BANG_EQUAL    TokenType = "!="
	EQUAL         TokenType = "="
	EQUAL_EQUAL   TokenType = "=="
	GREATER       TokenType = ">"
	GREATER_EQUAL TokenType = ">="
	LESS          TokenType = "<"
	LESS_EQUAL    TokenType = "<="
	// Keywords.
	AND    TokenType = "and"
	CLASS  TokenType = "class"
	ELSE   TokenType = "else"
	FALSE  TokenType = "false"
	FUN    TokenType = "fun"
	FOR    TokenType = "for"
	IF     TokenType = "if"
	NIL    TokenType = "nil"
	OR     TokenType = "or"
	PRINT  TokenType = "print"
	RETURN TokenType = "return"
	SUPER  TokenType = "super"
	THIS   TokenType = "this"
	TRUE   TokenType = "true"
	VAR    TokenType = "var"
	WHILE  TokenType = "while"
)

// In lexer/scanner.go, we use this list to iterate over the tokens in order of priority. Where the tokens retiurn later get higher priority.
var TokensList = []TokenType{
	EOF,

	IDENTIFIER,
	STRING,
	NUMBER,
	COMMENT,

	WHITESPACE,
	NEWLINE,
	LEFT_PAREN,
	RIGHT_PAREN,
	LEFT_BRACE,
	RIGHT_BRACE,
	COMMA,
	DOT,
	MINUS,
	PLUS,
	SEMICOLON,
	SLASH,
	STAR,

	BANG,
	BANG_EQUAL,
	EQUAL,
	EQUAL_EQUAL,
	GREATER,
	GREATER_EQUAL,
	LESS,
	LESS_EQUAL,

	AND,
	CLASS,
	ELSE,
	FALSE,
	FUN,
	FOR,
	IF,
	NIL,
	OR,
	PRINT,
	RETURN,
	SUPER,
	THIS,
	TRUE,
	VAR,
	WHILE,
}

type summaryOfAllDfaStates struct {
	IsAnyValidToken        bool
	IsAnyIntermediateToken bool
	ValidToken             TokenType
	IntermediateToken      TokenType
}
type DFAStatesManager struct {
	TokenDFAList           []DFA
	DfaResultForToken      []DfaResult
	CurrentLoopDfaResults  *summaryOfAllDfaStates
	PreviousLoopDfaResults *summaryOfAllDfaStates
}
type DFA interface {
	Step(input rune) DfaResult
	Reset()
}
type TokenType string
type DfaResult byte

func (resultsSummary *summaryOfAllDfaStates) initialize() {
	resultsSummary.IsAnyValidToken = false
	resultsSummary.IsAnyIntermediateToken = false
}
func (resultsSummary *summaryOfAllDfaStates) AreAllInvalid() bool {
	return !(resultsSummary.IsAnyIntermediateToken || resultsSummary.IsAnyValidToken)
}

func (stateManager *DFAStatesManager) Initialize() {
	stateManager.TokenDFAList = stateManager.GenerateDFAs()
	stateManager.DfaResultForToken = make([]DfaResult, len(TokensList))
	for i := range len(TokensList) {
		stateManager.DfaResultForToken[i] = VALID
	}
	stateManager.CurrentLoopDfaResults = &summaryOfAllDfaStates{}
	stateManager.PreviousLoopDfaResults = &summaryOfAllDfaStates{}
	stateManager.CurrentLoopDfaResults.initialize()
	stateManager.PreviousLoopDfaResults.initialize()
}
func (stateManager *DFAStatesManager) putCurrentSummaryInPrevious() {
	stateManager.PreviousLoopDfaResults.IsAnyValidToken = stateManager.CurrentLoopDfaResults.IsAnyValidToken
	stateManager.PreviousLoopDfaResults.IsAnyIntermediateToken = stateManager.CurrentLoopDfaResults.IsAnyIntermediateToken
	stateManager.PreviousLoopDfaResults.ValidToken = stateManager.CurrentLoopDfaResults.ValidToken
	stateManager.PreviousLoopDfaResults.IntermediateToken = stateManager.CurrentLoopDfaResults.IntermediateToken
}
func (stateManager *DFAStatesManager) Step(input rune) {
	// Execute a step in all the dfas with the current rune.
	for i := range len(TokensList) {
		// The token written after in the order of dfa.TokensList will get higher priority
		if stateManager.DfaResultForToken[i] == INVALID {
			continue
		}
		token := TokensList[i]
		stateManager.DfaResultForToken[i] = stateManager.TokenDFAList[i].Step(input)
		if stateManager.DfaResultForToken[i].IsValid() {
			stateManager.CurrentLoopDfaResults.IsAnyValidToken = true
			stateManager.CurrentLoopDfaResults.ValidToken = token
		}
		if stateManager.DfaResultForToken[i].IsIntermediate() {
			stateManager.CurrentLoopDfaResults.IsAnyIntermediateToken = true
			stateManager.CurrentLoopDfaResults.IntermediateToken = token
		}
	}
}
func (stateManager *DFAStatesManager) ClearCurrentLoopDfaResults() {
	stateManager.putCurrentSummaryInPrevious()
	stateManager.CurrentLoopDfaResults.initialize()
}
func (stateManager *DFAStatesManager) ResetAllDFAs() {
	for i := range len(TokensList) {
		stateManager.TokenDFAList[i].Reset()
		stateManager.DfaResultForToken[i] = VALID
	}
}
func (stateManager *DFAStatesManager) FullReset() {
	stateManager.ResetAllDFAs()
	stateManager.CurrentLoopDfaResults.initialize()
	stateManager.PreviousLoopDfaResults.initialize()
}
func (stateManager *DFAStatesManager) GenerateDFAs() []DFA {

	output := make([]DFA, len(TokensList))
	for i, token := range TokensList {
		if token == EOF {
			dfa := &EofDFA{}
			dfa.Initialize()
			output[i] = dfa
			continue // EOF is not a string, so we skip it
		}
		if token == IDENTIFIER {
			dfa := &IdentifierDFA{}
			dfa.Initialize()
			output[i] = dfa
			continue
		}
		if token == STRING {
			dfa := &StringDFA{}
			dfa.Initialize()
			output[i] = dfa
			continue
		}
		if token == NUMBER {
			dfa := &NumberDFA{}
			dfa.Initialize()
			output[i] = dfa
			continue
		}
		if token == COMMENT {
			dfa := &CommentDFA{}
			dfa.Initialize()
			output[i] = dfa
			continue
		}
		if token == WHITESPACE {
			dfa := &WhitespaceDFA{}
			dfa.Initialize()
			output[i] = dfa
			continue
		}
		if token == NEWLINE {
			dfa := &NewlineDFA{}
			dfa.Initialize()
			output[i] = dfa
			continue
		}
		dfa := &InputStringDFA{}
		dfa.Initialize(string(token))
		output[i] = dfa
	}

	return output
}

func (s *DfaResult) ToString() string {
	if int(*s) == 0 {
		return "INVALID"
	}
	if int(*s) == 1 {
		return "INTERMEDIATE"
	}
	if int(*s) == 2 {
		return "VALID"
	}
	return ""
}
func (s *DfaResult) IsValid() bool {
	return *s == VALID
}
func (s *DfaResult) IsIntermediate() bool {
	return *s == INTERMEDIATE
}
func (s *DfaResult) IsInvalid() bool {
	return *s == INVALID
}
