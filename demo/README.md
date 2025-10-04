# Demo Programs

This directory contains demonstration programs showing how to use the various components of the Lox interpreter. Each demo is a standalone Go program that you can run to see the functionality in action.

## Directory Structure

```
demo/
├── lexer_demos/          # Lexical analyzer demonstrations
├── parser_demos/         # Recursive descent parser demonstrations
└── streamable_parser_demos/  # LL(1) streaming parser demonstrations
```

## Running Demos

### Option 1: Modify main.go

Edit `/main.go` to call the demo function you want:

```go
func main() {
    // Uncomment the demo you want to run:
    
    // Lexer demos
    // lexer_demos.Sample_dfa_demo()
    // lexer_demos.Sample_scanner_demo()
    
    // Parser demos
    // parser_demos.Sample_parser_demo()
    // parser_demos.Sample_generate_parser_demo()
    
    // Streamable parser demos
    // streamable_parser_demos.Sample_streamable_parser_demo()
    // streamable_parser_demos.Sample_compute_firsts()
    // streamable_parser_demos.Sample_compute_follow()
    // streamable_parser_demos.Sample_ebnf_to_bnf()
    streamable_parser_demos.Sample_parser_writer()
}
```

Then run:
```bash
go run main.go
```

### Option 2: Import and Run Directly

Create your own test file:

```go
package main

import "github.com/VirajAgarwal1/lox/demo/lexer_demos"

func main() {
    lexer_demos.Sample_scanner_demo()
}
```

## Lexer Demos

Location: `demo/lexer_demos/`

### dfa.demo.go

**Purpose:** Demonstrates how DFAs (Deterministic Finite Automata) work for tokenization.

**What it shows:**
- Creating DFAs for different token types
- Feeding characters into DFAs one at a time
- How DFAs return VALID, INTERMEDIATE, or INVALID states
- Conflict resolution when multiple DFAs match

**Example usage:**
```go
func Sample_dfa_demo() {
    dfas := dfa.GenerateDFAs()
    
    // Test string "while"
    input := "while"
    for _, char := range input {
        for tokenType, automaton := range dfas {
            result := automaton.Step(char)
            if result.IsValid() {
                fmt.Printf("Matched: %s\n", tokenType)
            }
        }
    }
}
```

**Video:** See `dfa-demo.mp4` in the root directory for a visual walkthrough.

### scanner.demo.go

**Purpose:** Demonstrates the complete lexical analyzer scanning source code into tokens.

**What it shows:**
- Initializing the lexical analyzer
- Reading tokens one at a time from source code
- Token properties: type, lexeme, line, offset
- Handling different token types (keywords, operators, literals, identifiers)

**Example usage:**
```go
func Sample_scanner_demo() {
    source := "var x = 10 + 20;"
    reader := bufio.NewReader(strings.NewReader(source))
    
    scanner := &lexer.LexicalAnalyzer{}
    scanner.Initialize(reader)
    
    for {
        token, err := scanner.ReadToken()
        if err == io.EOF {
            break
        }
        fmt.Println(token.ToString())
    }
}
```

**Video:** See `scanner-demo.mp4` in the root directory for a visual demonstration.

**Expected output:**
```
|0|0| [var]Token -> `var`
|0|4| [IDENTIFIER]Token -> `x`
|0|6| [=]Token -> `=`
|0|8| [NUMBER]Token -> `10`
|0|11| [+]Token -> `+`
|0|13| [NUMBER]Token -> `20`
|0|15| [;]Token -> `;`
```

## Parser Demos

Location: `demo/parser_demos/`

### parser.demo.go

**Purpose:** Demonstrates using the generated recursive descent parser to parse and evaluate expressions.

**What it shows:**
- Creating a buffered lexer from source code
- Parsing expressions with operator precedence
- Evaluating the resulting AST
- Handling complex expressions with parentheses

**Example usage:**
```go
func Sample_parser_demo() {
    source := "1 + 2 * 3"
    
    // Create lexer
    lex := &lexer.LexicalAnalyzer{}
    lex.Initialize(bufio.NewReader(strings.NewReader(source)))
    
    // Create buffered lexer
    bufferedLex := &lexer.BufferedLexer{}
    bufferedLex.Initialize(lex)
    
    // Parse
    nodes, ok, err := parser.Parse_expression(bufferedLex)
    if err != nil || !ok {
        panic("Parse failed")
    }
    
    // Evaluate
    result := nodes[0].Evaluate()
    fmt.Printf("Result: %v\n", result.Inner)
}
```

**Expected output:**
```
Result: 7
```

### generate_parser.demo.go

**Purpose:** Demonstrates generating a parser from a grammar specification file.

**What it shows:**
- Reading an EBNF grammar file
- Parsing the grammar definition
- Generating Go code for the parser
- Writing the generated code to a file

**Example usage:**
```go
func Sample_generate_parser_demo() {
    // Read grammar file
    grammarFile, _ := os.Open("parser/lox.grammar")
    defer grammarFile.Close()
    
    // Parse grammar
    scanner := &lexer.LexicalAnalyzer{}
    scanner.Initialize(bufio.NewReader(grammarFile))
    
    grammarRules, _ := grammar.ProcessGrammarDefinition(scanner)
    
    // Generate parser
    outputFile, _ := os.Create("parser/generated_parser.go")
    defer outputFile.Close()
    
    grammar.WriteParserForGrammar(grammarRules, outputFile, "expression")
}
```

This creates a new `generated_parser.go` file with all the parsing functions.

## Streamable Parser Demos

Location: `demo/streamable_parser_demos/`

### ⚠️ IMPORTANT: Two-Step Process

