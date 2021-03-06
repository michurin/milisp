package milisp

import "fmt"

type expr struct {
	expr []Expression
	line int
	pos  int
}

func (e expr) String() string {
	return fmt.Sprintf("%s@%d:%d", e.expr, e.line, e.pos)
}

func (e expr) Eval(env Environment) (interface{}, error) {
	if len(e.expr) == 0 {
		return nil, nil
	}
	op, err := e.expr[0].Eval(env)
	if err != nil {
		return nil, err
	}
	operation, ok := op.(Operation)
	if !ok {
		return nil, fmt.Errorf("operation %T not executable: %s", op, e.expr[0])
	}
	res, err := operation.Perform(env, e.expr[1:])
	if err != nil {
		return nil, err
	}
	return res, err
}
