# Streamable Parser

A table-driven LL(1) predictive parser with streaming token consumption and automatic parser generation from EBNF grammars.

## What This Does

This parser implements a classic LL(1) parsing algorithm that uses precomputed FIRST and FOLLOW sets to make predictive parsing decisions. Unlike the recursive descent parser, this one uses an explicit stack and parsing table, making it suitable for streaming applications where you want fine-grained control over parsing progress.

**Key Features:**
- LL(1) table-driven parsing
- Automatic FIRST and FOLLOW set computation
- EBNF to BNF conversion
- Event-based parsing with start/end/leaf emissions
- Streaming token consumption
- Error recovery with token synchronization

## Architecture

### Parser Generation Pipeline

The parser generation involves several steps:

```
EBNF Grammar
    ↓
[EBNF to BNF Converter]
    ↓
BNF Grammar
    ↓
[FIRST Set Computer]
    ↓
[FOLLOW Set Computer]
    ↓
[Parser Code Generator]
    ↓
Generated Parser Code
```

### Components

1. **EBNF to BNF Converter** (`parser_generator/ebnf_to_bnf/`)
   - Converts extended grammar notation to basic form
   - Eliminates `*`, `+`, and optional operators
   - Introduces artificial non-terminals for repetitions

2. **FIRST/FOLLOW Set Computer** (`parser_generator/first_follow/`)
   - Computes FIRST sets for production sequences
   - Computes FOLLOW sets for non-terminals
   - Handles epsilon productions correctly

3. **Grammar File Parser** (`parser_generator/grammar_file_parser/`)
   - Parses grammar specification files
   - Validates grammar syntax
   - Builds internal grammar representation

4. **Parser Code Generator** (`parser_generator/parser_writer/`)
   - Generates Go code for the parser
   - Embeds FIRST and FOLLOW sets
   - Creates grammar rules data structure
   - Generates type-safe parsing functions

## How LL(1) Parsing Works

LL(1) stands for:
- **L**eft-to-right scanning
- **L**eftmost derivation
- **1** token lookahead

The parser uses a stack-based approach:

1. **Initialize**: Push start symbol onto stack
2. **Loop**:
   - Peek at top of stack and lookahead token
   - If stack top is terminal: match with lookahead
   - If stack top is non-terminal: use FIRST/FOLLOW to choose production
   - Replace stack top with production right-hand side
3. **Accept**: When stack is empty and input is exhausted

### Predictive Parsing

The key to LL(1) parsing is making the right prediction:

```
Given: Stack top = non-terminal N, Lookahead = token T

1. For each production N → Alpha:
   - If T is in FIRST(Aplha), use this production
   
2. If no production matches:
   - If ε is in FIRST(N) and T is in FOLLOW(N), use ε-production
   - Otherwise, syntax error
```

## Usage

### Generating a Parser

```go
package main

import (
    "github.com/VirajAgarwal1/lox/streamable_parser/parser_generator/parser_writer"
)

func main() {
    // Generate parser from grammar file
    err := parser_writer.WriteParserForGrammar(
        "input_grammar.txt",
        "generated_parser.go",
        "expression", // start symbol
        "myparser",   // package name
    )
    if err != nil {
        panic(err)
    }
}
```

### Using the Parser

```go
package main

import (
    "bufio"
    "io"
    "strings"
    "github.com/VirajAgarwal1/lox/streamable_parser"
    "github.com/VirajAgarwal1/lox/lexer"
)

func main() {
    // Create lexer
    source := "1 + 2 * 3"
    lex := &lexer.LexicalAnalyzer{}
    lex.Initialize(bufio.NewReader(strings.NewReader(source)))
    
    // Create buffered lexer
    bufferedLex := &lexer.BufferedLexicalAnalyzer{}
    bufferedLex.Initialize(lex)
    
    // Create parser
    parser := &streamable_parser.StreamableParser{}
    parser.Initialize(bufferedLex)
    
    // Parse in streaming fashion
    for {
        event := parser.Parse()
        if event == nil {
            continue // artificial non-terminals are filtered out
        }
        
        switch event.Type {
        case streamable_parser.EmitElemType_Start:
            println("Start:", event.Content)
            
        case streamable_parser.EmitElemType_End:
            println("End:", event.Content)
            
        case streamable_parser.EmitElemType_Leaf:
            println("Leaf:", event.Content)
            
        case streamable_parser.EmitElemType_Error:
            println("Error:", event.Content)
            return
        }
        
        // Check if parsing is complete
        if event.Type == streamable_parser.EmitElemType_End && 
           event.Content == "expression" {
            break
        }
    }
}
```

## Event-Based Parsing

The parser emits events as it processes the input:

### Event Types

```go
type EmitElem struct {
    Type    EmitElemType // start, end, leaf, or error
    Content string       // non-terminal name or error message
    Leaf    *lexer.Token // token for leaf events
}
```

**Start Event**: Emitted when expanding a non-terminal
```
EmitElem{Type: Start, Content: "expression"}
```

**End Event**: Emitted when finishing a non-terminal
```
EmitElem{Type: End, Content: "expression"}
```

