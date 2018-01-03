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

func print(e SExp) {
	switch e.(type) {
	case nil:
	default:
		fmt.Println(toString(e, true))
	}
}

func init() {
	initSpecialFormSet()
	initReplEnv()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("user> ")
	for scanner.Scan() {
		s, err := read(scanner)
		if err != nil {
			fmt.Println(err)
		} else {
			res, err := eval(s)
			if err != nil {
				fmt.Println(err)
			} else {
				print(res)
			}

		}
		fmt.Print("user> ")
	}
}
