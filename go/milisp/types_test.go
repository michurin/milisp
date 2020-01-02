package milisp_test

import (
	"errors"
	"testing"

	"github.com/michurin/milisp/go/milisp"
)

func TestEnvironment_symbolErrors(t *testing.T) {
	prErr := func(_ interface{}, err error) {
		if err == nil {
			t.Fatal("Have to be error")
		}
	}
	emptyEnv := milisp.Environment{}
	prErr(milisp.EvalCode(emptyEnv, "(X)"))
	prErr(milisp.EvalCode(milisp.Environment{"X": nil}, "(X)"))
	prErr(milisp.EvalCode(milisp.Environment{ // symbol exists, however evaluate with internal error
		"X": milisp.OpFunc(func(_ milisp.Environment, _ []milisp.Expression) (interface{}, error) {
			return nil, errors.New("error message")
		})}, "(X)"))
}
