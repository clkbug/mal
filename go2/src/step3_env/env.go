package main

// Env : Symbol -> S Expression
type Env map[Symbol]SExp

var replEnv Env

func (e Env) set(sym Symbol, sexp SExp) {
	e[sym] = sexp
}

func (e Env) get(sym Symbol) (SExp, bool) {
	v, ok := e[sym]
	return v, ok
}

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
	replEnv.set(Symbol("+"), plus)
	replEnv.set(Symbol("-"), minus)
	replEnv.set(Symbol("*"), times)
	replEnv.set(Symbol("/"), div)
}
