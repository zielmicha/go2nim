package main

func call1(a int, bar ...string) {
	println("call ", len(bar))
}

func call2(a int, bar ...string) int {
	println("call ", len(bar))
	return 12
}

func call3(bar ...interface{}) int {
	println("call ", len(bar))
	return 12
}

func main() {
	call1(6, "1", "2")
	call1(6, ([]string{"a", "b"})...)
	//call2(6, "1", "2")
	call3(1, 2)
}
