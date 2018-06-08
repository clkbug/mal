package main

import "errors"

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

func makeNewEnv(e Env) Env {
	return Env{
		env:     make(envInternal),
		nextEnv: &e,
	}
}

var replEnv Env

func init() {
	replEnv = Env{
		env:     makeEnvInternal(),
		nextEnv: nil,
	}
	plus := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) (SExp, error) {
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
		}),
	}
	minus := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) (SExp, error) {
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
		}),
	}
	times := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) (SExp, error) {
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
		}),
	}
	div := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) (SExp, error) {
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
		}),
	}
	list := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) (SExp, error) {
			return args, nil
		}),
	}
	listq := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) (SExp, error) {
			switch args[0].(type) {
			case List:
				return Bool(true), nil
			default:
				return Bool(false), nil
			}
		}),
	}
	emptyq := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) (SExp, error) {
			switch args[0].(type) {
			case List:
				return Bool(len(args[0].(List)) == 0), nil
			case Vector:
				return Bool(len(args[0].(Vector)) == 0), nil
			default:
				return Bool(false), nil
			}
		}),
	}
	count := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) (SExp, error) {
			switch args[0].(type) {
			case List:
				return Int(len(args[0].(List))), nil
			case Vector:
				return Int(len(args[0].(Vector))), nil
			default:
				return Int(0), nil
			}
		}),
	}
	not := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) (SExp, error) {
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
		}),
	}
	do := Closure{
		env: replEnv,
		fun: Func(func(env Env, args List) (SExp, error) {
			var retVal SExp
			for _, v := range args {
				retVal, _ = v.eval(env)
			}
			return retVal, nil
		}),
	}
	replEnv.set(Symbol("+"), plus)
	replEnv.set(Symbol("-"), minus)
	replEnv.set(Symbol("*"), times)
	replEnv.set(Symbol("/"), div)
	replEnv.set(Symbol("list"), list)
	replEnv.set(Symbol("list?"), listq)
	replEnv.set(Symbol("empty?"), emptyq)
	replEnv.set(Symbol("count"), count)
	replEnv.set(Symbol("not"), not)
	replEnv.set(Symbol("do"), do)
}
