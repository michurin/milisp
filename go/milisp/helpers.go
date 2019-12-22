package milisp

import (
	"fmt"
	"strconv"
)

// EvalFloat is a shortcut for Exec + cast to float.
func EvalFloat(env Environment, e Expression) (float64, error) {
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
		}
		return 0, nil
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

// EvalString is a shortcut for Exec + cast to string.
func EvalString(env Environment, e Expression) (string, error) {
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

// EvalCode is a shortcut for Compile+Eval. Useful if you want to execute code just once.
func EvalCode(env Environment, text string) (interface{}, error) {
	p, err := Compile(text)
	if err != nil {
		return nil, err
	}
	res, err := p.Eval(env)
	if err != nil {
		return nil, err
	}
	return res, nil
}
