package main

import "go/parser"
import "go/ast"
import "go2nim"
import "fmt"
import "os"

func main() {
	fn := os.Args[1]
	context := go2nim.NewContext()
	root, err := parser.ParseFile(context.Fset, fn, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	preamble := context.ConvertPreamble([]*ast.File{root})
	fmt.Println(preamble)
	body := context.ConvertBody(root)
	fmt.Println(body)
}
