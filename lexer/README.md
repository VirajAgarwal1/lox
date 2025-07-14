# Lexer Package

A Go package for lexical analysis that converts source code into tokens using deterministic finite automata (DFAs). This is a complete lexical analyzer that I built because I find the intersection of compiler theory and practical implementation fascinating, and wanted to create something that gives maximum flexibility for tokenizing programming languages.

## Demo
[**Demo Video Link**](../scanner-demo.mp4)

## Performance

This lexer is optimized to minimize memory allocations and maximize throughput. It uses prebuilt DFA structures, reuses internal buffers across invocations, and keeps allocation counts per tokenization cycle very low. The result is a fast and lightweight scanner that can handle large inputs efficiently. Benchmarks show that it's suitable for real-time or batch lexing tasks, such as language tooling or compiler pipelines.

## What This Does

This package provides a `LexicalAnalyzer` that takes a stream of runes and produces a stream of tokens. It's designed to be the lexical analysis component of a compiler or interpreter, sitting between your source code reader and your parser.

## Processing Pipeline

The conversion follows this beautiful order:
**Bytes → Runes → Lexemes → Tokens**

The lexical analyzer sits in the middle of this pipeline, taking care of the Runes → Tokens conversion. It reads runes one at a time from your source, builds up lexemes (sequences of runes), and outputs complete tokens with all their metadata.

## How to Use the Scanner

Using the lexical analyzer is straightforward:

```go
package main

import (
    "bufio"
    "io"
    "strings"
    "your-module/lexer"
)

func main() {
    source := "while (x < 10) { x = x + 1; }"
    reader := bufio.NewReader(strings.NewReader(source))
    
    scanner := lexer.LexicalAnalyzer{}
    scanner.Initialize(reader)
    
    // Keep reading tokens until EOF
    for {
        token, err := scanner.ReadToken()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            break
        }
        
        fmt.Println(token.ToString())
    }
}
```

The `ReadToken()` function creates a beautiful pipe: input is a `RuneReader` from the buffered io package, output is tokens one at a time. This means you can easily pipe the output further into a parser or any other component that needs tokens.

## Token Structure

Each token contains:
- `TypeOfToken`: The type of token (keyword, identifier, operator, etc.)
- `lexemme`: The actual sequence of runes that formed the token
- `Line`: Line number in the source (starts at 0)
- `Offset`: Character offset within the line

## Conflict Resolution

The scanner handles conflicts intelligently using two strategies:

1. **Maximal Munching**: When multiple lexemes of different lengths are found, it always chooses the longest one. For example, `>=` beats `>` even though both are valid at the `>` character.

2. **Priority-Based Resolution**: When two lexemes have the same length, the scanner uses the order in an arbitrary array (from the underlying DFA package) as priority. Later tokens get higher priority, so keywords like `while` will match before generic identifiers.

## Supported Token Types

The lexer recognizes a complete set of tokens for a typical programming language:

- **Literals**: identifiers, strings, numbers, comments
- **Single-char tokens**: parentheses, braces, operators, punctuation
- **Multi-char tokens**: comparison operators (`==`, `!=`, `<=`, `>=`)
- **Keywords**: `if`, `while`, `for`, `class`, `fun`, `and`, `or`, etc.
- **Whitespace**: spaces, newlines

## Important Limitations

**Forward-Only Processing**: This scanner can only move forward through the source, one rune at a time. There's no capability to move the pointer backwards or jump around. If you need that kind of control, you'll have to handle it in the source reader you provide.

**The "Twice Processing" Issue**: There's a known inefficiency where the first rune of each token gets processed twice by all the DFAs. This happens because of how the scanner detects token boundaries - it needs to consume one more rune to know the previous token is complete, then save that rune for the next token. I've got a TODO to fix this, but it would require exposing DFA state management from the underlying package, which could complicate the clean interface. For now, it's a small price to pay for the simplicity.

**Edge Case Coverage**: While I've added a bunch of tests, I know I'm not covering all possible edge cases. It's hard to think of ALL the edge cases when you're building something like this. The code is solid for typical use cases, but if you use this in production, please test thoroughly with your specific input patterns.

## Error Handling

The lexer provides helpful error messages when it encounters invalid tokens:
- Line and offset information for precise error location
- When possible, it suggests the most likely intended token type based on partial matches
- Handles EOF gracefully, returning an EOF token when the source is exhausted

## Design Philosophy

This is setup in a way to maximize flexibility for how you might want to setup lexical analysis for your programming language. Obviously, if I had the time to make a regex parser that could generate DFAs automatically, that would make this much more flexible. But for now, this gives you a solid foundation with the ability to handle both simple fixed-string tokens and complex variable-length patterns.

I have added a bunch of tests as well to make sure this code works as expected. Though, I'm still not covering a lot of edge cases. It is hard to think of ALL the edge cases. So, for now I hope if *you* end up using this, you will use this with the same care with which I built it.

## Contributing

This is a personal learning project, but if you find bugs or have suggestions, I'd love to hear about them. The code is designed to be readable and extensible, so adding new token types should be straightforward (though you'd need to extend the underlying DFA package as well).

Happy tokenizing!