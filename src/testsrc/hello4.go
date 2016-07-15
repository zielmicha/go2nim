package main

type Foo struct {
}

func (f *Foo) func1() {
	println("hello")
}

func main() {
	var a Foo
	a.func1()
}
