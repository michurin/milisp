package milisp

import (
	"fmt"
	"testing"
)

func TestTokenize_ok(t *testing.T) {
	for _, c := range []struct {
		text string
		res  string
	}{
		{
			text: "abc x)\n))(y",
			res:  "[SYM:abc@1:1 SYM:x@1:5 END:)@1:6 END:)@2:1 END:)@2:2 BEG:(@2:3 SYM:y@2:4]",
		},
		{
			text: "a",
			res:  "[SYM:a@1:1]",
		},
		{
			text: `"abc"`,
			res:  "[STR:abc@1:1]",
		},
		{
			text: `"a\"c"`,
			res:  `[STR:a"c@1:1]`,
		},
		{
			text: `x #`,
			res:  `[SYM:x@1:1]`,
		},
	} {
		c := c
		t.Run(c.text, func(t *testing.T) {
			p, err := tokenize(c.text)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
			if fmt.Sprintf("%v", p) != c.res {
				t.Errorf("Unexpected result: %v", p)
			}
		})
	}
}

func TestTokenize_invalid(t *testing.T) {
	for _, c := range []struct {
		text string
		err  string
	}{
		{
			text: `"`,
			err:  "unexpected EOF",
		},
		{
			text: `\`,
			err:  `unexpected char \ at 1:1`,
		},
		{
			text: `x"`,
			err:  `unexpected char " at 1:2`,
		},
		{
			text: `"\`,
			err:  "unexpected EOF",
		},
	} {
		c := c
		t.Run(c.text, func(t *testing.T) {
			p, err := tokenize(c.text)
			if err == nil {
				t.Errorf("Unexpected error: %s", err)
			}
			if p != nil {
				t.Errorf("Unexpected result: %v", p)
			}
			if err.Error() != c.err {
				t.Errorf("Unexpected error: %s", err)
			}
		})
	}
}
