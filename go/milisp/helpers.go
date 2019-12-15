package milisp

import (
	"fmt"
	"strconv"
)

func EvalFloat(e Expression, env Env) (float64, error) {
	if e == nil {
		return 0, fmt.Errorf("nil interface")
	}
	r, err := e.Eval(env)
	if err != nil {
		return 0, err
	}
	switch v := r.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("can not convert to float %v while executing %s", v, e)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("can not cast %T to float64 (not implemented): %s", r, e)
	}
}

func EvalString(e Expression, env Env) (string, error) {
	if e == nil {
		return "", fmt.Errorf("nil interface")
	}
	r, err := e.Eval(env)
	if err != nil {
		return "", err
	}
	switch v := r.(type) {
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("can not cast %T to string (not implemented): %s", r, e)
	}
}

func EvalCode(env Env, text string) (interface{}, error) {
	p, err := Prog(text)
	if err != nil {
		return nil, err
	}
	res, err := p.Eval(env)
	if err != nil {
		return nil, err
	}
	return res, nil
}
