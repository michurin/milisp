package milisp

import "fmt"

func Prog(text string) (Expression, error) {
	tokens, err := tokenize(text)
	if err != nil {
		return nil, fmt.Errorf("tokenizer error: %s", err)
	}
	expr, pos, err, _ := parse(tokens, 0) // we can drop finish-flag, err is enough
	if err != nil {
		return nil, fmt.Errorf("parser error: %s", err)
	}
	if pos < len(tokens) {
		return nil, fmt.Errorf("extra content after token %s", tokens[pos])
	}
	return expr, nil
}
