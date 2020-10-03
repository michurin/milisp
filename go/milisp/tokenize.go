package milisp

import (
	"fmt"
	"strconv"
)

const opNop = 0

// States, char classes and operations for tokenization FSM

const (
	sSpaces = iota
	sString
	sQuotedString
	sCharAfterSlash
	sComment
	sStop
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
	opNewToken = 1 << iota
	opAppendChar
	opSaveToken
	opSaveQuotedToken
	opOpenToken
	opCloseToken
	opStopOk
	opErrorChar
	opErrorEOF
)

// States, char classes and operations for positioning FSM

const (
	sCR = iota
	sNL
	sLine
)

const (
	cCR = iota
	cNL
	cTab
	cChar // including space
)

const (
	opNewLine = iota + 1
	opStepOne
	opStepTab
)

func tokenize(text string) ([]universalToken, error) {
	tokenState := sSpaces
	lineState := sLine
	line := 1
	pos := 1
	tokens := []universalToken(nil)
	chars := []rune(nil)
	var op int
	startLine := line
	startPos := pos
	for _, ch := range text + "\x1b" {
		// classify char
		tp, spTp := charType(ch)
		// tokenization
		tokenState, op = tokenizeStateTransitionFunction(tokenState, tp)
		if op&opErrorChar > 0 {
			return nil, fmt.Errorf("unexpected char %c at %d:%d", ch, line, pos)
		}
		if op&opErrorEOF > 0 {
			return nil, fmt.Errorf("unexpected EOF")
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
		if op&opStopOk > 0 { // have to be tha last operation
			break
		}
		// find out position of next char
		lineState, op = charPositionStateTransitionFunction(lineState, spTp)
		switch op {
		case opNewLine:
			line++
			pos = 1
		case opStepOne:
			pos++
		case opStepTab:
			pos += 8 - (pos-1)%8
		}
	}
	return tokens, nil
}

func tokenizeStateTransitionFunction(state int, symbol int) (int, int) {
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
			return sStop, opErrorChar
		case cEOF:
			return sStop, opStopOk
		case cOther:
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
			return sStop, opErrorChar
		case cEOF:
			return sStop, opSaveToken | opStopOk
		case cOther:
			return sString, opAppendChar
		}
	case sQuotedString:
		switch symbol {
		case cSlash:
			return sCharAfterSlash, opNop
		case cQuote:
			return sSpaces, opSaveQuotedToken
		case cEOF:
			return sStop, opErrorEOF
		case cOther, cSpace, cBracketOpen, cBracketClose, cNewLine, cCommentStart:
			return sQuotedString, opAppendChar
		}
	case sCharAfterSlash:
		switch symbol {
		case cEOF:
			return sStop, opErrorEOF
		case cSlash, cQuote, cOther, cSpace, cBracketOpen, cBracketClose, cNewLine, cCommentStart:
			return sQuotedString, opAppendChar
		}
	case sComment:
		switch symbol {
		case cNewLine:
			return sSpaces, opNop
		case cEOF:
			return sStop, opStopOk
		case cSlash, cQuote, cOther, cSpace, cBracketOpen, cBracketClose, cCommentStart:
			return sComment, opNop
		}
	}
	panic("impossible state")
}

func charPositionStateTransitionFunction(state int, symbol int) (int, int) {
	// According https://www.unicode.org/reports/tr14/tr14-32.html
	// Long story short: valid newline combinations are \n, \r and \r\n:
	// \n\n\n\n — 4 lines
	// \r\r\r\r — 4 lines
	// \r\n\r\n — 2 lines
	// \n\r\n\r — 3 lines (\n + \r\n + \r)
	// The special case is for \r only
	// All \n, \xb, \xc etc. are interpreted in the same way.
	switch state {
	case sCR:
		switch symbol {
		case cCR:
			return sCR, opNewLine
		case cNL:
			return sNL, opNop
		case cTab:
			return sLine, opStepTab
		case cChar:
			return sLine, opStepOne
		}
	case sNL:
		switch symbol {
		case cCR:
			return sCR, opNewLine
		case cNL:
			return sNL, opNewLine
		case cTab:
			return sLine, opStepTab
		case cChar:
			return sLine, opStepOne
		}
	case sLine:
		switch symbol {
		case cCR:
			return sCR, opNewLine
		case cNL:
			return sNL, opNewLine
		case cTab:
			return sLine, opStepTab
		case cChar:
			return sLine, opStepOne
		}
	}
	panic("impossible state")
}

func charType(ch rune) (int, int) {
	switch ch {
	case '\r':
		return cNewLine, cCR
	case '\n', 0xb, 0xc, 0x2028, 0x2029:
		return cNewLine, cNL
	case 0x9:
		return cSpace, cTab
	case 0x20:
		return cSpace, cChar
	case '(':
		return cBracketOpen, cChar
	case ')':
		return cBracketClose, cChar
	case '"':
		return cQuote, cChar
	case '\\':
		return cSlash, cChar
	case '#':
		return cCommentStart, cChar
	case 0x1b:
		return cEOF, cChar
	}
	return cOther, cChar
}
