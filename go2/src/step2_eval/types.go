package main

import "fmt"

// Undefined : Undefined symbol. When an error occurred, reader returns UNDEF and err
type Undefined int

func (u Undefined) toString() string  { return "*Undefined*" }
func (u Undefined) eval(env Env) SExp { return u }

// UNDEF : Undef
const UNDEF = Undefined(0)

// Int : integer
type Int int

func (i Int) toString() string  { return fmt.Sprint(i) }
func (i Int) eval(env Env) SExp { return i }

// Symbol : Symbol
type Symbol string

func (s Symbol) toString() string { return string(s) }
func (s Symbol) eval(env Env) SExp {
	v, ok := env.get(s)
	if ok {
		return v
	}
	return UNDEF
}

// Keyword : Keyword
type Keyword string

func (k Keyword) toString() string  { return ":" + string(k) }
func (k Keyword) eval(env Env) SExp { return k }

// StringLiteral : should be print with '"'
type StringLiteral string

func (s StringLiteral) toString() string  { return fmt.Sprintf("\"%s\"", s) }
func (s StringLiteral) eval(env Env) SExp { return s }

// List : e.g. (1 2 3)
type List []SExp

func (l List) toString() string {
	return toStringSexpSlice("(", []SExp(l), ")")
}

func (l List) eval(env Env) SExp {
	if len(l) == 0 {
		return l
	} else if len(l) > 1 {
		switch c := l[0].eval(env); c.(type) {
		case Closure: // apply
			args := make(List, len(l)-1)
			for i, elem := range l[1:] {
				args[i] = elem.eval(env)
			}

			return c.(Closure).apply(args)
		default:
			println("error: can't apply")
		}
	}
	return UNDEF
}

// Vector : e.g. [1 2 3]
type Vector []SExp

func (v Vector) toString() string {
	return toStringSexpSlice("[", []SExp(v), "]")
}

func (v Vector) eval(env Env) SExp {
	ret := make(Vector, len(v))
	for i, elem := range v {
		ret[i] = elem.eval(env)
	}
	return ret
}

// HashMap : {x 1, y 2}
type HashMap []SExp

func (hm HashMap) toString() string {
	return toStringSexpSlice("{", []SExp(hm), "}")
}

func (hm HashMap) eval(env Env) SExp {
	ret := make(HashMap, len(hm))
	for i, elem := range hm {
		ret[i] = elem.eval(env)
	}
	return ret
}

// Func : function
type Func func(env Env, args List) SExp

// Closure : function + environment
type Closure struct {
	env Env
	fun Func
}

func (c Closure) toString() string  { return "*Closure*" }
func (c Closure) eval(env Env) SExp { return c }

func (c Closure) apply(args List) SExp { return c.fun(c.env, args) }

// SExp : a S SExpression
type SExp interface {
	toString() string
	eval(Env) SExp
}

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
