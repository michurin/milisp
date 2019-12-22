package milisp_test

import (
	"fmt"

	"github.com/michurin/milisp/go/milisp"
)

func ExampleCompile_validSyntax() {
	for _, text := range []string{
		"A",
		"(A())", // no space
		"(A)",
		"(A(X Y)A)",
		"((X Y))",
	} {
		prog, err := milisp.Compile(text)
		if err != nil {
			panic(err)
		}
		fmt.Println(prog)
	}
	// Output:
	// SYM:A@1:1
	// [SYM:A@1:2 []@1:3]@1:1
	// [SYM:A@1:2]@1:1
	// [SYM:A@1:2 [SYM:X@1:4 SYM:Y@1:6]@1:3 SYM:A@1:8]@1:1
	// [[SYM:X@1:3 SYM:Y@1:5]@1:2]@1:1
}

func ExampleCompile_invalidSyntax() {
	for _, text := range []string{
		"A B",
		"\\",
		"(",
		"(()",
	} {
		prog, err := milisp.Compile(text)
		if err == nil || prog != nil {
			panic(text)
		}
		fmt.Println(err)
	}
	// Output:
	// extra content after token SYM:B@1:3
	// tokenizer error: unexpected char \ at 1:1
	// parser error: unexpected end of file
	// parser error: unexpected end of file; can not end expr started at BEG:(@1:1
}
