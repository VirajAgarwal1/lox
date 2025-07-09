package lexer_tests

import (
	"testing"

	dfa "github.com/VirajAgarwal1/lox/lexer/dfa"
)

// ----------------------------
// Helper Functions
// ----------------------------

func testDFAString(t *testing.T, inputDFA dfa.DFA, input string, expectedFinal dfa.DfaReturn, testName string) {
	t.Helper()

	inputDFA.Reset()
	lastResult := dfa.INVALID

	for i, r := range input {
		result := inputDFA.Step(r)
		lastResult = result
		t.Logf("%s: Step %d, char %q -> %s", testName, i, r, result.ToString())

		if result == dfa.INVALID {
			break
		}
	}

	if lastResult != expectedFinal {
		t.Errorf("%s: Expected %s, got %s for input %q", testName, expectedFinal.ToString(), lastResult.ToString(), input)
	}
}

func testDFASteps(t *testing.T, inputDFA dfa.DFA, input string, expected []dfa.DfaReturn, testName string) {
	t.Helper()
	inputDFA.Reset()

	if len(input) != len(expected) {
		t.Fatalf("%s: Input length %d != expected result length %d", testName, len(input), len(expected))
	}

	for i, r := range input {
		result := inputDFA.Step(r)
		if result != expected[i] {
			t.Errorf("%s: Step %d, char %q: expected %s, got %s", testName, i, r, expected[i].ToString(), result.ToString())
		}
	}
}

// ----------------------------
// Number DFA Tests
// ----------------------------

func TestNumberDFA(t *testing.T) {
	numberDFA := &dfa.NumberDFA{}
	numberDFA.Initialize()

	valid := []string{"0", "123", "9", "9999", "0.1", "123.456", "1.0"}
	intermediate := []string{"1.", "00083."}
	invalid := []string{"", ".", ".1", "a1", "1a", "1.2.3", "--1", "1 ", "1.1.1 111", "11 11.11", "1 1 . 11 1 1"}

	for _, input := range valid {
		t.Run("Valid_"+input, func(t *testing.T) {
			testDFAString(t, numberDFA, input, dfa.VALID, "NumberDFA_Valid")
		})
	}

	for _, input := range intermediate {
		t.Run("Intermediate_"+input, func(t *testing.T) {
			testDFAString(t, numberDFA, input, dfa.INTERMEDIATE, "NumberDFA_Intermediate")
		})
	}

	for _, input := range invalid {
		t.Run("Invalid_"+input, func(t *testing.T) {
			testDFAString(t, numberDFA, input, dfa.INVALID, "NumberDFA_Invalid")
		})
	}

	t.Run("Steps_123.45", func(t *testing.T) {
		expected := []dfa.DfaReturn{dfa.VALID, dfa.VALID, dfa.VALID, dfa.INTERMEDIATE, dfa.VALID, dfa.VALID}
		testDFASteps(t, numberDFA, "123.45", expected, "NumberDFA_123.45")
	})
}

// ----------------------------
// Identifier DFA Tests
// ----------------------------

func TestIdentifierDFA(t *testing.T) {
	identifierDFA := &dfa.IdentifierDFA{}
	identifierDFA.Initialize()

	valid := []string{"a", "_b", "a1", "CamelCase", "snake_case"}
	invalid := []string{"", "1a", " a", "_123abc", "a b", "a-b", "1_", "a _ b _", "🥰 😎", "🥰😎"}

	for _, input := range valid {
		t.Run("Valid_"+input, func(t *testing.T) {
			testDFAString(t, identifierDFA, input, dfa.VALID, "IdentifierDFA_Valid")
		})
	}
	for _, input := range invalid {
		t.Run("Invalid_"+input, func(t *testing.T) {
			testDFAString(t, identifierDFA, input, dfa.INVALID, "IdentifierDFA_Invalid")
		})
	}
}

// ----------------------------
// String DFA Tests
// ----------------------------

