package main

import (
	"errors"
	"fmt"
)

// SExp : a S SExpression
type SExp interface {
	toString() string
	printStr(isReadable bool) string
	eval(Env) (SExp, error)
	copy() SExp
	isSame(SExp) bool
}

// Undefined : Undefined symbol. When an error occurred, reader returns UNDEF and err
type Undefined int

func (u Undefined) toString() string           { return "*Undefined*" }
func (u Undefined) printStr(_ bool) string     { return u.toString() }
func (u Undefined) eval(env Env) (SExp, error) { return u, nil }
func (u Undefined) copy() SExp                 { return u }
func (u Undefined) isSame(s SExp) bool {
	switch s.(type) {
	case Undefined:
		return true
	}
	return false
}

// UNDEF : Undef
const UNDEF = Undefined(0)

// NilType : the type of nil
type NilType int

func (n NilType) toString() string           { return "nil" }
func (n NilType) printStr(_ bool) string     { return n.toString() }
func (n NilType) eval(env Env) (SExp, error) { return n, nil }
func (n NilType) copy() SExp                 { return n }
func (n NilType) isSame(s SExp) bool {
	switch s.(type) {
	case NilType:
		return true
	}
	return false
}

// NIL : Nil
const NIL = NilType(0)

// Bool : bool
type Bool bool

func (b Bool) toString() string           { return fmt.Sprint(b) }
func (b Bool) printStr(_ bool) string     { return b.toString() }
func (b Bool) eval(env Env) (SExp, error) { return b, nil }
func (b Bool) copy() SExp                 { return b }
func (b Bool) isSame(s SExp) bool {
	switch s := s.(type) {
	case Bool:
		return b == s
	}
	return false
}

// Int : integer
type Int int

func (i Int) toString() string           { return fmt.Sprint(i) }
func (i Int) printStr(_ bool) string     { return i.toString() }
func (i Int) eval(env Env) (SExp, error) { return i, nil }
func (i Int) copy() SExp                 { return i }
func (i Int) isSame(s SExp) bool {
	switch s := s.(type) {
	case Int:
		return i == s
	}
	return false
}

// Symbol : Symbol
type Symbol string

func (s Symbol) toString() string       { return string(s) }
func (s Symbol) printStr(_ bool) string { return s.toString() }
func (s Symbol) eval(env Env) (SExp, error) {
	v, ok := env.get(s)
	if ok {
		return v, nil
	}
	return UNDEF, errors.New("can't find Symbol " + s.toString())
}
func (s Symbol) copy() SExp { return s }
func (s Symbol) isSame(se SExp) bool {
	switch t := se.(type) {
	case Symbol:
		return s == t
	}
	return false
}

// Keyword : Keyword
type Keyword string

func (k Keyword) toString() string           { return ":" + string(k) }
func (k Keyword) printStr(_ bool) string     { return k.toString() }
func (k Keyword) eval(env Env) (SExp, error) { return k, nil }
func (k Keyword) copy() SExp                 { return k }
func (k Keyword) isSame(s SExp) bool {
	switch s := s.(type) {
	case Keyword:
		return s == k
	}
	return false
}

// StringLiteral : should be print with '"'
type StringLiteral string

func (s StringLiteral) toString() string {
	return "\"" + string(s) + "\""
}
func (s StringLiteral) printStr(isReadable bool) string {
	if isReadable {
		return fmt.Sprintf("\"%s\"", s.escape())
	}
	return string(s)
}
func (s StringLiteral) eval(env Env) (SExp, error) { return s, nil }
func (s StringLiteral) copy() SExp                 { return s }
func (s StringLiteral) isSame(t SExp) bool {
	switch t := t.(type) {
	case StringLiteral:
		return s == t
	}
	return false
}

func (s StringLiteral) escape() StringLiteral {
	str := string(s)
	ret := ""
	for _, r := range str {
		switch r {
		case '\n':
			ret += "\\n"
		case '"':
			ret += "\\\""
		case '\\':
			ret += "\\\\"
		default:
			ret += fmt.Sprintf("%c", r)
		}
	}
	return StringLiteral(ret)
}
func (s StringLiteral) unescape() StringLiteral {
	str := string(s)
	ret := ""
	bs := false
	for _, r := range str {
		if bs {
			bs = false
			switch r {
			case 'n':
				ret += "\n"
			case '"':
				ret += "\""
			case '\\':
				ret += "\\"
			default:
				panic("unescape!!!")
			}
		} else if r == '\\' {
			bs = true
		} else {
			ret += fmt.Sprintf("%c", r)
		}
	}
	return StringLiteral(ret)
}

// List : e.g. (1 2 3)
type List []SExp

func (l List) toString() string {
	return toStringSexpSlice("(", []SExp(l), ")", true)
}
func (l List) printStr(isReadable bool) string {
	return toStringSexpSlice("(", []SExp(l), ")", isReadable)
}

