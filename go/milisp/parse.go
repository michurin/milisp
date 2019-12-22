package milisp

import "fmt"

func parse(tokens []universalToken, pos int) (Expression, int, bool, error) {
	if pos >= len(tokens) {
		return nil, 0, true, fmt.Errorf("unexpected end of file")
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
			e, pos, finish, err = parse(tokens, pos)
			if err != nil {
				return nil, 0, true, err
			}
			if finish {
				return expr{
					expr: ee,
					line: firstToken.line,
					pos:  firstToken.pos,
				}, pos + 1, false, nil
			}
			if pos >= len(tokens) {
				return nil, 0, true, fmt.Errorf("unexpected end of file; can not end expr started at %s", firstToken)
			}
			ee = append(ee, e)
		}
	case tpClose:
		return nil, pos, true, nil // pos will be shifted on call side
	default:
		return firstToken, pos + 1, false, nil
	}
}
