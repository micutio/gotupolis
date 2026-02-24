package server

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"

	opt "github.com/micutio/goptional"
	ts "github.com/micutio/gotupolis/pkg/tuplespace"
)

type tokenType uint

// Token types.
const (
	typeInt      tokenType = 1
	typeFloat              = 2
	typeStr                = 3
	typeTuple              = 4
	typeWildcard           = 5
	typeNone               = 6
)

type token struct {
	typ tokenType
	val interface{}
}

// Lexer is used to parse tuple definitions from an input string into the corresponding data
// structure. It's internal state is the input string, represented as a slice of runes and the
// current position of the parsing process.
type Lexer struct {
	pos int
	buf []rune
}

// NewLexer Initialises a new Lexer instance with a string that is supposed to contain N tuple
// definitions, with 0<=N.
func NewLexer(buffer string) Lexer {
	return Lexer{
		pos: 0,
		buf: []rune(buffer),
	}
}

// IntoTuples runs the lexer over the entire input and returns a slice of tuples.
// Returns an error if the input cannot be parsed to completion.
func (l *Lexer) IntoTuples() ([]ts.Tuple, error) {
	var tuples []ts.Tuple

	for t, err := l.nextTuple(); ; t, err = l.nextTuple() {
		if t.IsPresent() {
			tuples = append(tuples, t.Get())
		} else {
			return tuples, err
		}
	}
}

// nextTuple returns the next tuple contained in the string.
// If the string is empty or the Lexer has reached the end or no tuple can be read for some reason,
// then an empty option will be returned.
// Returns an error if the remaining input does not contain a fully formed tuple or any of its
// constituent elements cannot be parsed.
func (l *Lexer) nextTuple() (opt.Maybe[ts.Tuple], error) {
	if l.pos >= len(l.buf) {
		return opt.NewNothing[ts.Tuple](), nil
	}

	tkn, err := l.matchToken()
	if err != nil {
		errWithContext := errors.New(fmt.Sprintf("malformed tuple: {%s}", err))
		return opt.NewNothing[ts.Tuple](), errWithContext
	}

	elem := l.elemFromToken(tkn)
	if elem.GetType() == ts.ElemTuple {
		tupleVal := elem.GetValue().(ts.Tuple)
		return opt.NewJust(tupleVal), nil
	}

	return opt.NewJust(ts.MakeTuple()), nil
}

// matchToken returns the next token in the string.
// Returns an error if the token cannot be serialised into the tuple element data types
// (integer, floating point, string, tuple, wildcard).
func (l *Lexer) matchToken() (token, error) {
	r := l.buf[l.pos]
	switch r {
	case '-':
		return l.parseNumber()
	case '"':
		return l.parseString()
	case '_':
		return l.parseWildcard(), nil
	case '(':
		return l.parseTuple()
	case ' ':
		fallthrough
	case ',':
		l.pos += 1
		return l.matchToken()
	default:
		if unicode.IsDigit(r) {
			return l.parseNumber()
		}
		err := errors.New(fmt.Sprintf("invalid symbol '%v'", r))
		l.pos += 1
		return token{typeNone, nil}, err
	}
}

// elemFromToken converts a token into the corresponding tuple element.
func (l *Lexer) elemFromToken(tkn token) ts.Elem {
	switch tkn.typ {
	case typeInt:
		return ts.I(tkn.val.(int))
	case typeFloat:
		return ts.F(tkn.val.(float64))
	case typeStr:
		return ts.S(tkn.val.(string))
	case typeWildcard:
		return ts.Any()
	case typeTuple:
		tupleTokens := tkn.val.([]token)
		var tupleElems []ts.Elem
		for _, optTkn := range tupleTokens {
			tupleElems = append(tupleElems, l.elemFromToken(optTkn))
		}
		return ts.T(ts.MakeTuple(tupleElems...))
	default:
		return ts.None()
	}
}

// parseNumber attempts to parse an integer or floating point number from the input string.
// Returns an error if the parsed token is not a valid number, e.g.: multiple decimal points or
// contains any non-numerical symbols.
func (l *Lexer) parseNumber() (token, error) {
	start := l.pos
	isFloat := false

	for l.pos < len(l.buf) {
		char := l.buf[l.pos]
		switch char {
		case '.':
			if isFloat {
				return token{}, errors.New("float number with double decimal points")
			}

			isFloat = true
			l.pos += 1
			break
		case '-':
			if l.pos == start {
				l.pos += 1
			} else {
				// only allow a minus at the start of the number
				continue
			}
			break
		default:
			if unicode.IsDigit(char) {
				l.pos += 1
			} else {
				continue
			}
		}
	}

	if isFloat {
		return l.parseFloat(start)
	}

	return l.parseInt(start)
}

func (l *Lexer) parseFloat(start int) (token, error) {
	strVal := string(l.buf[start:l.pos])
	f, err := strconv.ParseFloat(strVal, 8)
	if err != nil {
		return token{}, err
	}

	return token{typeFloat, f}, nil
}

func (l *Lexer) parseInt(start int) (token, error) {
	strVal := string(l.buf[start:l.pos])
	i, err := strconv.Atoi(strVal)
	if err != nil {
		return token{}, err
	}
	return token{typeInt, i}, nil
}

// parseString attempts to parse a string from the input.
// Returns an error if it does not encounter a closing quote mark before reaching the end of the
// input.
func (l *Lexer) parseString() (token, error) {
	l.pos += 1
	start := l.pos
	for l.buf[l.pos] != '"' {
		l.pos += 1

		if l.pos >= len(l.buf) {
			return token{}, errors.New("parseString: incomplete string")
		}
	}

	l.pos += 1

	return token{typeStr, string(l.buf[start : l.pos-1])}, nil
}

// parseWildcard parses a single wildcard character (underscore).
func (l *Lexer) parseWildcard() token {
	start := l.pos
	l.pos += 1
	return token{typeWildcard, string(l.buf[start:l.pos])}
}

// parseTuple parses a complete tuple from the input.
// Returns an error if it does not encounter a closing parens before reaching the end of the input.
func (l *Lexer) parseTuple() (token, error) {
	l.pos += 1
	var tupleItems []token
	for l.buf[l.pos] != ')' {
		nextToken, tknErr := l.matchToken()
		if tknErr != nil {
			return token{}, tknErr
		}

		if nextToken.typ != typeNone {
			tupleItems = append(tupleItems, nextToken)
		}

		if l.pos >= len(l.buf) {
			return token{}, errors.New("error: incomplete tuple")
		}
	}
	l.pos += 1
	return token{typeTuple, tupleItems}, nil
}
