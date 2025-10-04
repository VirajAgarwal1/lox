/*
# NOTE:
	Currently, sourcecode of Lox only supports reading it in **UTF-8** encoding. If provided sourcecode in other encodings then the program may behave in unexpected ways..
*/

// DEMO CODE
package main

import (
	streamable_parser_demos "github.com/VirajAgarwal1/lox/demo/streamable_parser_demos"
)

func main() {
	// # IMPORTANT: Parser generation requires TWO separate runs!
	// Go compiles code at build-time, so you can't generate and use code in the same execution.
	//
	// Step 1: Run this to GENERATE the parser:
	// streamable_parser_demos.Sample_generate_streamable_parser()
	//
	// Step 2: Comment out the line above, uncomment the line below, and run AGAIN:
	streamable_parser_demos.Sample_streamable_parser_demo()

	// Other useful demos:
	// streamable_parser_demos.Sample_compute_firsts()
	// streamable_parser_demos.Sample_compute_follow()
	// streamable_parser_demos.Sample_ebnf_to_bnf()
}
