# L O X 

A Go-based interpreter implementation for the Lox programming language, featuring custom-built lexer and parser components **built entirely from scratch**.

## Project Origin

This project started as an implementation of the Lox programming language from Robert Nystrom's excellent book ["Crafting Interpreters"](https://craftinginterpreters.com/). However, I wanted to take a different approach: **building the language in Go with a streaming architecture from the ground up**. Rather than following the book's implementation directly, I used it as inspiration to build every component from scratch—from DFA-based lexers to LL(1) parser generators—focusing on creating parsers that handle inputs incrementally without loading entire parse trees into memory.

## What Makes This Project Special

### 1. **Hand-Crafted DFA-Based Lexer** 
Built deterministic finite automata from scratch without regex libraries—hand-coding state transitions for every token type, implementing maximal munching and conflict resolution, optimized for minimal memory allocations.

### 2. **Recursive Descent Parser with Code Generation**
To generate parsers, I first built a parser that parses grammar files themselves—the classic bootstrap problem. The system reads EBNF grammar specifications, handles all constructs (terminals, non-terminals, operators, grouping), and automatically generates working Go parser code with proper precedence and associativity.

### 3. **Streamable LL(1) Predictive Parser** (The Heart of This Project)
The main goal: a parser that incrementally consumes tokens and emits parsing events without building the entire AST in memory. This required weeks of research into LL(1) theory, FIRST/FOLLOW set computation from academic papers, and designing an event-based architecture with proper state management and error recovery. The result: a parser that handles arbitrarily large inputs with bounded memory and processes events in real-time.

### 4. **Streamable Parser Code Generation** (The Ultimate Achievement)
The ultimate synthesis: a system that reads grammar files, automatically converts EBNF to BNF, computes FIRST/FOLLOW sets, and generates complete table-driven LL(1) parsers with streaming, error recovery, and event emission—all from a simple grammar specification. This required combining everything learned about parsing theory, code generation, and streaming architecture design.

### 5. **Comprehensive Documentation & Demos**
Because understanding is as important as building:
- Detailed READMEs explaining theory and implementation
- Working demo programs for every major component
- Complete test suites
- Example grammars and sample code

## Demo Videos - See It In Action! 

Want to see how this all works? Watch these video demonstrations directly below:

### 📹 **DFA-Based Tokenization in Action**
Watch how deterministic finite automata process input character-by-character, see state transitions and token recognition in real-time, and understand how maximal munching and conflict resolution work.

<video width="100%" controls>
  <source src="dfa-demo.mp4" type="video/mp4">
  Your browser does not support the video tag. <a href="dfa-demo.mp4">Download the video</a>
</video>

### 📹 **Complete Lexical Analysis**
Full demonstration of the lexer converting source code to tokens. See how keywords, operators, literals, and identifiers are recognized, and follow tokens through the scanning pipeline.

<video width="100%" controls>
  <source src="scanner-demo.mp4" type="video/mp4">
  Your browser does not support the video tag. <a href="scanner-demo.mp4">Download the video</a>
</video>

### 📹 **Automatic Parser Generation**
Watch the parser generator read a grammar file and produce working parser code. See EBNF to BNF conversion, FIRST/FOLLOW set computation in action, and observe complete LL(1) parser generation from start to finish.

<video width="100%" controls>
  <source src="streamable_parser%20code%20generation%20demo.mp4" type="video/mp4">
  Your browser does not support the video tag. <a href="streamable_parser%20code%20generation%20demo.mp4">Download the video</a>
</video>

### 📹 **Streaming Parser in Action**
See the generated LL(1) parser process expressions incrementally. Watch event-based parsing with start/end/leaf emissions and follow the parse tree construction in real-time.

<video width="100%" controls>
  <source src="streamable_parser%20parsing%20demo.mp4" type="video/mp4">
  Your browser does not support the video tag. <a href="streamable_parser%20parsing%20demo.mp4">Download the video</a>
</video>

These videos provide visual walkthroughs of the core components and are a great way to quickly understand what this project does!

## The Research Behind It

Building this required deep dives into automata theory (DFA construction, conflict resolution), parsing theory (LL(1), FIRST/FOLLOW sets), grammar design (EBNF, left factoring), code generation (AST manipulation, templates), and streaming architecture design (event-based systems, incremental parsing). The streaming parser particularly required studying academic papers on LL(1) algorithms and existing streaming parsers (SAX) for event design patterns.

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

**Important for Parser Generation Demos:**

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

**Completed (After Extensive Research & Development):**
- **Hand-crafted DFA-based lexical analyzer** - Built from scratch with custom state machines
- **Recursive descent parser generator** - Reads grammars, generates working parsers
- **LL(1) predictive parser** - Complete implementation with FIRST/FOLLOW computation
- **Predictive parser code generator** - The ultimate achievement: generates entire LL(1) parsers from grammar specs
- **EBNF to BNF converter** - Automatic grammar transformation
- **Error handling framework** - Stack traces and error propagation
- **Comprehensive test suites** - Ensuring correctness at every level
- **Extensive documentation** - Because understanding matters

**In Progress:**
- None currently (this represents the culmination of the learning journey so far)

**Future Work:**
- AST interpreter/evaluator
- Semantic analysis
- Runtime environment
- Type checking
- Optimization passes

## Important Limitations

**Character Encoding**: Currently, the source code must be in UTF-8 encoding. Other encodings may cause unexpected behavior.

**Forward-Only Processing**: The lexer can only move forward through the source. No backtracking or random access is supported at the scanner level.

**Parser Generation**: Code generation happens at build-time in Go. You cannot generate and use parser code in the same execution - you must run your program twice (once to generate, once to use).

**Edge Cases**: While extensively tested, not all possible edge cases are covered. Use with appropriate care.

## Design Philosophy

This project prioritizes:
1. **Deep Understanding** - Building everything from scratch to truly understand how it works
2. **Learning Through Doing** - Clear, readable code over clever optimizations so others can learn too
3. **Flexibility** - Extensible architecture supporting multiple parsing strategies
4. **Performance** - Efficient where it matters (minimal allocations, streaming)
5. **Knowledge Sharing** - Extensive documentation because teaching reinforces learning

## Why Build This?

**Because the best way to learn something is to build it yourself.**

You don't truly understand parsing until you've debugged ambiguous grammars, implemented FIRST/FOLLOW computation by hand, watched your DFAs work, and built parsers that generate parsers. The streaming architecture was particularly challenging—handling massive inputs efficiently with bounded memory required weeks of additional research, but taught me more than any textbook could.

This project proves that with curiosity and persistence, you can build complex systems from first principles. If you're learning about compilers, pick a challenging goal, commit to understanding it deeply, and iterate until it works.

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

Happy parsing!