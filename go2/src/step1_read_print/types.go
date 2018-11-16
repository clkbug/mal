package main

import "fmt"

// Undefined : Undefined symbol. When an error occurred, reader returns UNDEF and err
type Undefined int

func (u Undefined) toString() string {
	return "*Undefined*"
}

// UNDEF : Undef
const UNDEF = Undefined(0)

// Int : integer
type Int int

func (i Int) toString() string {
	return fmt.Sprint(i)
}

// Symbol : Symbol
type Symbol string

func (s Symbol) toString() string {
	return string(s)
}

// StringLiteral : should be print with '"'
type StringLiteral string

func (s StringLiteral) toString() string {
	return fmt.Sprintf("\"%s\"", s)
}

// List : e.g. (1 2 3)
type List []SExp

func (l List) toString() string {
	return toStringSexpSlice("(", []SExp(l), ")")
}

// Vector : e.g. [1 2 3]
type Vector []SExp

func (v Vector) toString() string {
	return toStringSexpSlice("[", []SExp(v), "]")
}

// HashMap : {x 1, y 2}
type HashMap []SExp

func (hm HashMap) toString() string {
	return toStringSexpSlice("{", []SExp(hm), "}")
}

// SExp : a S SExpression
type SExp interface {
	toString() string
}

func toStringSexpSlice(ls string, sexps []SExp, rs string) string {
	t := make([]byte, 0, 10)
	t = append(t, ls...)
	for i, v := range sexps {
		t = append(t, v.toString()...)
		if i != len(sexps)-1 {
			t = append(t, " "...)
		}
	}
	t = append(t, rs...)
	return string(t)
}
