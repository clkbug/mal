package main

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
		fun: Func(func(_ Env, args List) SExp {
			s := 0
			for _, v := range args {
				switch v.(type) {
				case Int:
					s += int(v.(Int))
				default:
					panic("+ expects int, but got ...")
				}
			}
			return Int(s)
		}),
	}
	minus := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) SExp {
			s := int(args[0].(Int))
			for _, v := range args[1:] {
				switch v.(type) {
				case Int:
					s -= int(v.(Int))
				default:
					panic("- expects int, but got ...")
				}
			}
			return Int(s)
		}),
	}
	times := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) SExp {
			s := 1
			for _, v := range args {
				switch v.(type) {
				case Int:
					s *= int(v.(Int))
				default:
					panic("* expects int, but got ...")
				}
			}
			return Int(s)
		}),
	}
	div := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) SExp {
			s := int(args[0].(Int))
			for _, v := range args[1:] {
				switch v.(type) {
				case Int:
					s /= int(v.(Int))
				default:
					panic("/ expects int, but got ...")
				}
			}
			return Int(s)
		}),
	}
	list := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) SExp {
			return args
		}),
	}
	not := Closure{
		env: replEnv,
		fun: Func(func(_ Env, args List) SExp {
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
			return Bool(!b)
		}),
	}
	replEnv.set(Symbol("+"), plus)
	replEnv.set(Symbol("-"), minus)
	replEnv.set(Symbol("*"), times)
	replEnv.set(Symbol("/"), div)
	replEnv.set(Symbol("list"), list)
	replEnv.set(Symbol("not"), not)
}
