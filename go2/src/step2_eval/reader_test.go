package main

import "testing"

func TestNext(test *testing.T) {
	check := func(r *Reader, e []Token) {
		t, err := r.next()
		if err != nil {
			test.Error(err)
		}
		for i := 0; t != ""; i++ {
			if t != e[i] {
				test.Errorf("Expected: %v\nbut actually got: %v\n", e[i], t)
			}
			t, err = r.next()
			if err != nil {
				test.Error(err)
			}
		}
	}

	r := initReader(" ( + 1 2 3 4 5) ")
	expected := []Token{"(", "+", "1", "2", "3", "4", "5", ")"}
	check(r, expected)

	r = initReader("(+ 1 2 ( * 3 4 5) 6 ( - 7 8 ( - 9 10)	\n  	  ) ) ")
	expected = []Token{"(", "+", "1", "2", "(", "*", "3", "4", "5", ")", "6", "(", "-", "7", "8", "(", "-", "9", "10", ")", ")", ")"}
	check(r, expected)

	r = initReader("(\"hoge\" fuga piyo); comment!!")
	expected = []Token{"(", "\"hoge\"", "fuga", "piyo", ")"}
	check(r, expected)

	r = initReader("\"hoge\"\"fuga\"~@")
	expected = []Token{"\"hoge\"", "\"fuga\"", "~@"}
	check(r, expected)

	r = initReader("hoge")
	expected = []Token{"hoge"}
	check(r, expected)

}
