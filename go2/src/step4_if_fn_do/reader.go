package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Reader is a reader
type Reader struct {
	s            []rune
	pos          int
	isReachedEND bool
}

// Token is the type of tokens
type Token string

const (
	// QUOTE 'x => x
	QUOTE = "'"
	// QUASIQUOTE `x  => (quasiquote x)
	QUASIQUOTE = "`"
	// UNQUOTE ~x => (unquote x)
	UNQUOTE = "~"
	// SPLICEUNQUOTE ~@x => (splice-unquote x)
	SPLICEUNQUOTE = "~@"
	// DEREF @x => (deref x)
	DEREF = "@"
	// META ^{a 1} [1 2 3] => (with-meta [1 2 3] {"a" 1})
	META = "^"
)

func initReader(s string) *Reader {
	return &Reader{
		s:            []rune(s),
		pos:          0,
		isReachedEND: false,
	}
}

// Next returns the token at the current position and increments the position.
func (r *Reader) next() (Token, error) {
	t, err := r.peek()

	if err != nil {
		return t, err
	}

	if r.isReachedEND {
		r.pos = len(r.s)
		return t, nil
	}

	for isSpace(r.s[r.pos]) {
		r.pos++
	}

	r.pos += len(t)

	return t, nil
}

// peek returns the toekn at the current position.
func (r *Reader) peek() (Token, error) {
	if r.isReachedEND {
		return "", nil
	} else if r.pos == len(r.s) {
		r.isReachedEND = true
		return "", nil
	}
	start := r.pos

	for isSpace(r.s[start]) {
		start++
		if start == len(r.s) {
			r.isReachedEND = true
			return "", nil
		}
	}

	switch r.s[start] {
	case '(', ')', '[', ']', '{', '}', '\'', '`', '@', '^':
		return runeToToken(r.s[start]), nil
	case ';':
		r.isReachedEND = true
		return "", nil
	case '~':
		if r.s[start+1] == '@' {
			return "~@", nil
		}
		return "~", nil

	case '"':
		end := start + 1
		for r.s[end] != '"' {
			if r.s[end] == '\\' {
				end++
			}
			end++

			if end >= len(r.s) {
				return "", errors.New("expected '\"', got EOF")
			}
		}
		return Token(r.s[start : end+1]), nil
	}

	end := start
	for !(isSpecial(r.s[end]) || isSpace(r.s[end])) {
		end++
		if end == len(r.s) {
			break
		}
	}
	return Token(r.s[start:end]), nil

}

func runeToString(c rune) string {
	var t [1]rune
	t[0] = c
	return string(t[:])
}

func runeToToken(c rune) Token {
	return Token(runeToString(c))
}

func isSpace(c rune) bool {
	switch c {
	case ' ', '\t', '\n', '\r', ',':
		return true
	default:
		return false
	}
}

func isSpecial(c rune) bool {
	return strings.ContainsAny(runeToString(c), "()[]{};\"'`@^")
}

func (r *Reader) readForm() (SExp, error) {
	t, err := r.peek()
	if err != nil {
		return UNDEF, err
	}
	switch t {
	case "(":
		return r.readSeq(")")
	case "[":
		return r.readSeq("]")
	case "{":
		return r.readSeq("}")
	case QUOTE, QUASIQUOTE, UNQUOTE, SPLICEUNQUOTE, DEREF:
		_, _ = r.next()
		s, e := r.readForm()
		if e != nil {
			return nil, e
		}
		qd := make([]SExp, 2)
		switch t {
		case QUOTE:
			qd[0] = Symbol("quote")
		case QUASIQUOTE:
			qd[0] = Symbol("quasiquote")
		case UNQUOTE:
			qd[0] = Symbol("unquote")
		case SPLICEUNQUOTE:
			qd[0] = Symbol("splice-unquote")
		case DEREF:
			qd[0] = Symbol("deref")
		}
		qd[1] = s
		return List(qd), nil
	case META:
		_, _ = r.next()
		s, e := r.readSeq("}")
		if e != nil {
			return nil, e
		}
		qd := make([]SExp, 3)
		qd[0] = Symbol("with-meta")
		qd[2] = s
		s, e = r.readForm()
		if e != nil {
			return nil, e
		}
		qd[1] = s
		return List(qd), nil
	case "":
		return UNDEF, nil
	default:
		return r.readAtom()
	}
}

func (r *Reader) readSeq(right string) (SExp, error) {
	r.next()
	l := make([]SExp, 0)
	for {
		t, err := r.peek()
		if err != nil {
			if strings.HasPrefix(err.Error(), "expected '\"'") {
				return UNDEF, fmt.Errorf("expected '%s', got EOF", right)
			}
		}
		if t == Token(right) {
			r.next()
			break
		} else if t == "" {
			return nil, fmt.Errorf("expected '%s', got EOF", right)
		}
		h, err := r.readForm()
		if err != nil {
			return UNDEF, err
		}
		l = append(l, h)
	}
	switch right {
	case ")":
		return List(l), nil
	case "]":
		return Vector(l), nil
	case "}":
		return HashMap(l), nil
	default:
		return UNDEF, errors.New("Invalid 'right' in reader.readSeq")
	}
}

func (r *Reader) readAtom() (SExp, error) {
	t, err := r.next()
	if err != nil {
		return UNDEF, err
	}
	if tmp := []rune(string(t)); strings.ContainsAny(runeToString(tmp[0]), "-0123456789") {
		i, e := strconv.Atoi(string(t))
		if e != nil {
			return Symbol(t), nil // when ParseInt fails, this works
		}
		return Int(i), nil
	} else if tmp[0] == '"' {
		return StringLiteral(string(tmp[1 : len(tmp)-1])).unescape(), nil
	} else if tmp[0] == ':' {
		return Keyword(tmp[1:]), nil
	} else if t == "true" {
		return Bool(true), nil
	} else if t == "false" {
		return Bool(false), nil
	} else if t == "nil" {
		return NIL, nil
	}
	return Symbol(t), nil
}