func TestStringDFA(t *testing.T) {
	stringDFA := &dfa.StringDFA{}
	stringDFA.Initialize()

	valid := []string{`""`, `"a"`, `"abcde wlekjfn pief"`, `"iiw oww 😇 😇😇😇😇😇"`}
	intermediate := []string{`"`, `" woiefw owk sej`, `"`, `"unclosed`}
	invalid := []string{"'single'", `"with\"quote"`, `"  www
	 " www"`, `222"efwf"`, `.    "   "`}

	for _, input := range valid {
		t.Run("Valid_"+input, func(t *testing.T) {
			testDFAString(t, stringDFA, input, dfa.VALID, "StringDFA_Valid")
		})
	}
	for _, input := range intermediate {
		t.Run("Intermediate_"+input, func(t *testing.T) {
			testDFAString(t, stringDFA, input, dfa.INTERMEDIATE, "StringDFA_Intermediate")
		})
	}
	for _, input := range invalid {
		t.Run("Invalid_"+input, func(t *testing.T) {
			testDFAString(t, stringDFA, input, dfa.INVALID, "StringDFA_Invalid")
		})
	}
}

// ----------------------------
// Comment DFA Tests
// ----------------------------

func TestCommentDFA(t *testing.T) {
	commentDFA := &dfa.CommentDFA{}
	commentDFA.Initialize()

	valid := []string{"//\n", "// comment\n", "//!@#$\n"}
	intermediate := []string{"/"}
	invalid := []string{"//\nnextline", " a//", "// comment\n nextline", "/ comm"}

	for _, input := range valid {
		t.Run("Valid_"+input, func(t *testing.T) {
			testDFAString(t, commentDFA, input, dfa.VALID, "CommentDFA_Valid")
		})
	}
	for _, input := range intermediate {
		t.Run("Intermediate_"+input, func(t *testing.T) {
			testDFAString(t, commentDFA, input, dfa.INTERMEDIATE, "CommentDFA_Intermediate")
		})
	}
	for _, input := range invalid {
		t.Run("Invalid_"+input, func(t *testing.T) {
			testDFAString(t, commentDFA, input, dfa.INVALID, "CommentDFA_Invalid")
		})
	}
}

// ----------------------------
// InputString DFA Tests
// ----------------------------

func TestInputStringDFA(t *testing.T) {
	cases := []struct {
		target string
		input  string
		expect dfa.DfaReturn
	}{
		{"if", "if", dfa.VALID},
		{"if", "i", dfa.INTERMEDIATE},
		{"==", "=", dfa.INTERMEDIATE},
		{"==", "==", dfa.VALID},
		{"==", "==x", dfa.INVALID},
		{"and", "and", dfa.VALID},
		{"and", "an", dfa.INTERMEDIATE},
		{"and", "andx", dfa.INVALID},
		{"or", "or", dfa.VALID},
		{"or", "o", dfa.INTERMEDIATE},
		{"or", "orx", dfa.INVALID},
		{"while", "while", dfa.VALID},
		{"while", "whil", dfa.INTERMEDIATE},
		{"while", "whilex", dfa.INVALID},
		{"for", "for", dfa.VALID},
		{"for", "fo", dfa.INTERMEDIATE},
		{"for", "forx", dfa.INVALID},
		{"else", "else", dfa.VALID},
		{"else", "el", dfa.INTERMEDIATE},
		{"else", "elsex", dfa.INVALID},
		{"function", "function", dfa.VALID},
	}

	for _, tc := range cases {
		t.Run(tc.target+"_"+tc.input, func(t *testing.T) {
			dfaForInputString := &dfa.InputStringDFA{}
			dfaForInputString.Initialize(tc.target)
			testDFAString(t, dfaForInputString, tc.input, tc.expect, "InputStringDFA")
		})
	}
}

// ----------------------------
// Miscellaneous Edge Tests
// ----------------------------

func TestResetFunctionality(t *testing.T) {
	numberDFA := &dfa.NumberDFA{}
	numberDFA.Initialize()
	numberDFA.Step('1')
	numberDFA.Step('2')
	numberDFA.Reset()
	if result := numberDFA.Step('3'); result != dfa.VALID {
		t.Errorf("Expected VALID after reset, got %s", result.ToString())
	}
}

