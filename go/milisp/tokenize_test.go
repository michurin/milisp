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
			text: "a",
			res:  "[SYM:a@1:1]",
		},
		{
			text: "a\nb\n\n c\rd\r\re\r\nf\n\rg\n\r\n\rh",
			res:  "[SYM:a@1:1 SYM:b@2:1 SYM:c@4:2 SYM:d@5:1 SYM:e@7:1 SYM:f@8:1 SYM:g@10:1 SYM:h@13:1]",
		},
		{
			text: "a\x0bb\x0b\x0b c\rd\r\re\r\x0bf\x0b\rg\x0b\r\x0b\rh",
			res:  "[SYM:a@1:1 SYM:b@2:1 SYM:c@4:2 SYM:d@5:1 SYM:e@7:1 SYM:f@8:1 SYM:g@10:1 SYM:h@13:1]",
		},
		{
			text: "a\n\tb\n1\tc\n1234567\td\n12345678\te\r\t\tf",
			res:  "[SYM:a@1:1 SYM:b@2:9 NUM:1@3:1 SYM:c@3:9 NUM:1234567@4:1 SYM:d@4:9 NUM:12345678@5:1 SYM:e@5:17 SYM:f@6:17]",
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
		{
			text: "abc x)\n))(y",
			res:  "[SYM:abc@1:1 SYM:x@1:5 END:)@1:6 END:)@2:1 END:)@2:2 BEG:(@2:3 SYM:y@2:4]",
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

func assertPanicFSM(f func(int, int) (int, int)) func(t *testing.T) {
	return func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("The code did not panic")
			}
		}()
		f(1<<30, 1<<30)
	}
}

func TestTokenize_panicFSM(t *testing.T) {
	t.Run("tokenizeStateTransitionFunction", assertPanicFSM(tokenizeStateTransitionFunction))
	t.Run("charPositionStateTransitionFunction", assertPanicFSM(charPositionStateTransitionFunction))
}
