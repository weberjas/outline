package test

import "fmt"

func printTest() {
	fmt.Println("vim-go")
}

func this_is_a_test() {
	fmt.Println("thisisatest")
}

type OneOfEach struct {
	name    string
	value   int
	decimal float64
}

func (o OneOfEach) Display() {
	fmt.Println(o.name)
	fmt.Println(o.value)
	fmt.Println(o.decimal)
}
