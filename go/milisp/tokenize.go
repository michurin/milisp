package milisp

import (
	"fmt"
	"strconv"
)

const (
	sSpaces = iota
	sString
	sQuotedString
	sCharAfterSlash
	sComment
	sStop
	sError
)

const (
	cSpace = iota + 1000
	cBracketOpen
	cBracketClose
	cQuote
	cSlash
	cEOF
	cOther
	cCommentStart
	cNewLine
)

const (
	opNop      = 0
	opNewToken = 1 << iota
	opAppendChar
	opSaveToken
	opSaveQuotedToken
	opOpenToken
	opCloseToken
)

func tokenize(text string) ([]universalToken, error) {
	state := sSpaces
	line := 1
	pos := 0
	tokens := []universalToken(nil)
	chars := []rune(nil)
	var op int
	startLine := line
	startPos := pos
	for _, ch := range text + "\x1b" {
		if ch == 10 {
			line++
			pos = -1 // eliminate \n itself
		}
		pos++
		tp := charType(ch)
		state, op = stateTransitionFunction(state, tp)
		if state == sError {
			if tp == cEOF {
				return nil, fmt.Errorf("unexpected EOF")
			}
			return nil, fmt.Errorf("unexpected char %c at %d:%d", ch, line, pos)
		}
		if op&opNewToken > 0 {
			startLine = line
			startPos = pos
			chars = nil
		}
		if op&opAppendChar > 0 {
			chars = append(chars, ch)
		}
		if op&opSaveToken > 0 {
			s := string(chars)
			f, err := strconv.ParseFloat(s, 64)
			if err == nil {
				tokens = append(tokens, universalToken{
					tp:   tpNumber,
					num:  f,
					str:  s,
					line: startLine,
					pos:  startPos,
				})
			} else {
				tokens = append(tokens, universalToken{
					tp:   tpSymbol,
					str:  string(chars),
					line: startLine,
					pos:  startPos,
				})
			}
		}
		if op&opSaveQuotedToken > 0 {
			tokens = append(tokens, universalToken{
				tp:   tpString,
				str:  string(chars),
				line: startLine,
				pos:  startPos,
			})
		}
		if op&opOpenToken > 0 {
			tokens = append(tokens, universalToken{
				tp:   tpOpen,
				str:  "(",
				line: line,
				pos:  pos,
			})
		}
		if op&opCloseToken > 0 {
			tokens = append(tokens, universalToken{
				tp:   tpClose,
				str:  ")",
				line: line,
				pos:  pos,
			})
		}
		if state == sStop { // have to be at the end
			break
		}
	}
	return tokens, nil
}

func stateTransitionFunction(state int, symbol int) (int, int) {
	switch state {
	case sSpaces:
		switch symbol {
		case cSpace, cNewLine:
			return sSpaces, opNop
		case cCommentStart:
			return sComment, opNop
		case cBracketOpen:
			return sSpaces, opOpenToken
		case cBracketClose:
			return sSpaces, opCloseToken
		case cQuote:
			return sQuotedString, opNewToken
		case cSlash:
			return sError, opNop
		case cEOF:
			return sStop, opNop
		default:
			return sString, opNewToken | opAppendChar
		}
	case sString:
		switch symbol {
		case cSpace, cNewLine:
			return sSpaces, opSaveToken
		case cBracketOpen:
			return sSpaces, opSaveToken | opOpenToken
		case cBracketClose:
			return sSpaces, opSaveToken | opCloseToken
		case cQuote, cSlash, cCommentStart: // ["/#] couldn't be part of token
			return sError, opNop
		case cEOF:
			return sStop, opSaveToken
		default:
			return sString, opAppendChar
		}
	case sQuotedString:
		switch symbol {
		case cSlash:
			return sCharAfterSlash, opNop
		case cQuote:
			return sSpaces, opSaveQuotedToken
		case cEOF:
			return sError, opNop
		default: // cOther, cSpace, cBracketOpen/Close, sNewLine, cCommentStart
			return sQuotedString, opAppendChar
		}
	case sCharAfterSlash:
		switch symbol {
		case cEOF:
			return sError, opNop
		default:
			return sQuotedString, opAppendChar
		}
	default: // case sComment:
		switch symbol {
		case cNewLine:
			return sSpaces, opNop
		case cEOF:
			return sStop, opNop
		default:
			return sComment, opNop
		}
	}
}

func charType(ch rune) int {
	switch ch {
	case 10:
		return cNewLine
	case 9, 12, 13, 32:
		return cSpace
	case '(':
		return cBracketOpen
	case ')':
		return cBracketClose
	case '"':
		return cQuote
	case '\\':
		return cSlash
	case '#':
		return cCommentStart
	case 27:
		return cEOF
	}
	return cOther
}
