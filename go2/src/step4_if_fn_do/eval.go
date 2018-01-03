package main

import (
	"errors"
	"fmt"
)

var specialFormsSet map[string]struct{}

func isSpecialForm(s SExp) bool {
	switch s.(type) {
	case Symbol:
		name := string(s.(Symbol))
		_, b := specialFormsSet[name]
		return b
	default:
		return false
	}
}
func initSpecialFormSet() {
	specialFormsSet = make(map[string]struct{})
	specialFormsSet["if"] = struct{}{}
	specialFormsSet["fn*"] = struct{}{}
	specialFormsSet["cond"] = struct{}{}
	specialFormsSet["or"] = struct{}{}
	specialFormsSet["and"] = struct{}{}
	specialFormsSet["def!"] = struct{}{}
	specialFormsSet["defmacro!"] = struct{}{}
	specialFormsSet["let*"] = struct{}{}
	specialFormsSet["do"] = struct{}{}
	specialFormsSet["quote"] = struct{}{}
	specialFormsSet["unquote"] = struct{}{}
	specialFormsSet["quasiquote"] = struct{}{}
	specialFormsSet["splice-unquote"] = struct{}{}

}

func eval(e SExp) (SExp, error) { return evalAst(replEnv, e) }

func evalAst(env Env, exp SExp) (SExp, error) {
	switch exp.(type) {
	case int, string, StringLiteral, Keyword:
		return exp, nil
	case Symbol:
		sym := exp.(Symbol)
		switch sym {
		case "true":
			return true, nil
		case "false":
			return false, nil
		case "nil":
			return exp, nil
		}
		val, ok := env.get(sym)
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

		if isSpecialForm(car) {
			return evalSpecialForm(env, list)
		}

		f, err := evalAst(env, car)
		if err != nil {
			return f, err
		}
		switch f.(type) {
		case func(List) SExp: // embedded functions
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
		case Fn:
			fun := f.(Fn)
			args := toSeq(fun.args)

			newEnv := make(map[Symbol]SExp)
			for k, v := range fun.env {
				newEnv[k] = v
			}
			for i, v := range cdr {
				argv, err := evalAst(env, v)
				if err != nil {
					return nil, err
				}
				newEnv[args[i].(Symbol)] = argv
			}
			return evalAst(newEnv, fun.body)

		}

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
	return exp, errors.New("eval!!!!!!!!!")
}

func evalSpecialForm(env Env, sexp SExp) (SExp, error) {
	s := []SExp(sexp.(List))
	name := s[0].(Symbol)
	switch name {
	case "def!":
		v := s[1].(Symbol)
		val, err := evalAst(env, s[2])
		if err != nil {
			return nil, err
		}
		gEnv.set(v, val)
		return val, nil
	case "let*":
		newEnv := make(map[Symbol]SExp)
		for k, v := range env {
			newEnv[k] = v
		}
		var tmp []SExp
		switch s[1].(type) {
		case List:
			tmp = []SExp(s[1].(List))
		case Vector:
			tmp = []SExp(s[1].(Vector))
		}
		for i := 0; i+1 < len(tmp); i += 2 {
			name := tmp[i].(Symbol)
			val, err := evalAst(newEnv, tmp[i+1])
			if err != nil {
				return nil, err
			}
			newEnv[name] = val
		}
		return evalAst(newEnv, s[2])
	case "if":
		p, err := evalAsBool(env, s[1])
		if err != nil {
			return nil, nil
		} else if p {
			return evalAst(env, s[2])
		} else if len(s) >= 4 {
			return evalAst(env, s[3])
		} else {
			return Symbol("nil"), nil
		}
	case "fn*":
		fn := Fn{
			args: s[1],
			body: s[2],
		}
		fn.env = make(map[Symbol]SExp)
		for k, v := range env {
			fn.env[k] = v
		}
		return fn, nil
	}
	return nil, nil
}

func evalAsBool(env Env, sexp SExp) (bool, error) {
	switch sexp.(type) {
	case bool:
		b := sexp.(bool)
		return b, nil
	case List:
		l := []SExp(sexp.(List))
		if len(l) == 0 {
			return false, nil
		} else {
			b, err := evalAst(env, sexp)
			if err != nil {
				return false, err
			}
			switch b.(type) {
			case bool:
				return b.(bool), nil
			case List:
				return true, nil
			default:
				return true, nil
			}
		}
	case Symbol:
		if b := string(sexp.(Symbol)); b == "true" {
			return true, nil
		} else if b == "false" || b == "nil" {
			return false, nil
		}
	}
	return true, nil
}
