package main

import "fmt"

func main() {
	//var a interface{} = 665
	//aint, ok := a.(int)

	b := interface{}("hello")
	c := interface{}("world")
	d := interface{}(665)

	sl := []interface{}{
		//a,
		b,
		c,
		d,
	}

	fmt.Println([]interface{}{"hello", "world", 88})
	fmt.Println("hello", "world", 100)
	fmt.Println(sl...)
}
