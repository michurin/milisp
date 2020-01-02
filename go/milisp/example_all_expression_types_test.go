package milisp_test

import (
	"fmt"

	"github.com/michurin/milisp/go/milisp"
)

func Example_allExpressionTypes() {
	// simple operation: sum all arguments
	opSumAll := milisp.OpFunc(func(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
		x := float64(0)
		for _, a := range args {
			res, err := milisp.EvalFloat(env, a)
			if err != nil {
				return nil, err
			}
			x += res
		}
		return x, nil
	})
	// nontrivial operation: returns other operation by name
	opGetOperationByName := milisp.OpFunc(func(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
		opName, err := milisp.EvalString(env, args[0])
		if err != nil {
			return nil, err
		}
		op, ok := env[opName]
		if !ok {
			return nil, fmt.Errorf("operation with name %s not exists in env", opName)
		}
		return op, nil
	})
	for _, params := range []struct {
		env  milisp.Environment
		text string
	}{
		// constants evaluate without context
		{nil, "1"},
		{nil, `"ok"`},
		{nil, "()"},
		// you are free to store in context any types
		{milisp.Environment{"x": []rune("ok")}, "x"},
		// simples operation
		{milisp.Environment{"+": opSumAll}, "(+ 1 2)"},
		// nontrivial expression
		{milisp.Environment{
			"+":  opSumAll,
			"op": opGetOperationByName,
			"x":  float64(1),
		}, `((op "+") x (+ 1 1))`},
	} {
		res, err := milisp.EvalCode(params.env, params.text)
		if err != nil {
			panic(err)
		}
		fmt.Println(params.text, "->", res)
	}
	// Output:
	// 1 -> 1
	// "ok" -> ok
	// () -> <nil>
	// x -> [111 107]
	// (+ 1 2) -> 3
	// ((op "+") x (+ 1 1)) -> 3
}
