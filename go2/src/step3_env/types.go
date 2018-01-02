package main

import "fmt"

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

// SExp : a S SExpression
type SExp interface {
}

func toString(s SExp) string {
	switch s.(type) {
	case nil:
		return ""
	case int:
		return fmt.Sprint(s)
	case StringLiteral:
		str := string(s.(StringLiteral))
		return fmt.Sprintf("\"%s\"", str)
	case Symbol:
		sym := s.(Symbol)
		return string(sym)
	case Keyword:
		kw := s.(Keyword)
		return string(kw)
	case List:
		sexps := s.(List)
		return toStringSexpSlice("(", []SExp(sexps), ")")
	case Vector:
		sexps := s.(Vector)
		return toStringSexpSlice("[", []SExp(sexps), "]")
	case HashMap:
		sexps := s.(HashMap)
		return toStringSexpSlice("{", []SExp(sexps), "}")
	case error:
		return s.(error).Error()
	default:
		fmt.Print(s)
		panic("undefined type")
	}
}

func toStringSexpSlice(ls string, sexps []SExp, rs string) string {
	t := make([]byte, 0, 10)
	t = append(t, ls...)
	for i, v := range sexps {
		t = append(t, toString(v)...)
		if i != len(sexps)-1 {
			t = append(t, " "...)
		}
	}
	t = append(t, rs...)
	return string(t)
}
