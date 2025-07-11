# DFA Package

A Go package for building deterministic finite automata (DFAs) to handle lexical analysis in programming language scanners. I built this because I find the intersection of compiler theory and practical implementation fascinating, and wanted to create something that gives maximum flexibility for different lexical analysis needs.

## What This Does

In this package I have written DFAs of 2 types:

1. **Fixed String DFAs**: These take in a string and then can check whether a sequence of characters (input in the DFA one at a time) match the string with which it was initialized. I use this for lexemes with fixed length and runes - things like keywords (`if`, `while`, `class`) and operators (`==`, `!=`, `<=`).

2. **Hand-written DFAs**: Custom DFAs for other types of lexemes which are more variable like whitespace, newlines, comments, identifiers, and strings. These handle the more complex patterns that need custom logic.

The way to use this is to pipe the sequence of runes into all the DFAs at once, and see which lexemes' DFA gives `VALID` output. The DFAs are setup so that if even one rune is not the right one, then further DFA steps taken will always output `INVALID`. The DFAs also output `INTERMEDIATE` sometimes when the sequence of runes could become the respective lexeme given the right sequence of runes follow after, but the current sequence doesn't actually satisfy the lexeme DFA yet.

This is setup in a way to maximize the flexibility for how one might want to setup the scanner (lexical analyzer) for their programming language. Obviously, though if I had the time to make a regex parser itself, then that would have made this much more flexible. I have added a bunch of tests as well to make sure this code is working as expected. Though, I'm still not covering a lot of edge cases. It is hard to think of ALL the edge cases. So, for now I hope if *you* end up using this, you will use this with the same care with which I built it.

## Core Interface

```go
type DFA interface {
    Step(input rune) DfaReturn
    Reset()
}
```

Every DFA implements this simple interface:
- `Step(input rune)` - Feed the DFA a single rune and get back the current state
- `Reset()` - Reset the DFA back to its initial state

## Return Values

```go
const (
    INVALID DfaReturn = iota
    INTERMEDIATE
    VALID
)
```

- **INVALID**: The current sequence cannot possibly match this lexeme
- **INTERMEDIATE**: The sequence so far is valid but incomplete - more runes needed
- **VALID**: The current sequence is a complete, valid match for this lexeme

## Token Types

The package includes DFAs for a complete set of tokens you'd expect in a typical programming language:

- **Literals**: identifiers, strings, numbers, comments
- **Single-char tokens**: parentheses, braces, operators, punctuation
- **Multi-char tokens**: comparison operators (`==`, `!=`, `<=`, `>=`)
- **Keywords**: `if`, `while`, `for`, `class`, `fun`, etc.
- **Whitespace**: spaces, newlines

## Usage Example

```go
package main

import "your-module/dfa"

func main() {
    // Generate all DFAs for the supported token types
    dfas := dfa.GenerateDFAs()
    
    // Test input
    input := "return"
    
    // Feed each character to all DFAs
    for i, char := range input {
        for tokenType, automaton := range dfas {
            result := automaton.Step(char)
            if result.IsValid() {
                fmt.Printf("Token '%s' matches %s\n", input[:i+1], tokenType)
            }
        }
    }
}
```

## How It Works

The `GenerateDFAs()` function creates a map of `TokenType` to `DFA` implementations. For fixed strings (like keywords and operators), it uses `InputStringDFA`. For more complex patterns, it uses specialized DFAs like `IdentifierDFA`, `StringDFA`, `NumberDFA`, etc.

When you're building a scanner, you'd typically:
1. Get all the DFAs using `GenerateDFAs()`
2. For each input character, call `Step()` on all active DFAs
3. Track which DFAs are still valid/intermediate
4. When a DFA returns `VALID`, you've found a complete token
5. Reset DFAs as needed and continue

## A Word of Caution

This is a hobby project born out of my interest in compiler construction. While I've tested it reasonably well, I haven't covered every edge case imaginable. If you decide to use this in your own projects, please test thoroughly and use it with the same care I put into building it.

The real dream would be to extend this with a proper regex parser that could generate DFAs automatically, but that's a project for another day. For now, this gives you a solid foundation for lexical analysis with the flexibility to handle both simple fixed-string tokens and complex variable-length patterns.

## Contributing

This is a personal learning project, but if you find bugs or have suggestions, I'd love to hear about them. The code is designed to be readable and extensible, so adding new token types or DFA patterns should be straightforward. If not, though you can contact me on my email or open an issue!