**Leaf Event**: Emitted when matching a terminal
```
EmitElem{Type: Leaf, Content: "+", Leaf: <token>}
```

**Error Event**: Emitted on syntax errors
```
EmitElem{Type: Error, Content: "parse error from 1,5 to 1,7"}
```

### Example Event Stream

For input `1 + 2`:

```
Start: expression
Start: term
Start: factor
Leaf: 1
End: factor
Leaf: +
Start: factor
Leaf: 2
End: factor
End: term
End: expression
```

## EBNF to BNF Conversion

The parser automatically converts EBNF constructs:

### Kleene Star (Zero or More)

```ebnf
// EBNF
statement_list -> statement*

// Converts to BNF
statement_list -> __repeat_statement_0
__repeat_statement_0 -> statement __repeat_statement_0 | ε
```

### Plus (One or More)

```ebnf
// EBNF  
digit_seq -> "DIGIT"+

// Converts to BNF
digit_seq -> "DIGIT" __repeat_digit_0
__repeat_digit_0 -> "DIGIT" __repeat_digit_0 | ε
```

### Grouping

```ebnf
// EBNF
expr -> term ( ( "+" | "-" ) term )*

// Converts to BNF with intermediate non-terminals
expr -> term __repeat_expr_0
__repeat_expr_0 -> __group_expr_1 term __repeat_expr_0 | ε
__group_expr_1 -> "+" | "-"
```

Artificial non-terminals (prefixed with `__`) are filtered out from events.

## FIRST and FOLLOW Sets

### FIRST Sets

FIRST(Aplha) = set of terminals that can begin strings derived from Aplha

**Computation Rules:**
1. If Aplha is a terminal, FIRST(Aplha) = {Aplha}
2. If Aplha is a non-terminal with production Aplha → β₁ β₂ ... βₙ:
   - Add FIRST(β₁) to FIRST(Aplha)
   - If β₁ can derive ε, add FIRST(β₂)
   - Continue until a symbol can't derive ε

### FOLLOW Sets

FOLLOW(A) = set of terminals that can appear immediately after A

**Computation Rules:**
1. Add $ (end-of-input) to FOLLOW(start symbol)
2. For production A → Aplha B β:
   - Add FIRST(β) - {ε} to FOLLOW(B)
   - If β can derive ε, add FOLLOW(A) to FOLLOW(B)
3. Repeat until no changes

### Example

```ebnf
expression -> term ("+" term)*
term       -> factor ("*" factor)*
factor     -> "NUMBER" | "(" expression ")"
```

**FIRST Sets:**
```
FIRST(expression) = { NUMBER, ( }
FIRST(term)       = { NUMBER, ( }
FIRST(factor)     = { NUMBER, ( }
```

**FOLLOW Sets:**
```
FOLLOW(expression) = { ), $ }
FOLLOW(term)       = { +, ), $ }
FOLLOW(factor)     = { *, +, ), $ }
```

## Error Recovery

The parser implements panic-mode error recovery:

### For Terminal Mismatch

When expected terminal doesn't match:
1. Emit error event with location range
2. Consume tokens until matching terminal is found
3. Continue parsing

### For Non-terminal Prediction Failure

When no production matches:
1. Emit error event with location range
2. Consume tokens until token in FOLLOW set is found
3. Pop non-terminal from stack
4. Continue parsing

This allows the parser to detect multiple errors in a single parse.

## Data Structures

### Grammar Representation

```go
type GrammarSequence struct {
    Elements []Grammar_element           // sequence of symbols
    FirstSet map[dfa.TokenType]struct{} // FIRST set for this sequence
}

type ProductionRule struct {
    Sequences []GrammarSequence          // alternative productions
    FollowSet map[dfa.TokenType]struct{} // FOLLOW set for this non-terminal
}

var grammarRules map[string]ProductionRule
```

### Stack Elements

```go
type StackElem struct {
    Type         StackElemType
    NonTermName  string         // for non-terminals
    TerminalType dfa.TokenType  // for terminals
}

const (
    StackElemType_Start  // expand non-terminal
    StackElemType_End    // finish non-terminal
    StackElemType_Leaf   // match terminal
)
```

## Performance Considerations

**Memory Usage:**
- Stack grows with nesting depth
- O(max_depth) memory for stack
- Grammar rules embedded in binary (no runtime parsing)

**Speed:**
- O(n) parsing time where n = number of tokens
- Constant-time stack operations
- Single token lookahead (no backtracking)

**Future Optimization:**
- Disk-backed stack for deep parsing (TODO in code)
- Would swap stack segments to disk when reaching 80% capacity

## Comparison with Recursive Descent Parser

| Feature | Streamable Parser | Recursive Descent |
|---------|------------------|-------------------|
| **Parsing Strategy** | Table-driven LL(1) | Recursive functions |
| **Stack** | Explicit stack | Call stack |
| **Memory** | O(depth) | O(depth) |
| **Streaming** | Event-based | Return AST |
| **Error Recovery** | Built-in | Manual |
| **Speed** | Fast | Very fast |
| **Code Size** | Compact | Large for big grammars |
| **Flexibility** | Standard LL(1) only | Can handle more grammars |

