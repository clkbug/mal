package main

import "errors"

const IF = "if"
const COND = "cond"
const OR = "or"
const DEF = "def!"
const DEFMACRO = "defmacro!"
const LET = "let*"

var specialFormMap = map[string]struct{}{
	IF: struct{}{}, COND: struct{}{}, OR: struct{}{},
	DEF: struct{}{}, DEFMACRO: struct{}{}, LET: struct{}{},
}

func isSpecialForm(s SExp) (string, bool) {
	switch s.(type) {
	case Symbol:
		_, ok := specialFormMap[string(s.(Symbol))]
		return string(s.(Symbol)), ok
	default:
		return "", false
	}
}

func evalIf(env Env, l List) (SExp, error) {
	cond := true
	switch l[0].(type) {
	case NilType:
		cond = false
	case Bool:
		cond = bool(l[0].(Bool))
	case List:
		list := l[0].(List)
		cond = len(list) == 0
	}
	if cond {
		return l[1].eval(env)
	} else if len(l) >= 3 {
		return l[2].eval(env)
	} else {
		return NIL, nil
	}
}

func evalDef(env Env, l List) (SExp, error) {
	switch l[0].(type) {
	case Symbol:
		s := l[0].(Symbol)
		v, err := l[1].eval(env)
		if err != nil {
			return UNDEF, err
		}
		env.set(s, v)
		return v, nil
	default:
		return UNDEF, errors.New("'(def! SYMBOL EXP)'")
	}
}

func evalLet(env Env, l List) (SExp, error) {
	switch l[0].(type) {
	case List:
		vars := l[0].(List)
		body := l[1]
		tmpEnv := makeNewEnv(env)
		if len(vars)%2 != 0 {
			return UNDEF, errors.New("Syntax Error: let*'s bind")
		}
		for i := 0; i < len(vars); i += 2 {
			evalLetBindOne(tmpEnv, vars[i:i+2])
		}
		return body.eval(tmpEnv)

	case Vector:
		vars := l[0].(Vector).toList()
		body := l[1]
		tmpEnv := makeNewEnv(env)
		if len(vars)%2 != 0 {
			return UNDEF, errors.New("Syntax Error: let*'s bind")
		}
		for i := 0; i < len(vars); i += 2 {
			evalLetBindOne(tmpEnv, vars[i:i+2])
		}
		return body.eval(tmpEnv)
	default:
		return UNDEF, errors.New("Syntax error: let*")
	}
}

func evalLetBindOne(env Env, l List) error {
	switch l[0].(type) {
	case Symbol:
		vname := l[0].(Symbol)
		value, err := l[1].eval(env)
		if err != nil {
			return err
		}
		env.set(vname, value)
		return nil
	default:
		return errors.New("Syntax error: let*'s bind")
	}
}
