package milisp_test

import (
	"fmt"

	"github.com/michurin/milisp/go/milisp"
)

func opVector(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	result := make([]float64, len(args)) // you are free to use float32, that more native for many boosting tools
	for i, a := range args {
		r, err := milisp.EvalFloat(env, a)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}

func opAnd(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	for _, a := range args {
		r, err := a.Eval(env)
		if err != nil {
			return nil, err
		}
		if !r.(bool) {
			return false, nil
		}
	}
	return true, nil
}

func opIn(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	val, err := milisp.EvalString(env, args[0])
	if err != nil {
		return nil, err
	}
	rawList, err := args[1].Eval(env)
	if err != nil {
		return nil, err
	}
	list := rawList.([]string) //nolint:forcetypeassert // do not forget to check here
	for _, v := range list {
		if v == val {
			return true, nil
		}
	}
	return false, nil
}

func opSetStringList(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
	varName, err := milisp.EvalString(env, args[0])
	if err != nil {
		return nil, err
	}
	list := make([]string, len(args)-1)
	for i, a := range args[1:] { // check
		list[i], err = milisp.EvalString(env, a)
		if err != nil {
			return nil, err
		}
	}
	env[varName] = list
	return nil, nil
}

// Simplest example how to calculate ML vector from raw features
// using one-hot encode approach.
func Example_oneHotFeaturesWithConstants() {
	// simplest expressions for one-hot encode the categorical features
	// this code have to be attached to model as meta
	text := `
	(vector
	    (and (in phoneCountryCode UK) (in phoneAreaCode LDN))
	    (and (in phoneCountryCode IL) (in phoneAreaCode TLV))
	    (and (in phoneCountryCode RU) (in phoneAreaCode MSK))
    )`
	env := milisp.Environment{
		// functions
		"vector": milisp.OpFunc(opVector),
		"and":    milisp.OpFunc(opAnd),
		"in":     milisp.OpFunc(opIn),
		// constants (we use custom type without any limitations)
		"UK":  []string{"+44"},
		"IL":  []string{"+972"},
		"RU":  []string{"+7"},
		"LDN": []string{"020"},
		"TLV": []string{"3"},
		"MSK": []string{"095", "495"},
		// data (raw features values; we have to substitute corresponding input on every run prediction)
		"phoneCountryCode": "+44",
		"phoneAreaCode":    "020",
	}
	res, err := milisp.EvalCode(env, text)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	// Output: [1 0 0]
}

// It's not uncommon when you need convey some constants as part
// of training/predicting pipeline. You can use lisp code to
// initialize environment once. And then you are free to reuse
// this context multiply times.
func Example_oneHotFeaturesWithConstantsAsModelAttr() {
	// one-run program to initialize constants in env
	textConsts := `
	(prog
		(set_str_list "UK" "+44")
		(set_str_list "IL" "+972")
		(set_str_list "RU" "+7")
		(set_str_list "LDN" "020")
		(set_str_list "TLV" "3")
		(set_str_list "MSK" "095" "495")
	)
	`
	// simplest expressions for one-hot encode the categorical features
	text := `
	(vector
	    (and (in phoneCountryCode UK) (in phoneAreaCode LDN))
	    (and (in phoneCountryCode IL) (in phoneAreaCode TLV))
	    (and (in phoneCountryCode RU) (in phoneAreaCode MSK))
    )`
	env := milisp.Environment{
		// functions
		"vector":       milisp.OpFunc(opVector),
		"and":          milisp.OpFunc(opAnd),
		"in":           milisp.OpFunc(opIn),
		"prog":         milisp.OpFunc(evalAllReturnLastResult),
		"set_str_list": milisp.OpFunc(opSetStringList),
		// data
		"phoneCountryCode": "+44",
		"phoneAreaCode":    "020",
	}
	_, err := milisp.EvalCode(env, textConsts)
	if err != nil {
		panic(err)
	}
	res, err := milisp.EvalCode(env, text)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	// Output: [1 0 0]
}
