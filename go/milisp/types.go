package milisp

type Env map[string]interface{}

type Expression interface {
	Eval(Env) (interface{}, error)
}

type Operation interface {
	Perform(Env, []Expression) (interface{}, error)
}

type OpFunc func(Env, []Expression) (interface{}, error)

func (f OpFunc) Perform(e Env, expr []Expression) (interface{}, error) {
	return f(e, expr)
}
