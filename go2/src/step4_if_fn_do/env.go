package main

import (
	"errors"
	"strings"
)

type envInternal map[Symbol]SExp

func makeEnvInternal() envInternal { return make(envInternal) }

// Env : Environment
type Env struct {
	env     envInternal
	nextEnv *Env
}

func (e Env) set(sym Symbol, sexp SExp) {
	e.env[sym] = sexp
}

func (e Env) get(sym Symbol) (SExp, bool) {
	v, ok := e.env[sym]
	if ok || e.nextEnv == nil {
		return v, ok
	}
	return e.nextEnv.get(sym)
}

func (e Env) del(sym Symbol) {
	delete(e.env, sym)
}

func makeNewEnv(e Env) Env {
	return Env{
		env:     make(envInternal),
		nextEnv: &e,
	}
}

func (e Env) copy() Env {
	var ne Env
	if e.nextEnv == nil {
		ne.env = make(envInternal)
	} else {
		ne = e.nextEnv.copy()
	}
	for k, v := range e.env {
		ne.env[k] = v.copy()
	}
	return ne
}

var replEnv Env

func init() {
	replEnv = Env{
		env:     makeEnvInternal(),
		nextEnv: nil,
	}
	plus := CoreFunc(
		func(args List) (SExp, error) {
			s := 0
			for _, v := range args {
				switch v.(type) {
				case Int:
					s += int(v.(Int))
				default:
					return UNDEF, errors.New("invalid +'s argument")
				}
			}
			return Int(s), nil
		})
	minus := CoreFunc(
		func(args List) (SExp, error) {
			s := int(args[0].(Int))
			for _, v := range args[1:] {
				switch v.(type) {
				case Int:
					s -= int(v.(Int))
				default:
					return UNDEF, errors.New("invalid -'s argument")
				}
			}
			return Int(s), nil
		})
	times := CoreFunc(
		func(args List) (SExp, error) {
			s := 1
			for _, v := range args {
				switch v.(type) {
				case Int:
					s *= int(v.(Int))
				default:
					return UNDEF, errors.New("invalid *'s argument")
				}
			}
			return Int(s), nil
		})
	div := CoreFunc(
		func(args List) (SExp, error) {
			s := int(args[0].(Int))
			for _, v := range args[1:] {
				switch v.(type) {
				case Int:
					s /= int(v.(Int))
				default:
					return UNDEF, errors.New("invalid +'s argument")
				}
			}
			return Int(s), nil
		})
	cmp := func(f func(x, y int) bool) CoreFunc {
		return CoreFunc(func(args List) (SExp, error) {
			if len(args) < 2 {
				return UNDEF, errors.New("few arguments for <,<=,>,>=")
			}
			switch x := args[0].(type) {
			case Int:
				switch y := args[1].(type) {
				case Int:
					return Bool(f(int(x), int(y))), nil
				}
			}
			return UNDEF, errors.New("arguments for '<' should be Int")
		})
	}
	lt := cmp(func(x, y int) bool { return x < y })
	le := cmp(func(x, y int) bool { return x <= y })
	gt := cmp(func(x, y int) bool { return x > y })
	ge := cmp(func(x, y int) bool { return x >= y })
	eq := CoreFunc(func(args List) (SExp, error) {
		if len(args) < 2 {
			return UNDEF, errors.New("few arguments for =")
		}
		return Bool(args[0].isSame(args[1])), nil
	})
	list := CoreFunc(
		func(args List) (SExp, error) {
			return args, nil
		})
	listq := CoreFunc(
		func(args List) (SExp, error) {
			switch args[0].(type) {
			case List:
				return Bool(true), nil
			default:
				return Bool(false), nil
			}
		})
	emptyq := CoreFunc(
		func(args List) (SExp, error) {
			switch args[0].(type) {
			case List:
				return Bool(len(args[0].(List)) == 0), nil
			case Vector:
				return Bool(len(args[0].(Vector)) == 0), nil
			default:
				return Bool(false), nil
			}
		})
	count := CoreFunc(
		func(args List) (SExp, error) {
			switch args[0].(type) {
			case List:
				return Int(len(args[0].(List))), nil
			case Vector:
				return Int(len(args[0].(Vector))), nil
			default:
				return Int(0), nil
			}
		})
	not := CoreFunc(func(args List) (SExp, error) {
		b := true
		switch args[0].(type) {
		case NilType:
			b = false
		case Bool:
			b = bool(args[0].(Bool))
		case List:
			list := args[0].(List)
			b = len(list) == 0
		}
		return Bool(!b), nil
	})
	do := CoreFunc(
		func(args List) (SExp, error) {
			return args[len(args)-1], nil
		})
	prn := CoreFunc(
		func(args List) (SExp, error) {
			s := make([]string, len(args))
			for i, a := range args {
				s[i] = a.toString()
			}
			println(strings.Join(s, " "))
			return NIL, nil
		})
	replEnv.set(Symbol("+"), plus)
	replEnv.set(Symbol("-"), minus)
	replEnv.set(Symbol("*"), times)
	replEnv.set(Symbol("/"), div)
	replEnv.set(Symbol("<"), lt)
	replEnv.set(Symbol("<="), le)
	replEnv.set(Symbol(">"), gt)
	replEnv.set(Symbol(">="), ge)
	replEnv.set(Symbol("="), eq)
	replEnv.set(Symbol("list"), list)
	replEnv.set(Symbol("list?"), listq)
	replEnv.set(Symbol("empty?"), emptyq)
	replEnv.set(Symbol("count"), count)
	replEnv.set(Symbol("not"), not)
	replEnv.set(Symbol("do"), do)
	replEnv.set(Symbol("prn"), prn)
}
