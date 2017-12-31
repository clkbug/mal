package main

import "fmt"

// Symbol : Symbol
type Symbol string

// StringLiteral : should be print with '"'
type StringLiteral string

// SExp : a S SExpression
type SExp interface {
}

func toString(s SExp) string {
	switch s.(type) {
	case int:
		return fmt.Sprint(s)
	case StringLiteral:
		str := string(s.(StringLiteral))
		return fmt.Sprintf("\"%s\"", str)
	case Symbol:
		sym := s.(Symbol)
		return string(sym)
	case []SExp:
		sexps := s.([]SExp)
		t := make([]byte, 0, 10)
		t = append(t, "("...)
		for i, v := range sexps {
			t = append(t, toString(v)...)
			if i != len(sexps)-1 {
				t = append(t, " "...)
			}
		}
		t = append(t, ")"...)
		return string(t)
	default:
		fmt.Print(s)
		panic("undefined type")
	}
}
