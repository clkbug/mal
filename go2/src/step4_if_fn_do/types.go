package main

import (
	"errors"
	"fmt"
)

// SExp : a S SExpression
type SExp interface {
	toString() string
	eval(Env) (SExp, error)
}

// Undefined : Undefined symbol. When an error occurred, reader returns UNDEF and err
type Undefined int

func (u Undefined) toString() string           { return "*Undefined*" }
func (u Undefined) eval(env Env) (SExp, error) { return u, nil }

// UNDEF : Undef
const UNDEF = Undefined(0)

// NilType : the type of nil
type NilType int

func (n NilType) toString() string           { return "nil" }
func (n NilType) eval(env Env) (SExp, error) { return n, nil }

// NIL : Nil
const NIL = NilType(0)

// Bool : bool
type Bool bool

func (b Bool) toString() string           { return fmt.Sprint(b) }
func (b Bool) eval(env Env) (SExp, error) { return b, nil }

// Int : integer
type Int int

func (i Int) toString() string           { return fmt.Sprint(i) }
func (i Int) eval(env Env) (SExp, error) { return i, nil }

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

// Keyword : Keyword
type Keyword string

func (k Keyword) toString() string           { return ":" + string(k) }
func (k Keyword) eval(env Env) (SExp, error) { return k, nil }

// StringLiteral : should be print with '"'
type StringLiteral string

func (s StringLiteral) toString() string           { return fmt.Sprintf("\"%s\"", s) }
func (s StringLiteral) eval(env Env) (SExp, error) { return s, nil }

// List : e.g. (1 2 3)
type List []SExp

func (l List) toString() string {
	return toStringSexpSlice("(", []SExp(l), ")")
}

func (l List) eval(env Env) (SExp, error) {
	if len(l) == 0 {
		return l, nil
	} else if len(l) > 1 {
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
			default:
				panic("can't reach here... eval special form")
			}
		}
		switch c, err := l[0].eval(env); c.(type) {
		case Closure: // apply
			args := make(List, len(l)-1)
			for i, elem := range l[1:] {
				args[i], err = elem.eval(env)
				if err != nil {
					return UNDEF, err
				}
			}

			return c.(Closure).apply(args), nil
		default:
			println("error: can't apply")
		}
	}
	return UNDEF, errors.New("..........")
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

// Func : function
type Func func(env Env, args List) SExp

// Closure : function + environment
type Closure struct {
	env Env
	fun Func
}

func (c Closure) toString() string           { return "*Closure*" }
func (c Closure) eval(env Env) (SExp, error) { return c, nil }

func (c Closure) apply(args List) SExp { return c.fun(c.env, args) }

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
