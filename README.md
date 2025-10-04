# L O X 

A Go-based interpreter implementation for the Lox programming language, featuring custom-built lexer and parser components **built entirely from scratch**.

## A Journey Through Compiler Construction

This project represents a deep dive into the fascinating world of compiler theory, where I built every component from the ground up to truly understand how programming languages work. What started as curiosity about "how does code become execution?" became an intensive learning journey through lexical analysis, parsing theory, grammar design, and code generation.

### The Learning Experience

## What Makes This Project Special

### 1. **Hand-Crafted DFA-Based Lexer** 
I designed and implemented deterministic finite automata from scratch, not using any regex libraries or existing tokenizer frameworks. This meant:
- Researching DFA theory and state machine design
- Hand-coding state transitions for every token type
- Implementing maximal munching and conflict resolution algorithms
- Optimizing for performance with minimal memory allocations

### 2. **Recursive Descent Parser with Code Generation**
Building a parser is one thing - building a *parser generator* is another level entirely. Here's the meta-challenge: **to generate parsers, I first had to build a parser that parses grammar files themselves.**

This was a fascinating recursive problem:
- Built a complete lexer/parser specifically for grammar file syntax (EBNF notation)
- The grammar parser had to handle terminals, non-terminals, operators (`*`, `+`, `or`), grouping, and all EBNF constructs
- Implemented parser combinators for composable parsing logic
- Designed AST structures to represent grammar rules internally
- **Automatic code generation** that transforms parsed grammar into working Go parser code
- Proper operator precedence and associativity handling in the generated parsers

**The bootstrap problem:** I needed a working parser to parse grammars that define how to build parsers. This meant getting the grammar parser right was critical - any bugs would propagate to all generated parsers.

### 3. **Streamable LL(1) Predictive Parser** (The Heart of This Project)
**This was the main goal and the most ambitious undertaking of the entire project.** Making a parser that's truly *streamable* - that can incrementally consume tokens and emit parsing events without building the entire AST in memory - required far more planning, design, and research than I ever anticipated.

**Why streamable parsing is hard:**
- Traditional parsers build complete ASTs before returning anything
- Streaming means emitting events (start non-terminal, end non-terminal, match terminal) as parsing progresses
- Requires careful state management with an explicit stack
- Error recovery must work incrementally without losing position
- The event stream must be meaningful and usable for building parse trees incrementally

**The research journey:**
- Spent weeks diving deep into LL(1) parsing theory and predictive parsing algorithms
- Studied how to implement FIRST and FOLLOW set computation from academic papers
- Designed an event-based architecture from scratch that others could actually use
- Figured out how to handle error recovery in a streaming context
- Built an EBNF to BNF converter because the streaming parser needed pure BNF
- Tested and refined the streaming behavior to ensure it was actually useful

**Was it worth it?** Absolutely. The streaming parser can parse arbitrarily large inputs with bounded memory, handle incremental parsing scenarios, and emit events that can be processed in real-time. But more importantly, building it taught me more about parsing than any textbook could.

### 4. **Streamable Parser Code Generation** (The Ultimate Achievement)
If building a streamable parser was hard, building a *generator* for streamable parsers was the ultimate synthesis challenge:
- Had to build a parser for grammar files (the meta-parser mentioned above)
- Automatically converts EBNF to BNF (non-trivial transformation)
- Computes FIRST and FOLLOW sets algorithmically (translating theory to code)
- Generates efficient, table-driven parsing code with embedded parsing tables
- The generated code must handle streaming, error recovery, and event emission
- All of this happens automatically from a simple grammar file

**This required synthesizing everything learned about:**
- Parsing theory (LL(1), FIRST/FOLLOW, grammar transformations)
- Code generation (templates, AST manipulation, symbol tables)
- Compiler design (multi-pass processing, intermediate representations)
- The streaming architecture I spent so long designing

The generator took the longest to get right. Every edge case in grammar handling, every nuance of LL(1) restrictions, every detail of the streaming event protocol had to be encoded into the generator's logic. But when it finally worked - when I could write a grammar file and get a complete, working, streaming parser in return - that was an incredible moment.

### 5. **Performance Engineering**
Every component is optimized for real-world use:
- Minimal memory allocations in hot paths
- Efficient buffer management
- Reusable data structures
- Benchmarked and profiled for bottlenecks

### 6. **Comprehensive Documentation & Demos**
Because understanding is as important as building:
- Detailed READMEs explaining theory and implementation
- Working demo programs for every major component
- Complete test suites
- Example grammars and sample code

## The Research Behind It

