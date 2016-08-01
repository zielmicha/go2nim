package typeswitch

func main() {
	var a interface{} = 5
	switch v := a.(type) {
	case int:
		println("int", v)
	case string:
		println("string", v)
	case float, complex128:
		println("floaty", v)
	}
}
