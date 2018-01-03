package main

import (
	"strings"
)

func quoteMeta(s string) string {
	s = strings.Replace(s, "\n", "\\n", -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	return s
}

func unquoteMeta(s string) string {
	s = strings.Replace(s, "\\\\", "\\", -1)
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\\"", "\"", -1)
	return s
}

func toSeq(s SExp) []SExp {
	switch s.(type) {
	case List:
		return []SExp(s.(List))
	case Vector:
		return []SExp(s.(Vector))
	case HashMap:
		return []SExp(s.(HashMap))
	}
	panic("not Seq...")
}
