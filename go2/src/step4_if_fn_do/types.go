package main

import (
	"errors"
	"fmt"
)

// SExp : a S SExpression
type SExp interface {
	toString() string
	eval(Env) (SExp, error)
	copy() SExp
}

// Undefined : Undefined symbol. When an error occurred, reader returns UNDEF and err
type Undefined int

func (u Undefined) toString() string           { return "*Undefined*" }
func (u Undefined) eval(env Env) (SExp, error) { return u, nil }
func (u Undefined) copy() SExp                 { return u }

// UNDEF : Undef
const UNDEF = Undefined(0)

// NilType : the type of nil
type NilType int

func (n NilType) toString() string           { return "nil" }
func (n NilType) eval(env Env) (SExp, error) { return n, nil }
func (n NilType) copy() SExp                 { return n }

// NIL : Nil
const NIL = NilType(0)

// Bool : bool
type Bool bool

func (b Bool) toString() string           { return fmt.Sprint(b) }
func (b Bool) eval(env Env) (SExp, error) { return b, nil }
func (b Bool) copy() SExp                 { return b }

// Int : integer
type Int int

func (i Int) toString() string           { return fmt.Sprint(i) }
func (i Int) eval(env Env) (SExp, error) { return i, nil }
func (i Int) copy() SExp                 { return i }

// Symbol : Symbol
type Symbol string

func (s Symbol) toString() string { return string(s) }
func (s Symbol) eval(env Env) (SExp, error) {
	v, ok := env.get(s)
	if ok {
		return v, nil
	}
	return UNDEF, errors.New("can't find Symbol " + s.toString())
}
func (s Symbol) copy() SExp { return s }

// Keyword : Keyword
type Keyword string

func (k Keyword) toString() string           { return ":" + string(k) }
func (k Keyword) eval(env Env) (SExp, error) { return k, nil }
func (k Keyword) copy() SExp                 { return k }

// StringLiteral : should be print with '"'
type StringLiteral string

func (s StringLiteral) toString() string           { return fmt.Sprintf("\"%s\"", s) }
func (s StringLiteral) eval(env Env) (SExp, error) { return s, nil }
func (s StringLiteral) copy() SExp                 { return s }

// List : e.g. (1 2 3)
type List []SExp

func (l List) toString() string {
	return toStringSexpSlice("(", []SExp(l), ")")
}

func (l List) eval(env Env) (SExp, error) {
	if len(l) == 0 {
		return l, nil
	} else {
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

			return c.(CoreFunc).apply(args)
		default:
			println("error: can't apply")
		}
	}
	return UNDEF, errors.New("..........")
}

func (l List) copy() SExp {
	t := make(List, len(l))
	for i, v := range l {
		t[i] = v.copy()
	}
	return t
}

// Vector : e.g. [1 2 3]
type Vector []SExp

func (v Vector) toString() string {
	return toStringSexpSlice("[", []SExp(v), "]")
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
	t := make(List, len(v))
	for i, val := range v {
		t[i] = val.copy()
	}
	return t
}

func (v Vector) toList() List {
	return List(v)
}

// HashMap : {x 1, y 2}
type HashMap []SExp

func (hm HashMap) toString() string {
	return toStringSexpSlice("{", []SExp(hm), "}")
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

// CoreFunc : function + environment
type CoreFunc struct {
	fun func(args List) (SExp, error)
}

func (c CoreFunc) toString() string           { return "*CoreFunc*" }
func (c CoreFunc) eval(env Env) (SExp, error) { return c, nil }
func (c CoreFunc) copy() SExp                 { return c }

func (c CoreFunc) apply(args List) (SExp, error) { return c.fun(args) }

func toStringSexpSlice(ls string, sexps []SExp, rs string) string {
	t := make([]byte, 0, 10)
	t = append(t, ls...)
	for i, v := range sexps {
		t = append(t, v.toString()...)
		if i != len(sexps)-1 {
			t = append(t, " "...)
		}
	}
	t = append(t, rs...)
	return string(t)
}
