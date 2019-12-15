package milisp_test

import (
	"fmt"

	"github.com/michurin/milisp/go/milisp"
)

func opVector(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	result := make([]float64, len(expr)-1) // you are free to use float32, that more native for many boosting tools
	for i := 0; i < len(expr)-1; i++ {
		r, err := milisp.EvalFloat(expr[i+1], env)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}

func opAnd(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	for i := 0; i < len(expr)-1; i++ {
		r, err := expr[i+1].Eval(env)
		if err != nil {
			return nil, err
		}
		if !r.(bool) {
			return false, nil
		}
	}
	return true, nil
}

func opIn(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	val, err := milisp.EvalString(expr[1], env)
	if err != nil {
		return nil, err
	}
	rawList, err := expr[2].Eval(env)
	if err != nil {
		return nil, err
	}
	list := rawList.([]string) // check
	for _, v := range list {
		if v == val {
			return true, nil
		}
	}
	return false, nil
}

func opSetStringList(env milisp.Env, expr []milisp.Expression) (interface{}, error) {
	varName, err := milisp.EvalString(expr[1], env)
	if err != nil {
		return nil, err
	}
	list := make([]string, len(expr)-2)
	for i, e := range expr[2:] { // check
		list[i], err = milisp.EvalString(e, env)
		if err != nil {
			return nil, err
		}
	}
	env[varName] = list
	return nil, nil
}

func ExampleFeaturesUsingConstants() {
	text := `
	(vector
	    (and (in phoneCountryCode UK) (in phoneAreaCode LDN))
	    (and (in phoneCountryCode IL) (in phoneAreaCode TLV))
	    (and (in phoneCountryCode RU) (in phoneAreaCode MSK))
    )`
	env := map[string]interface{}{
		// functions
		"vector": milisp.OpFunc(opVector),
		"and":    milisp.OpFunc(opAnd),
		"in":     milisp.OpFunc(opIn),
		// const
		"UK":  []string{"+44"},
		"IL":  []string{"+972"},
		"RU":  []string{"+7"},
		"LDN": []string{"020"},
		"TLV": []string{"3"},
		"MSK": []string{"095", "495"},
		// data
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

// The same, but we set constants from program too. I could be attached to
// model as well as transformation operations. And here is example
// how to apply number of programs to one context.
func ExampleFeaturesWithoutConstants() {
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
	text := `
	(vector
	    (and (in phoneCountryCode UK) (in phoneAreaCode LDN))
	    (and (in phoneCountryCode IL) (in phoneAreaCode TLV))
	    (and (in phoneCountryCode RU) (in phoneAreaCode MSK))
    )`
	env := map[string]interface{}{
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
