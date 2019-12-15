package milisp_test

import (
	"fmt"
	"testing"

	"github.com/michurin/milisp/go/milisp"
)

func TestEvalFloat(t *testing.T) {
	expr, err := milisp.Prog("X")
	if err != nil {
		t.Error(err)
	}
	for _, c := range []struct {
		expr    milisp.Expression
		set     bool
		val     interface{}
		num     float64
		isError bool
	}{
		{expr, true, true, 1., false},
		{expr, true, false, 0., false},
		{expr, true, 100, 100., false},
		{expr, true, "100", 100., false},
		{expr, true, "x", 0., true},
		{expr, true, nil, 0., true},
		{expr, false, nil, 0., true},
		{nil, false, nil, 0., true},
	} {
		c := c
		t.Run(fmt.Sprintf("%v-%f", c.val, c.num), func(t *testing.T) {
			env := map[string]interface{}{}
			if c.set {
				env["X"] = c.val
			}
			f, err := milisp.EvalFloat(c.expr, env)
			if c.isError {
				if err == nil {
					t.Failed()
				}
			} else {
				if err != nil {
					t.Error(err)
				}
				if f != c.num {
					t.Error(f)
				}
			}
		})
	}
}

func TestEvalString(t *testing.T) {
	expr, err := milisp.Prog("X")
	if err != nil {
		t.Error(err)
	}
	for _, c := range []struct {
		expr    milisp.Expression
		set     bool
		val     interface{}
		str     string
		isError bool
	}{
		{expr, true, "A", "A", false},
		{expr, true, nil, "", true},
		{expr, false, nil, "", true},
		{nil, false, nil, "", true},
	} {
		c := c
		t.Run(fmt.Sprintf("%v-%s", c.val, c.str), func(t *testing.T) {
			env := map[string]interface{}{}
			if c.set {
				env["X"] = c.val
			}
			s, err := milisp.EvalString(c.expr, env)
			if c.isError {
				if err == nil {
					t.Failed()
				}
			} else {
				if err != nil {
					t.Error(err)
				}
				if s != c.str {
					t.Error(s)
				}
			}
		})
	}
}

func TestEvalCode(t *testing.T) {
	env := map[string]interface{}{}
	for _, text := range []string{
		")",  // parse error
		"()", // runtime error
	} {
		text := text
		t.Run(text, func(t *testing.T) {
			res, err := milisp.EvalCode(env, text)
			if res != nil {
				t.Failed()
			}
			if err == nil {
				t.Failed()
			}
		})
	}
}