// ----------------------------
// Benchmarks
// ----------------------------

func BenchmarkNumberDFA(b *testing.B) {
	numberDFA := &dfa.NumberDFA{}
	numberDFA.Initialize()

	for b.Loop() {
		numberDFA.Reset()
		for _, r := range "123.456" {
			numberDFA.Step(r)
		}
	}
}

// ----------------------------
// Whitespace DFA Tests
// ----------------------------

func TestWhitespaceDFA(t *testing.T) {
	whitespaceDFA := &dfa.WhitespaceDFA{}
	whitespaceDFA.Initialize()

	valid := []string{" ", "  ", "   ", "\t", "\t\t", " \t ", "\t \t", "    \t  \t  ", "\r"}
	invalid := []string{"", "a", " a", "a ", "\n", " \n", "\t\n", "1", " 1", "\t1", ".", " .", "\t."}

	for _, input := range valid {
		t.Run("Valid_"+input, func(t *testing.T) {
			testDFAString(t, whitespaceDFA, input, dfa.VALID, "WhitespaceDFA_Valid")
		})
	}

	for _, input := range invalid {
		t.Run("Invalid_"+input, func(t *testing.T) {
			testDFAString(t, whitespaceDFA, input, dfa.INVALID, "WhitespaceDFA_Invalid")
		})
	}

	// Test individual steps for mixed whitespace
	t.Run("Steps_SpaceTab", func(t *testing.T) {
		expected := []dfa.DfaReturn{dfa.VALID, dfa.VALID}
		testDFASteps(t, whitespaceDFA, " \t", expected, "WhitespaceDFA_SpaceTab")
	})
}

// ----------------------------
// Newline DFA Tests
// ----------------------------

func TestNewlineDFA(t *testing.T) {
	newlineDFA := &dfa.NewlineDFA{}
	newlineDFA.Initialize()

	valid := []string{"\n", "\n\n", "\n\n\n"}
	invalid := []string{"", " ", "\t", "a", "\na", "a\n", " \n", "\n ", "\t\n", "\n\t", "1\n", "\n1"}

	for _, input := range valid {
		t.Run("Valid_"+input, func(t *testing.T) {
			testDFAString(t, newlineDFA, input, dfa.VALID, "NewlineDFA_Valid")
		})
	}

	for _, input := range invalid {
		t.Run("Invalid_"+input, func(t *testing.T) {
			testDFAString(t, newlineDFA, input, dfa.INVALID, "NewlineDFA_Invalid")
		})
	}

	// Test individual steps for multiple newlines
	t.Run("Steps_DoubleNewline", func(t *testing.T) {
		expected := []dfa.DfaReturn{dfa.VALID, dfa.VALID}
		testDFASteps(t, newlineDFA, "\n\n", expected, "NewlineDFA_DoubleNewline")
	})
}

// ----------------------------
// Combined Whitespace Tests
// ----------------------------

func TestWhitespaceNewlineSeparation(t *testing.T) {
	// Test that whitespace and newline are properly separated
	whitespaceDFA := &dfa.WhitespaceDFA{}
	whitespaceDFA.Initialize()

	newlineDFA := &dfa.NewlineDFA{}
	newlineDFA.Initialize()

	// Whitespace should not accept newlines
	t.Run("Whitespace_RejectsNewline", func(t *testing.T) {
		testDFAString(t, whitespaceDFA, "\n", dfa.INVALID, "WhitespaceDFA_RejectsNewline")
	})

	// Newline should not accept spaces/tabs
	t.Run("Newline_RejectsSpace", func(t *testing.T) {
		testDFAString(t, newlineDFA, " ", dfa.INVALID, "NewlineDFA_RejectsSpace")
	})

	t.Run("Newline_RejectsTab", func(t *testing.T) {
		testDFAString(t, newlineDFA, "\t", dfa.INVALID, "NewlineDFA_RejectsTab")
	})
}
