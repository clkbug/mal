package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func read() string {
	return ""
}
func eval(s string) string { return s }
func print(s string) {
	fmt.Println(s)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("user> ")
	for scanner.Scan() {
		s := scanner.Text()
		s = strings.TrimRight(s, "\n")
		if len(s) == 0 {
			break
		}
		print(eval(s))
		fmt.Print("user> ")
	}
}
