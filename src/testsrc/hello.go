package main
import "fmt"

type Bar int
type FooFunc func(a int) (boo string)
type TupleReturningFunc func(a int) (boo string, foo string)

type SomeStruct struct {
	a int
	b int
	// hello world
}


type SomeIface interface {
	Operation1() int
	Operation2(a int) string
}

const (
	one = 1
	two = 2
)

// nope

func (b Bar) Rebarize() {
	var a int = (5+5)*2
}

func main() {
    fmt.Println("hello world")
	fmt.Println("hello world 2")
}
