package milisp_test

import (
	"fmt"

	"github.com/michurin/milisp/go/milisp"
)

func evalAllReturnLastResult(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	res := interface{}(nil)
	err := error(nil)
	for _, e := range expr[1:] { // check len in real life
		res, err = e.Eval(env)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func mulAll(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	if len(expr) <= 1 { // command itself and at least one argument needed
		return nil, fmt.Errorf("too few args: %s", expr)
	}
	x := float64(1)
	for _, e := range expr[1:] {
		res, err := milisp.EvalFloat(e, env)
		if err != nil {
			return nil, err
		}
		x *= res
	}
	return x, nil
}

func sumAll(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	x := float64(0)
	for _, e := range expr[1:] {
		res, err := milisp.EvalFloat(e, env)
		if err != nil {
			return nil, err
		}
		x += res
	}
	return x, nil
}

func setVar(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	varName, err := milisp.EvalString(expr[1], env)
	if err != nil {
		return nil, err
	}
	varValue, err := expr[2].Eval(env)
	if err != nil {
		return nil, err
	}
	env[varName] = varValue
	return nil, nil
}

func loop(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	varName, err := milisp.EvalString(expr[1], env)
	if err != nil {
		return nil, err
	}
	first, err := milisp.EvalFloat(expr[2], env)
	if err != nil {
		return nil, err
	}
	last, err := milisp.EvalFloat(expr[3], env)
	if err != nil {
		return nil, err
	}
	body := expr[4]
	for i := int(first); i <= int(last); i++ {
		env[varName] = float64(i)
		_, err = body.Eval(env)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func ifGtOne(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	value, err := milisp.EvalFloat(expr[1], env)
	if err != nil {
		return nil, err
	}
	res := interface{}(nil)
	if value > 1. {
		res, err = expr[2].Eval(env)
	} else {
		res, err = expr[3].Eval(env)
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

func functionDefinition(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	funcName, err := milisp.EvalString(expr[1], env)
	if err != nil {
		return nil, err
	}
	argName, err := milisp.EvalString(expr[2], env)
	if err != nil {
		return nil, err
	}
	env[funcName] = function{
		argName: argName,
		body:    expr[3],
	}
	return nil, nil
}

func functionCall(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	funcName, err := milisp.EvalString(expr[1], env)
	if err != nil {
		return nil, err
	}
	argValue, err := expr[2].Eval(env)
	if err != nil {
		return nil, err
	}
	f := env[funcName].(function) // check, check, check...
	localEnv := map[string]interface{}{}
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

func ExampleFactorialLoop() {
	text := `
	(prog                     # execute all floowing expressions and return result of last
	    (set "x" 1)           # x = 1
	    (loop "i" 1 N         # for i = 1; i <= N; i++
	        (set "x" (* x i)) # x = x * i
	    )
	    x                     # return x
	)`
	env := map[string]interface{}{
		"prog": milisp.OpFunc(evalAllReturnLastResult),
		"set":  milisp.OpFunc(setVar),
		"loop": milisp.OpFunc(loop),
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

func ExampleFactorialRecursive() {
	text := `
    (prog
        (def "F" "x" (if_gt_one   # if x > 1 then F(x-1) else 1
            x
            (* x (call "F" (+ x -1)))
            1
        ))
        (call "F" N)
    )`
	env := map[string]interface{}{
		"prog":      milisp.OpFunc(evalAllReturnLastResult),
		"set":       milisp.OpFunc(setVar),
		"def":       milisp.OpFunc(functionDefinition),
		"call":      milisp.OpFunc(functionCall),
		"if_gt_one": milisp.OpFunc(ifGtOne),
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
