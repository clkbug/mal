package main

import (
	"fmt"
)

// Env : Symbol -> S Expression
type Env map[Symbol]SExp

var replEnv Env

func initReplEnv() {
	replEnv = make(map[Symbol]SExp)

	plus := func(args List) SExp {
		s := 0
		for _, v := range args {
			switch v.(type) {
			case int:
				s += v.(int)
			default:
				panic("+ expects int, but got ...")
			}
		}
		return s
	}
	minus := func(args List) SExp {
		s := args[0].(int)
		for _, v := range args[1:] {
			switch v.(type) {
			case int:
				s -= v.(int)
			default:
				panic("- expects int, but got ...")
			}
		}
		return s
	}
	times := func(args List) SExp {
		s := 1
		for _, v := range args {
			switch v.(type) {
			case int:
				s *= v.(int)
			default:
				panic("+ expects int, but got ...")
			}
		}
		return s
	}
	div := func(args List) SExp {
		s := args[0].(int)
		for _, v := range args[1:] {
			switch v.(type) {
			case int:
				s /= v.(int)
			default:
				panic("/ expects int, but got ...")
			}
		}
		return s
	}
	replEnv[Symbol("+")] = plus
	replEnv[Symbol("-")] = minus
	replEnv[Symbol("*")] = times
	replEnv[Symbol("/")] = div
}

func eval(e SExp) (SExp, error) { return evalAst(replEnv, e) }

func evalAst(env Env, exp SExp) (SExp, error) {
	switch exp.(type) {
	case int, string, Keyword:
		return exp, nil
	case Symbol:
		sym := exp.(Symbol)
		val, ok := env[sym]
		if ok {
			return val, nil
		}
		return val, fmt.Errorf("Error: '%s' not found", sym)
	case List:
		list := exp.(List)
		if len(list) == 0 {
			return list, nil
		}
		car := list[0]
		cdr := list[1:]
		f, err := evalAst(env, car)
		if err != nil {
			return f, err
		}
		fun := f.(func(List) SExp)
		args := make([]SExp, len(cdr))
		for i, v := range cdr {
			var err error
			args[i], err = evalAst(env, v)
			if err != nil {
				return nil, err
			}
		}
		return fun(args), nil
	case Vector:
		vect := exp.(Vector)
		res := make([]SExp, len(vect))
		for i, v := range vect {
			var err error
			res[i], err = evalAst(env, v)
			if err != nil {
				return nil, err
			}
		}
		return Vector(res), nil
	case HashMap:
		hash := exp.(HashMap)
		res := make([]SExp, len(hash))
		for i, v := range hash {
			var err error
			res[i], err = evalAst(env, v)
			if err != nil {
				return nil, err
			}
		}
		return HashMap(res), nil
	}
	return exp, nil
}
