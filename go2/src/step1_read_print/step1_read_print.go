package main

import (
	"bufio"
	"fmt"
	"os"
)

func read(scanner *bufio.Scanner) (SExp, error) {
	s := scanner.Text()
	if len(s) == 0 {
		os.Exit(0)
	}
	r := initReader(s)
	return r.readForm()
}

func eval(e SExp) SExp { return e }

func print(e SExp) {
	switch e.(type) {
	case nil:
		fmt.Println("hoge")

	}
	fmt.Println(e.toString())
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("user> ")
	for scanner.Scan() {
		s, err := read(scanner)
		if err != nil {
			fmt.Println(err)
		} else {
			print(eval(s))
		}
		fmt.Print("user> ")
	}
}
