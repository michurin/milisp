package milisp_test

import (
	"fmt"

	"github.com/michurin/milisp/go/milisp"
)

func evalAllReturnLastResult(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	res := interface{}(nil)
	for _, a := range args { // check len in real life
		var err error
		res, err = a.Eval(env)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func mulAll(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("too few args: %s", args)
	}
	x := float64(1)
	for _, a := range args {
		res, err := milisp.EvalFloat(env, a)
		if err != nil {
			return nil, err
		}
		x *= res
	}
	return x, nil
}

func sumAll(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	x := float64(0)
	for _, a := range args {
		res, err := milisp.EvalFloat(env, a)
		if err != nil {
			return nil, err
		}
		x += res
	}
	return x, nil
}

func setVar(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	varName, err := milisp.EvalString(env, args[0])
	if err != nil {
		return nil, err
	}
	varValue, err := args[1].Eval(env)
	if err != nil {
		return nil, err
	}
	env[varName] = varValue
	return nil, nil
}

func loop(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	varName, err := milisp.EvalString(env, args[0])
	if err != nil {
		return nil, err
	}
	first, err := milisp.EvalFloat(env, args[1])
	if err != nil {
		return nil, err
	}
	last, err := milisp.EvalFloat(env, args[2])
	if err != nil {
		return nil, err
	}
	body := args[3]
	for i := int(first); i <= int(last); i++ {
		env[varName] = float64(i)
		_, err = body.Eval(env)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func ifGtOne(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	value, err := milisp.EvalFloat(env, args[0])
	if err != nil {
		return nil, err
	}
	var res interface{}
	if value > 1. { // example of laziness, we don't evaluate unnecessary argument
		res, err = args[1].Eval(env)
	} else {
		res, err = args[2].Eval(env)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

type function struct {
	argName string
	body    milisp.Expression
}

func functionDefinition(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	funcName, err := milisp.EvalString(env, args[0])
	if err != nil {
		return nil, err
	}
	argName, err := milisp.EvalString(env, args[1])
	if err != nil {
		return nil, err
	}
	env[funcName] = function{
		argName: argName,
		body:    args[2],
	}
	return nil, nil
}

func functionCall(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	funcName, err := milisp.EvalString(env, args[0])
	if err != nil {
		return nil, err
	}
	argValue, err := args[1].Eval(env)
	if err != nil {
		return nil, err
	}
	f := env[funcName].(function) //nolint:forcetypeassert // do not forget to check here
	localEnv := milisp.Environment{}
	for k, v := range env {
		localEnv[k] = v
	}
	localEnv[f.argName] = argValue
	res, err := f.body.Eval(localEnv)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func Example_factorialLoop() {
	text := `
	(prog                     # execute all following expressions and return result of last
	    (set "x" 1)           # x = 1
	    (loop "i" 1 N         # for i = 1; i <= N; i++
	        (set "x" (* x i)) # x = x * i
	    )
	    x                     # return x
	)`
	env := milisp.Environment{
		// take a look inside examples file for implementations
		"prog": milisp.OpFunc(evalAllReturnLastResult),
		"set":  milisp.OpFunc(setVar), // it shows how to create new variables in env
		"loop": milisp.OpFunc(loop),   // it shows how to mutate variables
		"*":    milisp.OpFunc(mulAll),
		"N":    5.,
	}
	res, err := milisp.EvalCode(env, text)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	// Output: 120
}

func Example_factorialRecursive() {
	text := `
	(prog
	    (def "F" "x" (if_gt_one   # if x > 1 then F(x-1) else 1
	        x
	        (* x (call "F" (+ x -1)))
	        1
	    ))
	    (call "F" N)
	)`
	env := milisp.Environment{
		// take a look inside examples file for implementations
		"prog":      milisp.OpFunc(evalAllReturnLastResult),
		"set":       milisp.OpFunc(setVar),
		"def":       milisp.OpFunc(functionDefinition),
		"call":      milisp.OpFunc(functionCall), // local scopes (local copies of env)
		"if_gt_one": milisp.OpFunc(ifGtOne),      // lazy and conditional execution
		"*":         milisp.OpFunc(mulAll),
		"+":         milisp.OpFunc(sumAll),
		"N":         5.,
	}
	res, err := milisp.EvalCode(env, text)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	// Output: 120
}
