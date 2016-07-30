package checkdefer

import "fmt"

func main() {
	defer (func() {
		fmt.Println("recovering...")
		v := recover()
		fmt.Println("recovered")
	})()
	panic("oops")
}
