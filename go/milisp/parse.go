package milisp

import "fmt"

func parse(tokens []universalToken, pos int) (Expression, int, error, bool) {
	if pos >= len(tokens) {
		return nil, 0, fmt.Errorf("unexpected end of file"), true
	}
	firstToken := tokens[pos]
	switch firstToken.tp {
	case tpOpen:
		ee := []Expression(nil)
		var e Expression
		var err error
		var finish bool
		pos++
		for {
			e, pos, err, finish = parse(tokens, pos)
			if err != nil {
				return nil, 0, err, true
			}
			if finish {
				return expr{
					expr: ee,
					line: firstToken.line,
					pos:  firstToken.pos,
				}, pos + 1, nil, false
			}
			if pos >= len(tokens) {
				return nil, 0, fmt.Errorf("unexpected end of file; can not end expr started at %s", firstToken), true
			}
			ee = append(ee, e)
		}
	case tpClose:
		return nil, pos, nil, true // pos will be shifted on call side
	default:
		return firstToken, pos + 1, nil, false
	}
}
