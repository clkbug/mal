package main

import (
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

		if isSpecialForm(car) {
			return evalSpecialForm(env, list)
		}

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
		env.set(v, val)
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
	}
	return nil, nil
}