**Use Streamable When:**
- You need incremental parsing results
- Memory is constrained (can implement disk backing)
- Grammar is definitely LL(1)
- You want standard error recovery

**Use Recursive Descent When:**
- You need the full AST at once
- Grammar is complex (beyond LL(1))
- Call stack depth is not a concern
- Speed is critical

## Limitations

### LL(1) Grammar Restrictions

LL(1) parsers have specific requirements for the grammars they can parse. The following restrictions **must** be satisfied:

#### 1. No Left Recursion

**Left recursion** is when a non-terminal can derive a string that starts with itself:

```ebnf
// INVALID - Direct left recursion
expression -> expression "+" term

// INVALID - Indirect left recursion  
A -> B "x"
B -> A "y"
```

**Why it's problematic:** LL(1) parsers make decisions based on lookahead. Left recursion would cause infinite recursion when trying to expand the non-terminal.

**Solution:** Convert to right recursion or use repetition:
```ebnf
// VALID - Using repetition (EBNF)
expression -> term ("+" term)*

// VALID - Right recursion (BNF)
expression -> term expression_rest
expression_rest -> "+" term expression_rest | ε
```

#### 2. No Left Factoring Required

**Left factoring** is needed when multiple productions for the same non-terminal start with the same symbols:

```ebnf
// INVALID - Common prefix
statement -> "if" expression "then" statement "else" statement
statement -> "if" expression "then" statement

// INVALID - Another example
factor -> "(" expression ")"
factor -> "(" ")"
```

**Why it's problematic:** The parser can't decide which production to use with only one token of lookahead since both start with the same symbol.

**Solution:** Factor out the common prefix:
```ebnf
// VALID - Factored
statement -> "if" expression "then" statement else_part
else_part -> "else" statement | ε

// VALID - Factored
factor -> "(" optional_expression ")"
optional_expression -> expression | ε
```

#### 3. Disjoint FIRST Sets

For any non-terminal with multiple productions, the FIRST sets must be disjoint:

```ebnf
// INVALID - FIRST sets overlap
term -> "NUMBER"
term -> "NUMBER" "+" term
// Both have "NUMBER" in FIRST set
```

**Why it's problematic:** Parser can't determine which production to choose.

**Solution:** Restructure grammar or combine productions:
```ebnf
// VALID - Disjoint FIRST sets
term -> "NUMBER" optional_addition
optional_addition -> "+" term | ε
```

#### 4. Single Token Lookahead Sufficiency

The grammar must be parseable with only **one token of lookahead**. Some grammars inherently require more lookahead:

```ebnf
// Problematic - needs 2+ tokens lookahead
S -> "a" "b" "c"
S -> "a" "b" "d"
// Need to see both "a" and "b" to decide
```

### Grammar Requirements Summary

Your grammar **MUST**:
- ✅ Be free of left recursion (both direct and indirect)
- ✅ Be left-factored (no common prefixes in alternatives)
- ✅ Have disjoint FIRST sets for all production alternatives
- ✅ Be parseable with single token lookahead
- ✅ Have properly computed FOLLOW sets for epsilon productions

If your grammar violates any of these, the parser generator may:
- Fail to generate correct code
- Generate code that doesn't parse correctly
- Produce parsing conflicts at runtime

**Use the demo programs** (`sample_compute_firsts.go`, `sample_compute_follow.go`) to verify your grammar is LL(1) before generating a parser.

### Other Limitations

**Error Messages:**
- Currently basic location information
- Could be improved with better context

**Grammar Size:**
- Generated code size grows with grammar complexity
- Large grammars create large Go files

**No Ambiguity Resolution:**
- Unlike some parser generators, this doesn't automatically resolve conflicts
- You must manually fix grammar issues

## Future Improvements

1. **Better Error Messages**
   - Show source code snippet
   - Suggest likely fixes
   - Track multiple errors before stopping

2. **Disk-Backed Stack**
   - Handle extremely deep nesting
   - Swap inactive stack frames to disk

3. **Optimization**
   - Compress FIRST/FOLLOW sets
   - Use integers instead of strings for symbols
   - Pre-compile grammar to binary format

4. **Extended Features**
   - Support for attributes and actions
   - Semantic predicates
   - Parameterized non-terminals

## Demo Programs

See `demo/streamable_parser_demos/` for examples:
- `parser.demo.go` - Basic parsing demonstration
- `sample_compute_firsts.go` - FIRST set computation
- `sample_compute_follow.go` - FOLLOW set computation  
- `sample_ebnf_to_bnf.go` - Grammar conversion
- `sample_parser_writer.go` - Parser generation

## See Also

- [Parser Documentation](../parser/README.MD) - Alternative recursive descent parser
- [Lexer Documentation](../lexer/README.md) - Token generation
- [Demo Programs](../demo/streamable_parser_demos/) - Example usage

---

LL(1) parsing is elegant in its simplicity - a single token of lookahead is all you need!


