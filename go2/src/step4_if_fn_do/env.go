package main

import (
	"reflect"
)

// Env : Symbol -> S Expression
type Env map[Symbol]SExp

var replEnv Env
var gEnv Env

func (e Env) set(sym Symbol, sexp SExp) {
	e[sym] = sexp
}

func (e Env) get(sym Symbol) (SExp, bool) {
	v, ok := e[sym]
	if !ok {
		v, ok = gEnv[sym]
	}
	return v, ok
}

func initReplEnv() {
	gEnv = make(map[Symbol]SExp)
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
	eq := func(args List) SExp {
		for i := 0; i < 2; i++ {
			switch args[i].(type) {
			case List:
				args[i] = []SExp(args[i].(List))
			case Vector:
				args[i] = []SExp(args[i].(Vector))
			case HashMap:
				args[i] = []SExp(args[i].(HashMap))
			}
		}
		return reflect.DeepEqual(args[0], args[1])
	}
	lt := func(args List) SExp {
		return args[0].(int) < args[1].(int)
	}
	le := func(args List) SExp {
		return args[0].(int) <= args[1].(int)
	}
	gt := func(args List) SExp {
		return args[0].(int) > args[1].(int)
	}
	ge := func(args List) SExp {
		return args[0].(int) >= args[1].(int)
	}
	listP := func(args List) SExp {
		switch args[0].(type) {
		case List:
			return true
		default:
			return false
		}
	}
	list := func(args List) SExp {
		return args
	}
	empty := func(args List) SExp {
		switch args[0].(type) {
		case List:
			if len([]SExp(args[0].(List))) == 0 {
				return true
			}
		case Vector:
			if len([]SExp(args[0].(Vector))) == 0 {
				return true
			}
		case HashMap:
			if len([]SExp(args[0].(HashMap))) == 0 {
				return true
			}
		}
		return false
	}
	not := func(args List) SExp {
		switch args[0].(type) {
		case bool:
			b := args[0].(bool)
			return !b
		case Symbol:
			s := string(args[0].(Symbol))
			if s == "false" {
				return true
			}
		}
		return false
	}
	count := func(args List) SExp {
		switch args[0].(type) {
		case List:
			return len([]SExp(args[0].(List)))
		case Vector:
			return len([]SExp(args[0].(Vector)))
		case HashMap:
			return len([]SExp(args[0].(HashMap)))
		case Symbol:
			if args[0] == Symbol("nil") {
				return 0
			}
		}
		return 1
	}
	str := func(args List) SExp {
		if a := []SExp(args); len(a) == 0 {
			return StringLiteral("")
		} else {
			s := make([]byte, 0, 10)
			for _, v := range a {
				switch v.(type) {
				case StringLiteral:
					s = append(s, string(v.(StringLiteral))...)
				default:
					s = append(s, toString(v, false)...)
				}
			}
			return StringLiteral(string(s))
		}
	}
	replEnv.set(Symbol("+"), plus)
	replEnv.set(Symbol("-"), minus)
	replEnv.set(Symbol("*"), times)
	replEnv.set(Symbol("/"), div)
	replEnv.set(Symbol("="), eq)
	replEnv.set(Symbol("<"), lt)
	replEnv.set(Symbol("<="), le)
	replEnv.set(Symbol(">"), gt)
	replEnv.set(Symbol(">="), ge)
	replEnv.set(Symbol("list?"), listP)
	replEnv.set(Symbol("list"), list)
	replEnv.set(Symbol("empty?"), empty)
	replEnv.set(Symbol("not"), not)
	replEnv.set(Symbol("count"), count)
	replEnv.set(Symbol("str"), str)
}
