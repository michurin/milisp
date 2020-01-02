package milisp

import "testing"

func TestImpossibleCase(t *testing.T) {
	u := universalToken{
		tp:   tpOpen,
		num:  0,
		str:  "",
		line: 0,
		pos:  0,
	}
	r, err := u.Eval(Environment{})
	if r != nil {
		t.Failed()
	}
	if err == nil {
		t.Failed()
	}
}
