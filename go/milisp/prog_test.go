package milisp_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michurin/milisp/go/milisp"
)

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

func ExampleParse() {
	prOk(milisp.Prog("A"))
	prOk(milisp.Prog("(A())")) // no space
	prOk(milisp.Prog("(A)"))
	prOk(milisp.Prog("(A(X Y)A)"))
	prOk(milisp.Prog("((X Y))"))
	prErr(milisp.Prog("A B"))
	prErr(milisp.Prog("\\"))
	prErr(milisp.Prog("("))
	prErr(milisp.Prog("(()"))
	// Output:
	// SYM:A@1:1
	// [SYM:A@1:2 []@1:3]@1:1
	// [SYM:A@1:2]@1:1
	// [SYM:A@1:2 [SYM:X@1:4 SYM:Y@1:6]@1:3 SYM:A@1:8]@1:1
	// [[SYM:X@1:3 SYM:Y@1:5]@1:2]@1:1
	// extra content after token SYM:B@1:3
	// tokenizer error: unexpected char \ at 1:1
	// parser error: unexpected end of file
	// parser error: unexpected end of file; can not end expr started at BEG:(@1:1
}

func evalToString(text string, env milisp.Env) (string, error) {
	e, err := milisp.Prog(text)
	if err != nil {
		panic(err)
	}
	return milisp.EvalString(e, env)
}

func ExampleEvalUnknownSymbol() {
	emptyEnv := map[string]interface{}{}
	prErr(evalToString("(X)", emptyEnv))
	prErr(evalToString("()", emptyEnv))
	prErr(evalToString("(X)", map[string]interface{}{"X": nil}))
	prErr(evalToString("(X)", map[string]interface{}{
		"X": milisp.OpFunc(func(_ milisp.Env, _ []milisp.Expression) (interface{}, error) {
			return nil, errors.New("message")
		})}))
	// Output:
	// runtime error: unknown symbol: SYM:X@1:2
	// empty expression
	// operation <nil> not executable: SYM:X@1:2
	// message
}

func ExampleEvalConcat() {
	e, err := milisp.Prog(`(concat "A" (concat VAR "Q") "B")`)
	if err != nil {
		panic(err)
	}
	res, err := milisp.EvalString(e, map[string]interface{}{
		"concat": milisp.OpFunc(func(env milisp.Env, e []milisp.Expression) (interface{}, error) {
			s := make([]string, len(e)-1)
			for i, x := range e[1:] {
				s[i], err = milisp.EvalString(x, env)
				if err != nil {
					return nil, err
				}
			}
			return "<" + strings.Join(s, ",") + ">", nil
		}),
		"VAR": "P",
	})
	prOk(res, err)
	// Output: <A,<P,Q>,B>
}

func ExampleEvalMath() {
	e, err := milisp.Prog("(* 2 (+ x y z))")
	if err != nil {
		panic(err)
	}
	res, err := milisp.EvalFloat(e, map[string]interface{}{
		"+": milisp.OpFunc(func(env milisp.Env, e []milisp.Expression) (interface{}, error) {
			s := float64(0)
			for _, x := range e[1:] {
				t, err := milisp.EvalFloat(x, env)
				if err != nil {
					return nil, err
				}
				s += t
			}
			return s, nil
		}),
		"*": milisp.OpFunc(func(env milisp.Env, e []milisp.Expression) (interface{}, error) {
			if len(e) < 3 {
				return nil, fmt.Errorf("too few args")
			}
			m := float64(1)
			for _, x := range e[1:] {
				t, err := milisp.EvalFloat(x, env)
				if err != nil {
					return nil, err
				}
				m *= t
			}
			return m, nil
		}),
		"x": 1,
		"y": 2,
		"z": 3,
	})
	prOk(res, err)
	// Output: 12
}