This project required extensive research across multiple domains:
- **Automata Theory**: DFA construction, state minimization, conflict resolution
- **Parsing Theory**: LL(1) grammars, recursive descent, predictive parsing, FIRST/FOLLOW sets
- **Grammar Design**: EBNF notation, left recursion elimination, left factoring, grammar transformations
- **Code Generation**: AST construction, symbol table management, template-based generation, meta-programming
- **Language Implementation**: Operator precedence, error recovery, streaming architectures
- **Streaming Parser Design**: Event-based architectures, incremental parsing, bounded memory consumption, state management

The streaming parser alone required:
- Reading multiple academic papers on LL(1) parsing algorithms
- Studying existing streaming XML parsers (SAX) for event design patterns
- Experimenting with different stack representations and event protocols
- Weeks of planning how to make the streaming behavior actually useful
- Extensive testing with different grammar patterns to find edge cases

Each feature represents countless hours of reading, planning, implementing, debugging, and refining. The streamable parser generator was particularly intense - it's one thing to understand how parsers work, but quite another to build a system that automatically generates them with all the correct streaming behavior, error recovery, and event emission.

## Architecture

The interpreter follows the classic compiler pipeline:

```
Source Code (UTF-8)
    ↓
Lexer (DFA-based) → Tokens
    ↓
Parser (Recursive Descent or Predictive Parser) → AST
    ↓
[Interpreter - Future Work]
```

## Project Structure

### Core Components

- **`lexer/`** - Lexical analysis converting source code to tokens
  - DFA-based tokenization using deterministic finite automata
  - Supports all Lox token types: keywords, operators, literals, identifiers
  - Maximal munching and priority-based conflict resolution
  - See [lexer/README.md](lexer/README.md) for details

- **`parser/`** - Recursive descent parser with operator precedence
  - Generates parsers from EBNF grammar specifications
  - Handles expressions with proper precedence and associativity
  - Parser combinator-based implementation
  - See [parser/README.MD](parser/README.MD) for details

- **`streamable_parser/`** - Predictive parser
  - Predictive parser with streaming token consumption
  - Automatic FIRST and FOLLOW set computation
  - EBNF to BNF conversion
  - Event-based parsing with start/end/leaf emissions
  - See [streamable_parser/README.md](streamable_parser/README.md) for details

- **`errorhandler/`** - Error reporting utilities
  - Stack trace generation with file and line information
  - Error wrapping and context preservation
  - See [errorhandler/README.md](errorhandler/README.md) for details

### Supporting Directories

- **`demo/`** - Demonstration programs showing component usage
  - Lexer demos, parser demos, and streamable parser examples
  - See [demo/README.md](demo/README.md) for usage examples

- **`tests/`** - Comprehensive test suite
  - Unit tests for lexer, parser, and streamable parser
  - Test fixtures and sample Lox programs
  - See [tests/README.md](tests/README.md) for running tests

## Quick Start

### Prerequisites

- Go 1.24.4 or later
- UTF-8 encoded source files

### Installation

```bash
git clone https://github.com/VirajAgarwal1/lox.git
cd lox
go mod download
```

### Running Demos

**⚠️ Important for Parser Generation Demos:**

Parser generation requires **TWO separate runs** because Go compiles at build-time:
1. **First run:** Generates the parser code (writes `generated_parser.go`)
2. **Second run:** Uses the generated parser (compiles and runs the new code)

See detailed instructions in [demo/README.md](demo/README.md).

```bash
# Step 1: Generate parser code
go run main.go  # with generation function

# Step 2: Use the generated parser (modify main.go first)
go run main.go  # with usage function
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./lexer/...
go test ./parser/...
go test ./streamable_parser/...
```

## Current Status

✅ **Completed (After Extensive Research & Development):**
- ✨ **Hand-crafted DFA-based lexical analyzer** - Built from scratch with custom state machines
- ✨ **Recursive descent parser generator** - Reads grammars, generates working parsers
- ✨ **LL(1) predictive parser** - Complete implementation with FIRST/FOLLOW computation
- ✨ **Predictive parser code generator** - The ultimate achievement: generates entire LL(1) parsers from grammar specs
- ✨ **EBNF to BNF converter** - Automatic grammar transformation
- ✨ **Error handling framework** - Stack traces and error propagation
- ✨ **Comprehensive test suites** - Ensuring correctness at every level
- ✨ **Extensive documentation** - Because understanding matters

🚧 **In Progress:**
- None currently (this represents the culmination of the learning journey so far)

📋 **Future Work:**
- AST interpreter/evaluator
- Semantic analysis
- Runtime environment
- Type checking
- Optimization passes

## Important Limitations

