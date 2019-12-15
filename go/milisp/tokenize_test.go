package milisp

import "fmt"

func prErr(_ interface{}, err error) {
	if err == nil {
		panic("Have to be error")
	}
	fmt.Println(err)
}

func prOk(x interface{}, err error) {
	if err != nil {
		panic(err)
	}
	fmt.Println(x)
}

func ExampleToken() {
	prOk(tokenize("abc x)\n))(y"))
	prOk(tokenize("a"))
	// Output:
	// [SYM:abc@1:1 SYM:x@1:5 END:)@1:6 END:)@2:1 END:)@2:2 BEG:(@2:3 SYM:y@2:4]
	// [SYM:a@1:1]
}

func ExampleToken_Quoted() {
	prOk(tokenize(`"abc"`))
	// Output: [STR:abc@1:1]
}

func ExampleToken_QuotedEscape() {
	prOk(tokenize(`"a\"c"`))
	// Output: [STR:a"c@1:1]
}

func ExampleToken_QuotedErrorEOF() {
	prErr(tokenize(`"`))
	// Output: unexpected EOF
}

func ExampleToken_SpaceErrorNakedSlash() {
	prErr(tokenize(`\`))
	// Output: unexpected char \ at 1:1
}

func ExampleToken_StringErrorQuote() {
	prErr(tokenize(`x"`))
	// Output: unexpected char " at 1:2
}

func ExampleToken_QuotedErrorSlashEOF() {
	prErr(tokenize(`"\`))
	// Output: unexpected EOF
}

func ExampleToken_CommentEOF() {
	prOk(tokenize(`x #`))
	// Output: [SYM:x@1:1]
}

func ExampleToken_CommentNewLine() {
	prOk(tokenize("x #\ny"))
	// Output: [SYM:x@1:1 SYM:y@2:1]
}

func ExampleToken_CommentBody() {
	prOk(tokenize("x #xx\ny"))
	// Output: [SYM:x@1:1 SYM:y@2:1]
}
