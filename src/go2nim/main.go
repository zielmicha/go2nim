package main

import "fmt"
import "strings"
import "os"
import "go/parser"
import "go/token"
import "go/ast"

type Context struct {
	fset *token.FileSet
	resultVariables map[string]bool
	iotaValue int
}

func (c *Context) copy() *Context {
	newc := *c
	return &newc
}

func fromLiteralString(item *ast.BasicLit) string {
	if item.Kind != token.STRING {
		panic("unexpected item type")
	}
	return item.Value
}

func (c *Context) convertType(expr ast.Expr) string {
	var node interface{} = expr
	switch node := node.(type) {
	case *ast.Ident:
		// TODO: translate builtin identifiers
		return node.Name
	case *ast.FuncType:
		return "(proc(" + c.convertParamList(node.Params, nil) + "): " + c.convertReturnType(node.Results) + ")"
	default:
		ast.Print(c.fset, expr)
		panic("bad type expression")
	}
}

func (c *Context) convertParamList(fields *ast.FieldList, recv *ast.FieldList) string {
	if fields == nil {
		return ""
	}

	result := make([]string, len(fields.List))
	for i, field := range fields.List {
		if len(field.Names) != 1 {
			ast.Print(c.fset, fields)
			panic("unexpected number of names")
		}
		result[i] = field.Names[0].Name + ": " + c.convertType(field.Type)
	}

	if recv != nil {
		receiver := recv.List[0]
		recvParam := receiver.Names[0].Name + ": " + c.convertType(receiver.Type)
		result = append([]string{recvParam}, result...)
	}

	return strings.Join(result, ", ")
}

func (c *Context) convertReturnType(fields *ast.FieldList) string {
	if fields == nil || len(fields.List) == 0 {
		return "void"
	} else if len(fields.List) == 1 {
		return c.convertType(fields.List[0].Type)
	} else {
		return "composite return type"
	}
}

func (c *Context) convertImports(decls []ast.Decl) string {
	result := ""
	for _, decl := range decls {
		if gdecl, ok := decl.(*ast.GenDecl); ok && gdecl.Tok == token.IMPORT {
			result += "import " + fromLiteralString(gdecl.Specs[0].(*ast.ImportSpec).Path) + "\n"
		}
	}
	return result
}

func (c *Context) convertTypeDecls(decls []ast.Decl) string {
	result := "\ntype\n"
	for _, decl := range decls {
		if gdecl, ok := decl.(*ast.GenDecl); ok && gdecl.Tok == token.TYPE {
			tspec := gdecl.Specs[0].(*ast.TypeSpec)
			result += "  " + tspec.Name.Name + "* = " + c.convertType(tspec.Type) + "\n"
		}
	}
	return result
}

func (c *Context) convertFuncName(name string) string {
	firstLetter := strings.ToLower(string(name[0]))
	if firstLetter != string(name[0]) {
		return firstLetter + name[1:]
	} else {
		return name
	}
}

func (c *Context) convertFuncDecl(fdecl *ast.FuncDecl) string {
	paramList := c.convertParamList(fdecl.Type.Params, fdecl.Recv)
	name := c.convertFuncName(fdecl.Name.Name)
	isPublic := name != fdecl.Name.Name
	if isPublic {
		name += "*"
	}
	return "proc " + name + "(" + paramList + "): " + c.convertReturnType(fdecl.Type.Results)
}

func (c *Context) convertFuncDecls(decls []ast.Decl) string {
	result := "\n"
	for _, decl := range decls {
		if fdecl, ok := decl.(*ast.FuncDecl); ok {
			result += c.convertFuncDecl(fdecl) + "\n"
		}
	}
	return result
}

func (c *Context) convertPreamble(files []*ast.File) string {
	preamble := ""
	allDecls := make([]ast.Decl, 0)
	for _, file := range files {
		allDecls = append(allDecls, file.Decls...)
	}
	preamble += c.convertImports(allDecls)
	preamble += c.convertTypeDecls(allDecls)
	preamble += c.convertFuncDecls(allDecls)
	return preamble
}

func main() {
	fn := os.Args[1]
	context := &Context{}
	context.fset = token.NewFileSet()
	root, err := parser.ParseFile(context.fset, fn, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Print(context.fset, root)
	var preamble = context.convertPreamble([]*ast.File{root})
	fmt.Println(preamble)
}
