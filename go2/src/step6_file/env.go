package main

import (
	"errors"
	"io/ioutil"
	"os"
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
		func(args List, _ Env) (SExp, error) {
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
		func(args List, _ Env) (SExp, error) {
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
		func(args List, _ Env) (SExp, error) {
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
		func(args List, _ Env) (SExp, error) {
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
		return CoreFunc(func(args List, _ Env) (SExp, error) {
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
	eq := CoreFunc(func(args List, _ Env) (SExp, error) {
		if len(args) < 2 {
			return UNDEF, errors.New("few arguments for =")
		}
		return Bool(args[0].isSame(args[1])), nil
	})
	list := CoreFunc(
		func(args List, _ Env) (SExp, error) {
			return args, nil
		})
	listq := CoreFunc(
		func(args List, _ Env) (SExp, error) {
			switch args[0].(type) {
			case List:
				return Bool(true), nil
			default:
				return Bool(false), nil
			}
		})
	emptyq := CoreFunc(
		func(args List, _ Env) (SExp, error) {
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
		func(args List, _ Env) (SExp, error) {
			switch args[0].(type) {
			case List:
				return Int(len(args[0].(List))), nil
			case Vector:
				return Int(len(args[0].(Vector))), nil
			default:
				return Int(0), nil
			}
		})
	not := CoreFunc(func(args List, _ Env) (SExp, error) {
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
		func(args List, _ Env) (SExp, error) {
			return args[len(args)-1], nil
		})
	prstr := CoreFunc(
		func(args List, _ Env) (SExp, error) {
			s := printStrList(args, true, " ")
			return StringLiteral(s), nil
		})
	prn := CoreFunc(
		func(args List, _ Env) (SExp, error) {
			println(printStrList(args, true, " "))
			return NIL, nil
		})
	str := CoreFunc(
		func(args List, _ Env) (SExp, error) {
			s := ""
			for _, a := range args {
				s += a.printStr(false)
			}
			return StringLiteral(s), nil
		})
	printlnCF := CoreFunc(
		func(args List, _ Env) (SExp, error) {
			println(printStrList(args, false, " "))
			return NIL, nil
		})
	readString := CoreFunc(
		func(args List, _ Env) (SExp, error) {
			s, err := args[0].(StringLiteral)
			if !err {
				return NIL, errors.New("invalid read-string arg")
			}
			if len(s) == 0 {
				os.Exit(0)
			}
			r := initReader(string(s))
			return r.readForm()
		})
	evalCore := CoreFunc(
		func(args List, env Env) (SExp, error) {
			return args[0].eval(env)
		})
	slurp := CoreFunc(
		func(args List, env Env) (SExp, error) {
			fn, ok := args[0].(StringLiteral)
			if !ok {
				return NIL, errors.New("invalid slurp arg")
			}
			fp, err := os.Open(string(fn))
			if err != nil {
				return NIL, err
			}
			defer fp.Close()
			contents, err := ioutil.ReadAll(fp)
			if err != nil {
				return NIL, err
			}
			return StringLiteral(contents), nil
		})
	loadFile := CoreFunc(
		func(args List, env Env) (SExp, error) {
			s, err := slurp.apply(args, env)
			if err != nil {
				return NIL, err
			}
			str, ok := s.(StringLiteral)
			if !ok {
				return NIL, errors.New("!")
			}
			r := initReader(string(str))
			var val SExp = NIL
			for !r.isReachedEND {
				sexp, err := r.readForm()
				if err != nil {
					return NIL, err
				}
				buf, err := sexp.eval(env)
				if err != nil {
					return val, err
				}
				if buf != UNDEF {
					val = buf
				}
			}
			return val, nil
		})
	atom := CoreFunc(
		func(args List, env Env) (SExp, error) {
			return Atom{
				ref: args[0],
			}, nil
		})
	atomq := CoreFunc(
		func(args List, env Env) (SExp, error) {
			switch args[0].(type) {
			case Atom:
				return Bool(true), nil
			}
			return Bool(false), nil
		})
	deref := CoreFunc(
		func(args List, env Env) (SExp, error) {
			switch a := args[0].(type) {
			case Atom:
				return a.ref, nil
			}
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
	replEnv.set(Symbol("str"), str)
	replEnv.set(Symbol("pr-str"), prstr)
	replEnv.set(Symbol("println"), printlnCF)
	replEnv.set(Symbol("read-string"), readString)
	replEnv.set(Symbol("eval"), evalCore)
	replEnv.set(Symbol("slurp"), slurp)
	replEnv.set(Symbol("load-file"), loadFile)
	replEnv.set(Symbol("atom"), atom)
	replEnv.set(Symbol("atom?"), atomq)
	replEnv.set(Symbol("deref"), deref)
}

func printStrList(sexps List, isReadable bool, sep string) string {
	s := make([]string, len(sexps))
	for i, e := range sexps {
		s[i] = e.printStr(isReadable)
	}
	return strings.Join(s, sep)
}
