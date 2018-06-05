package main

// Env : Environment
type Env map[Symbol]SExp

func (e Env) set(sym Symbol, sexp SExp) {
	e[sym] = sexp
}

func (e Env) get(sym Symbol) (SExp, bool) {
	v, ok := e[sym]
	return v, ok
}

var replEnv Env

func initREPLEnv() {
	replEnv = make(map[Symbol]SExp)
	plus := Func(func(args List) SExp {
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
	})
	minus := Func(func(args List) SExp {
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
	})
	times := Func(func(args List) SExp {
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
	})
	div := Func(func(args List) SExp {
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
	})
	replEnv.set(Symbol("+"), plus)
	replEnv.set(Symbol("-"), minus)
	replEnv.set(Symbol("*"), times)
	replEnv.set(Symbol("/"), div)
}
