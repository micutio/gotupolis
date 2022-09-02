package gotupolis

import (
	"fmt"
	"strconv"
	"unicode"

	opt "github.com/micutio/goptional"
	ts "github.com/micutio/gotupolis/pkg/tuplespace"
)

type tokenType uint

// Token types
const (
	T_INT      tokenType = 1
	T_FLOAT              = 2
	T_STRING             = 3
	T_TUPLE              = 4
	T_WILDCARD           = 5
)

type token struct {
	typ tokenType
	val string
}

type Lexer struct {
	pos int
	buf []rune
}

// Initialise a new Lexer instance with a string that is supposed to contain N tuple definitions,
// with 0<=N.
func NewLexer(buffer string) Lexer {
	return Lexer{
		pos: 0,
		buf: []rune(buffer),
	}
}

// nextTuple returns the next tuple contained in the string.
// If the string is empty or the Lexer has reached the end or no tuple can be read for some reason,
// then an empty option will be returned.
func (l *Lexer) nextTuple() opt.Maybe[ts.Tuple] {
	if l.pos >= len(l.buf) {
		return opt.NewNothing[ts.Tuple]()
	}

	tkn := l.matchToken()
	elem := l.elemFromToken(tkn)
	if elem.GetType() == ts.TUPLE {
		tupleVal := elem.GetValue().(ts.Tuple)
		return opt.NewJust(ts.Tuple(tupleVal))
	} else {
		return opt.NewJust(ts.MakeTuple())
	}
}

func (l *Lexer) matchToken() opt.Maybe[token] {
	switch r := l.buf[l.pos]; {
	case unicode.IsDigit(r) || r == '-':
		return l.parseNumber()
	case r == '"':
		return l.parseString()
	case r == '_':
		return l.parseWildcard()
	case r == '(':
		return l.parseTuple()
	case r == ',':
		l.pos += 1
		return opt.NewNothing[token]()
	default:
		fmt.Printf("invalid symbol '%v'", r)
		l.pos += 1
		return opt.NewNothing[token]()
	}
}

func (l *Lexer) elemFromToken(t opt.Maybe[token]) ts.Elem {
	if t.IsPresent() {
		switch t.Get().typ {
		case T_INT:
			i, err := strconv.Atoi(t.Get().val)
			if err != nil {
				panic(err)
			} else {
				return ts.I(i)
			}
		case T_FLOAT:
			f, err := strconv.ParseFloat(t.Get().val, 8)
			if err != nil {
				panic(err)
			} else {
				return ts.F(f)
			}
		case T_STRING:
			return ts.S(t.Get().val)
		case T_WILDCARD:
			return ts.Any()
		}
	}
	return ts.None()
}

func (l *Lexer) parseNumber() opt.Maybe[token] {
	start := l.pos
	isFloat := false
	for l.pos < len(l.buf) {
		switch r := l.buf[l.pos]; {
		case unicode.IsDigit(r):
			l.pos += 1
		case r == '.':
			if isFloat {
				panic("float number with double decimal points")
			} else {
				isFloat = true
				l.pos += 1
			}
		default:
			break
		}
	}

	var typ tokenType
	if isFloat {
		typ = T_FLOAT
	} else {
		typ = T_INT
	}
	return opt.NewJust(token{typ, string(l.buf[start:l.pos])})
}

func (l *Lexer) parseString() opt.Maybe[token] {
	l.pos += 1
	start := l.pos
	for l.buf[l.pos] != '"' {
		l.pos += 1

		if l.pos >= len(l.buf) {
			fmt.Println("error: incomplete string!")
			return opt.NewNothing[token]()
		}
	}

	l.pos += 1

	return opt.NewJust(token{T_STRING, string(l.buf[start:l.pos])})
}

func (l *Lexer) parseWildcard() opt.Maybe[token] {
	start := l.pos
	l.pos += 1
	return opt.NewJust(token{T_WILDCARD, string(l.buf[start:l.pos])})
}

func (l *Lexer) parseTuple() opt.Maybe[token] {
	// TODO
}
