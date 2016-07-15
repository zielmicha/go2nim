package main

func main() {
	if true { goto l1 }
	println("one")
	if true { goto l2 }
l1:
	println("two")
l2:
	println("three")
}
