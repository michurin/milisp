package milisp

import "fmt"

const (
	tpSymbol = iota
	tpNumber
	tpString
	tpOpen
	tpClose
)

type universalToken struct {
	tp   int
	num  float64
	str  string
	line int
	pos  int
}

func (t universalToken) String() string {
	return fmt.Sprintf("%s:%s@%d:%d", []string{"SYM", "NUM", "STR", "BEG", "END"}[t.tp], t.str, t.line, t.pos)
}

func (t universalToken) Eval(env Environment) (interface{}, error) {
	switch t.tp {
	case tpSymbol:
		x, ok := env[t.str]
		if !ok {
			return nil, fmt.Errorf("runtime error: unknown symbol: %s", t)
		}
		return x, nil
	case tpNumber:
		return t.num, nil
	case tpString:
		return t.str, nil
	default: // case tpOpen, tpClose:
		return nil, fmt.Errorf("runtime error: impossible symbol: %s", t)
	}
}
