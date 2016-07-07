package main

import "go/parser"
import "go/ast"
import "go2nim"
import "fmt"
import "os"

func main() {
	action := os.Args[1]
	fns := os.Args[2:]
	context := go2nim.NewContext()

	files := []*ast.File{}
	for _, fn := range fns {
		root, err := parser.ParseFile(context.Fset, fn, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		files = append(files, root)
	}
	printPreamble := action == "preamble" || action == "all"
	printBody := action == "body" || action == "all"
	if printPreamble {
		preamble := context.ConvertPreamble(files)
		fmt.Println(preamble)
	}
	if printBody {
		for _, file := range files {
			body := context.ConvertBody(file)
			fmt.Println(body)
		}
	}
	if !printBody && !printPreamble {
		panic("invalid action")
	}
}
