package main

import (
	"fmt"
	"reflect"
	"strings"
)

// Symbol : Symbol
type Symbol string

// Keyword : Keyword
type Keyword string

// StringLiteral : should be print with '"'
type StringLiteral string

// List : e.g. (1 2 3)
type List []SExp

// Vector : e.g. [1 2 3]
type Vector []SExp

// HashMap : {x 1, y 2}
type HashMap []SExp

// Fn : (fn* (x y z) (+ x y z))
type Fn struct {
	args SExp
	body SExp
	env  Env
}

// SExp : a S SExpression
type SExp interface {
}

func toString(s SExp, addQuotation bool) string {
	switch s.(type) {
	case nil:
		return ""
	case int, bool, string:
		return fmt.Sprint(s)
	case StringLiteral:
		str := string(s.(StringLiteral))
		if addQuotation {
			return fmt.Sprintf("\"%s\"", str)
		} else {
			return str
		}
	case Symbol:
		sym := s.(Symbol)
		return string(sym)
	case Keyword:
		kw := s.(Keyword)
		return string(kw)
	case List:
		sexps := s.(List)
		return toStringSexpSlice("(", []SExp(sexps), ")", addQuotation)
	case Vector:
		sexps := s.(Vector)
		return toStringSexpSlice("[", []SExp(sexps), "]", addQuotation)
	case HashMap:
		sexps := s.(HashMap)
		return toStringSexpSlice("{", []SExp(sexps), "}", addQuotation)
	case Fn:
		fn := s.(Fn)
		return fmt.Sprintf("(fn* %s %s)", toString(fn.args, addQuotation), toString(fn.body, addQuotation))
	case error:
		return s.(error).Error()
	default:
		fmt.Println(s, reflect.TypeOf(s))
		panic("undefined type")
	}
}

func toStringSexpSlice(ls string, sexps []SExp, rs string, addQuotation bool) string {
	t := make([]string, 0, 10)
	for _, v := range sexps {
		t = append(t, toString(v, addQuotation))
	}
	return ls + strings.Join(t, " ") + rs
}