func (l List) eval(env Env) (SExp, error) {
	if len(l) == 0 {
		return l, nil
	}
	if v, ok := isSpecialForm(l[0]); ok {
		switch v {
		case IF:
			return evalIf(env, l[1:])
		case COND:
		case OR:
		case DEF:
			return evalDef(env, l[1:])
		case DEFMACRO:
		case LET:
			return evalLet(env, l[1:])
		case FN:
			return evalFn(env, l[1:])
		default:
			panic("can't reach here... eval special form")
		}
	}
	switch c, err := l[0].eval(env); c.(type) {
	case CoreFunc: // apply
		args := make(List, len(l)-1)
		for i, elem := range l[1:] {
			args[i], err = elem.eval(env)
			if err != nil {
				return UNDEF, err
			}
		}
		return c.(CoreFunc).apply(args, env)
	case Closure:
		args := make(List, len(l)-1)
		for i, elem := range l[1:] {
			args[i], err = elem.eval(env)
			if err != nil {
				return UNDEF, err
			}
		}
		return c.(Closure).apply(args)
	default:
		println("error: can't apply\n\t" + l.toString())
	}

	return UNDEF, errors.New("eval?")
}

func (l List) copy() SExp {
	return l
}

func (l List) isSame(s SExp) bool {
	switch s := s.(type) {
	case List:
		if len(l) != len(s) {
			return false
		}
		for i := 0; i < len(l); i++ {
			if !l[i].isSame(s[i]) {
				return false
			}
		}
		return true
	case Vector:
		return s.toList().isSame(l)
	}
	return false
}

// Vector : e.g. [1 2 3]
type Vector []SExp

func (v Vector) toString() string {
	return toStringSexpSlice("[", []SExp(v), "]", true)
}
func (v Vector) printStr(isReadable bool) string {
	return toStringSexpSlice("[", []SExp(v), "]", isReadable)
}

func (v Vector) eval(env Env) (SExp, error) {
	ret := make(Vector, len(v))
	for i, elem := range v {
		var err error
		ret[i], err = elem.eval(env)
		if err != nil {
			return UNDEF, err
		}
	}
	return ret, nil
}

func (v Vector) copy() SExp {
	return v
}

func (v Vector) isSame(s SExp) bool {
	switch s := s.(type) {
	case List:
		return v.toList().isSame(s)
	case Vector:
		if len(v) != len(s) {
			return false
		}
		for i := 0; i < len(v); i++ {
			if !v[i].isSame(s[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (v Vector) toList() List {
	return List(v)
}

// HashMap : {x 1, y 2}ya
type HashMap []SExp

func (hm HashMap) toString() string {
	return toStringSexpSlice("{", []SExp(hm), "}", true)
}
func (hm HashMap) printStr(isReadable bool) string {
	return toStringSexpSlice("{", []SExp(hm), "}", isReadable)
}

func (hm HashMap) eval(env Env) (SExp, error) {
	ret := make(HashMap, len(hm))
	for i, elem := range hm {
		var err error
		ret[i], err = elem.eval(env)
		if err != nil {
			return UNDEF, err
		}
	}
	return ret, nil
}

func (hm HashMap) copy() SExp {
	t := make(HashMap, len(hm))
	for key, val := range hm {
		t[key] = val.copy()
	}
	return t
}

func (hm HashMap) isSame(s SExp) bool {
	switch s := s.(type) {
	case HashMap:
		if len(hm) != len(s) {
			return false
		}
		for i := 0; i < len(hm); i++ {
			if !hm[i].isSame(s[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// CoreFunc : function
type CoreFunc func(args List, env Env) (SExp, error)

func (c CoreFunc) toString() string           { return "*CoreFunc*" }
func (c CoreFunc) printStr(_ bool) string     { return "*CoreFunc*" }
func (c CoreFunc) eval(env Env) (SExp, error) { return c, nil }
func (c CoreFunc) copy() SExp                 { return c }
func (c CoreFunc) isSame(s SExp) bool         { return false } // Function isn't comparable

func (c CoreFunc) apply(args List, env Env) (SExp, error) { return c(args, env) }

// Closure : environment + arg List + body
type Closure struct {
	name   Symbol // recursive function use
	env    Env
	params []Symbol
	body   SExp
}

func (c Closure) toString() string           { return "*Closure*" }
func (c Closure) printStr(_ bool) string     { return "*Closure*" }
func (c Closure) eval(env Env) (SExp, error) { return c, nil }
func (c Closure) isSame(s SExp) bool         { return false } // Closure isn't comparable
func (c Closure) copy() SExp {
	return c
}
func (c Closure) apply(args List) (SExp, error) {
	ne := makeNewEnv(c.env)
	for i, p := range c.params {
		if string(p) == "&" {
			ne.set(c.params[i+1], args[i:])
			break
		} else {
			ne.set(p, args[i])
		}
	}
	return c.body.eval(ne)
}

func toStringSexpSlice(ls string, sexps []SExp, rs string, isReadable bool) string {
	t := make([]byte, 0, 10)
	t = append(t, ls...)
	for i, v := range sexps {
		t = append(t, v.printStr(isReadable)...)
		if i != len(sexps)-1 {
			t = append(t, " "...)
		}
	}
	t = append(t, rs...)
	return string(t)
}

// Atom : atom
type Atom struct {
	ref SExp
}

func (a Atom) toString() string {
	return "(atom " + a.ref.toString() + ")"
}
func (a Atom) printStr(b bool) string {
	return "(atom " + a.ref.printStr(b) + ")"
}
func (a Atom) eval(env Env) (SExp, error) { return a, nil }
func (a Atom) isSame(s SExp) bool {
	switch s := s.(type) {
	case Atom:
		if a == s {
			return true
		}
	}
	return false
}
func (a Atom) copy() SExp { return a }
