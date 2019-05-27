package main

import "fmt"

func main() {
	s := "hello"
	fmt.Println(len(s))
	fmt.Println(s[0])
	fmt.Printf("%q\n", s[0])
	fmt.Printf("%b", s[0])
	r := "hello, 你好"
	fmt.Println(len(r))
}
