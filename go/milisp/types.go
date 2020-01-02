package milisp

// Environment define the scope of execution: variables and operations
type Environment map[string]interface{}

// Expression is only block of program. Everything is expression
type Expression interface {
	Eval(env Environment) (interface{}, error)
}

// Operation is a user defined thing that evaluate expression
type Operation interface {
	Perform(env Environment, args []Expression) (interface{}, error)
}

// OpFunc is a helper type to use function as Operation interface
type OpFunc func(env Environment, args []Expression) (interface{}, error)

// Perform operation function
func (f OpFunc) Perform(env Environment, args []Expression) (interface{}, error) {
	return f(env, args)
}
