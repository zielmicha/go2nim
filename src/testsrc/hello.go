package main
import "fmt"

type Bar int
type FooFunc func(a int) (boo string)

func (b Bar) Rebarize() {

}

func main() {
    fmt.Println("hello world")
}