**Parser generation and usage must happen in TWO separate runs!**

**Why?** Go compiles code at build-time, not run-time. When you generate parser code, it writes Go source files to disk, but those files won't be compiled into your binary until the next build. You cannot generate code and immediately use it in the same execution.

**The Process:**
1. **First run:** Generate the parser code → writes `generated_parser.go`
2. **Second run:** Use the generated parser → reads the compiled version from the previous run

### parser.demo.go

**Purpose:** Demonstrates LL(1) streaming parser generation and usage.

#### `Sample_generate_streamable_parser()` - Step 1: Generate

Generates a StreamableParser from a grammar file. Run this FIRST.

**Example usage in main.go:**
```go
func main() {
    // Step 1: Generate the parser
    streamable_parser_demos.Sample_generate_streamable_parser()
}
```

Run:
```bash
go run main.go
```

**Output:**
```
=== Step 1: Generate Parser Code ===

✓ Parser generated successfully!
  Output: /path/to/streamable_parser/generated_parser.go

NEXT: Change main.go to call Sample_streamable_parser_demo() and run again.
```

#### `Sample_streamable_parser_demo()` - Step 2: Use

Uses the generated parser to parse expressions. Run this SECOND (after generation).

**Example usage in main.go:**
```go
func main() {
    // Step 2: Use the generated parser
    streamable_parser_demos.Sample_streamable_parser_demo()
}
```

Run:
```bash
go run main.go
```

**Output:**
```
=== Step 2: Use the Generated Parser ===

Parsing: 1+2*3-4/"hello"+true

Parse tree:
└── ⟨expression⟩
    ├── ⟨comma⟩
    │   ├── ⟨equality⟩
    │   │   ├── ⟨comparison⟩
    │   │   │   ├── ⟨term⟩
    │   │   │   │   ├── ⟨factor⟩
    │   │   │   │   │   ├── ⟨unary⟩
    │   │   │   │   │   │   └── Leaf(1)
    ...
```

### sample_compute_firsts.go

**Purpose:** Demonstrates computing FIRST sets for grammar rules.

**What it shows:**
- Reading a grammar file
- Converting EBNF to BNF
- Computing FIRST sets for each production
- Displaying the FIRST sets

**Why it's useful:**
FIRST sets are essential for LL(1) parsing - they tell the parser which production to choose based on the lookahead token.

**Example output:**
```
FIRST sets:
expression: {NUMBER, IDENTIFIER, (}
term: {NUMBER, IDENTIFIER, (}
factor: {NUMBER, IDENTIFIER, (}
```

### sample_compute_follow.go

**Purpose:** Demonstrates computing FOLLOW sets for non-terminals.

**What it shows:**
- Computing FOLLOW sets from grammar rules
- How FOLLOW sets depend on production rules
- Handling epsilon productions

**Why it's useful:**
FOLLOW sets help the parser handle epsilon productions and detect syntax errors.

**Example output:**
```
FOLLOW sets:
expression: {$, )}
term: {+, -, $, )}
factor: {*, /, +, -, $, )}
```

### sample_ebnf_to_bnf.go

**Purpose:** Demonstrates converting EBNF grammars to BNF form.

**What it shows:**
- Reading EBNF grammar with `*`, `+`, and grouping
- Converting to pure BNF form
- Introducing artificial non-terminals
- Preserving grammar semantics

**Example:**
```
EBNF:  expr -> term ("+" term)*

BNF:   expr -> term __repeat_expr_0
       __repeat_expr_0 -> "+" term __repeat_expr_0 | ε
```


## Video Demonstrations

The project includes video demonstrations:

- **`dfa-demo.mp4`** - Visual explanation of how DFAs work for tokenization
- **`scanner-demo.mp4`** - Complete lexical analysis demonstration

These videos provide visual walkthroughs

## Common Patterns

### Pattern 1: Lexer Initialization

```go
source := "your source code here"
reader := bufio.NewReader(strings.NewReader(source))

scanner := &lexer.LexicalAnalyzer{}
scanner.Initialize(reader)
```

### Pattern 2: Token Reading Loop

```go
for {
    token, err := scanner.ReadToken()
    if err == io.EOF {
        break
    }
    if err != nil {
        panic(err)
    }
    // Process token
}
```

### Pattern 3: Parser Usage

```go
// Create buffered lexer
bufferedLex := &lexer.BufferedLexer{}
bufferedLex.Initialize(scanner)

// Parse
nodes, ok, err := parser.Parse_expression(bufferedLex)
if err != nil || !ok {
    panic("Parse failed")
}
```

### Pattern 4: Streamable Parser Events

```go
parser := &streamable_parser.StreamableParser{}
parser.Initialize(bufferedLexer)

for {
    event := parser.Parse()
    if event == nil {
        continue // Skip artificial non-terminals
    }
    
    switch event.Type {
    case EmitElemType_Start:
        // Begin non-terminal
    case EmitElemType_End:
        // Finish non-terminal
    case EmitElemType_Leaf:
        // Terminal matched
    case EmitElemType_Error:
        // Syntax error
        return
    }
}
```

## Tips for Experimenting

1. **Start with lexer demos** - Understanding tokenization is fundamental
2. **Try different inputs** - Modify source strings to see different behaviors
3. **Add print statements** - See intermediate states during parsing
4. **Break things** - Try invalid syntax to see error handling
5. **Read the generated code** - Look at `generated_parser.go` to understand what's generated

## See Also

- [Lexer Documentation](../lexer/README.md)
- [Parser Documentation](../parser/README.MD)
- [Streamable Parser Documentation](../streamable_parser/README.md)
- [Tests Directory](../tests/README.md) - Unit tests for all components