**Character Encoding**: Currently, the source code must be in UTF-8 encoding. Other encodings may cause unexpected behavior.

**Forward-Only Processing**: The lexer can only move forward through the source. No backtracking or random access is supported at the scanner level.

**Parser Generation**: Code generation happens at build-time in Go. You cannot generate and use parser code in the same execution - you must run your program twice (once to generate, once to use).

**Edge Cases**: While extensively tested, not all possible edge cases are covered. Use with appropriate care in production environments.

## Design Philosophy

This project prioritizes:
1. **Deep Understanding** - Building everything from scratch to truly understand how it works
2. **Learning Through Doing** - Clear, readable code over clever optimizations so others can learn too
3. **Flexibility** - Extensible architecture supporting multiple parsing strategies
4. **Performance** - Efficient where it matters (minimal allocations, streaming)
5. **Knowledge Sharing** - Extensive documentation because teaching reinforces learning

## Why Build This?

**Because the best way to learn something is to build it yourself.**

Using existing parser generators like YACC or ANTLR is fine for production work, but you don't truly understand parsing theory until you've:
- Debugged why a grammar is ambiguous
- Implemented FIRST/FOLLOW set computation by hand
- Watched your DFAs accept and reject tokens
- Generated working code from grammar specifications
- Dealt with left recursion and left factoring issues
- Designed a streaming architecture that actually works
- Built a parser to parse the grammars that define how to build parsers

**The main challenge: Making it streamable.** I initially thought "I'll build a parser" but quickly realized I wanted something more ambitious - a parser that could handle massive inputs efficiently, emit events incrementally, and not require loading entire parse trees into memory. This single design goal drove weeks of additional research and implementation. Was it harder than I expected? Absolutely. Did I learn exponentially more because of it? Without question.

This project is proof that with curiosity, persistence, and a lot of research, you can build complex systems from first principles. Every bug fixed taught me something. Every feature implemented deepened my understanding. Every design decision for the streaming architecture revealed new insights into how parsers work. The streaming parser, in particular, required a level of planning and forethought I hadn't anticipated - but solving those challenges was incredibly rewarding.

**If you're learning about compilers, I encourage you to build something like this too.** Pick a challenging goal (like streaming), commit to understanding it deeply, and iterate until it works. The journey is just as valuable as the destination.

## Future Improvements

### Lexer

1. **DFA Optimization**
   - Trie-based approach for exact keyword matching would be faster
   - Token grouping by character class (digits, letters, operators) before DFA dispatch
   
2. **Scanner Efficiency**
   - Skip DFAs that return INVALID from subsequent token processing
   - Fix the "twice processing" issue where first rune of each token is processed twice

### Parser

1. **Error Recovery**
   - Better error messages with suggestions
   - Panic mode recovery for syntax errors
   
2. **Performance**
   - Memoization for parser combinators
   - Lazy evaluation strategies

### General

1. **Regex-based DFA Generation**
   - Automatically generate DFAs from regex patterns
   - Would make token definitions much more flexible

## What I Learned

Building this project taught me:
- **Automata Theory**: How DFAs work in practice, not just theory
- **Parsing Algorithms**: The subtle differences between LL(1), recursive descent, and other parsing strategies
- **Grammar Design**: Why some grammars work and others don't, and how to fix them
- **Code Generation**: How to write programs that write programs (and parsers that parse grammar definitions)
- **Algorithm Implementation**: Translating academic papers into working code
- **Systems Thinking**: How all the pieces fit together in a compiler pipeline
- **Streaming Architecture Design**: How to design event-based systems that maintain state correctly
- **Debugging Complex Systems**: Finding bugs in generated code is a unique challenge
- **Performance Optimization**: Where optimizations matter and where they don't
- **Planning and Design**: The streaming parser taught me that upfront design is crucial for complex systems
- **Persistence**: When building the streaming parser took far longer than expected, pushing through was worth it
- **Documentation**: How to explain complex topics clearly

**Most importantly:** I learned that I *can* understand and build complex systems if I'm willing to put in the time and effort. The streaming parser was the perfect example - it seemed impossibly complex at first, required way more research and planning than I anticipated, but breaking it down into smaller problems and iterating steadily made it achievable. And the depth of understanding gained from building it from scratch is something no tutorial could provide.

## Contributing

This is primarily a **learning project**, but that's exactly why contributions are welcome! 

If you're also learning about compilers:
- Try building a feature
- Add tests
- Improve documentation
- Share what you learned

Bug reports and suggestions are always appreciated. The code is designed to be readable and extensible specifically so others can learn from it.

---

Happy parsing! 🚀