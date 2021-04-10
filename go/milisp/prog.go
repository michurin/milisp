package milisp

import "fmt"

// Compile text of lisp program to ready to run internal representation.
func Compile(text string) (Expression, error) {
	tokens, err := tokenize(text)
	if err != nil {
		return nil, fmt.Errorf("tokenizer error: %s", err)
	}
	expr, pos, _, err := parse(tokens, 0) // we can drop finish-flag, err is enough
	if err != nil {
		return nil, fmt.Errorf("parser error: %s", err)
	}
	if pos < len(tokens) {
		return nil, fmt.Errorf("extra content after token %s", tokens[pos])
	}
	return expr, nil
}
