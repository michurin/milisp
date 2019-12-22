package milisp_test

import (
	"errors"
	"fmt"

	"github.com/michurin/milisp/go/milisp"
)

func ExampleEnvironment_oneAtomPrograms() {
	prOk := func(result interface{}, err error) {
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
	emptyEnv := map[string]interface{}{"X": "VALUE_OF_X"}
	prOk(milisp.EvalCode(emptyEnv, "X"))   // symbol: value will be taken from environment
	prOk(milisp.EvalCode(emptyEnv, `"X"`)) // string
	prOk(milisp.EvalCode(emptyEnv, "1"))   // float
	// Output:
	// VALUE_OF_X
	// X
	// 1
}

func ExampleEnvironment_symbolErrors() {
	prErr := func(_ interface{}, err error) {
		if err == nil {
			panic("Have to be error")
		}
		fmt.Println(err)
	}
	emptyEnv := map[string]interface{}{}
	prErr(milisp.EvalCode(emptyEnv, "(X)"))
	prErr(milisp.EvalCode(emptyEnv, "()"))
	prErr(milisp.EvalCode(map[string]interface{}{"X": nil}, "(X)"))
	prErr(milisp.EvalCode(map[string]interface{}{ // symbol exists, however evaluate with internal error
		"X": milisp.OpFunc(func(_ milisp.Environment, _ []milisp.Expression) (interface{}, error) {
			return nil, errors.New("error message")
		})}, "(X)"))
	// Output:
	// runtime error: unknown symbol: SYM:X@1:2
	// empty expression
	// operation <nil> not executable: SYM:X@1:2
	// error message
}

func ExampleOpFunc_valuesAndExecutables() {
	e, err := milisp.Compile("(* 2 (+ x y z))")
	if err != nil {
		panic(err)
	}
	env := map[string]interface{}{
		"+": milisp.OpFunc(func(env milisp.Environment, e []milisp.Expression) (interface{}, error) {
			s := float64(0)
			for _, x := range e[1:] {
				t, err := milisp.EvalFloat(env, x)
				if err != nil {
					return nil, err
				}
				s += t
			}
			return s, nil
		}),
		"*": milisp.OpFunc(func(env milisp.Environment, e []milisp.Expression) (interface{}, error) {
			if len(e) < 3 {
				return nil, fmt.Errorf("too few args")
			}
			m := float64(1)
			for _, x := range e[1:] {
				t, err := milisp.EvalFloat(env, x)
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
	}
	res, err := milisp.EvalFloat(env, e)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	// Output: 12
}
