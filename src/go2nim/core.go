package go2nim

import "strings"
import "go/token"
import "go/ast"

type Loop struct {
	name string
	breakStmt string
	continueStmt string
}

type Context struct {
	Fset *token.FileSet
	resultVariables map[string]bool
	resultNames []string
	publicNames map[string]bool
	assignedTo map[string]bool
	loops []Loop
	loopCounter int
	iotaValue int

	currentScopeVariables map[string]bool
}

func NewContext() *Context {
	return &Context{
		Fset: token.NewFileSet(),
		publicNames: map[string]bool{},
		assignedTo: map[string]bool{},
		currentScopeVariables: map[string]bool{}}
}

func (c *Context) copy() *Context {
	newc := *c
	return &newc
}

func (c *Context) newScope() *Context {
	newc := c.copy()
	newc.currentScopeVariables = map[string]bool{}
	return newc
}

func (c *Context) variableDeclared(name string) {
	c.currentScopeVariables[name] = true
}

func fromLiteralString(item *ast.BasicLit) string {
	if item.Kind != token.STRING {
		panic("unexpected item type")
	}
	return item.Value[1:len(item.Value)-1] // TODO
}

func (c *Context) convertComment(comments *ast.CommentGroup) string {
	if comments == nil {
		return ""
	}
	result := ""
	for _, comment := range comments.List {
		text := comment.Text
		if strings.HasPrefix(text, "//") {
			result += " #" + text[2:]
		} else {
			panic("unknown comment syntax")
		}
	}
	return result
}

func (c *Context) ConvertPreamble(files []*ast.File) string {
	preamble := "include gosupport\n"
	allDecls := make([]ast.Decl, 0)
	for _, file := range files {
		allDecls = append(allDecls, file.Decls...)
	}
	c.walkDecls(allDecls)
	preamble += c.convertImports(allDecls)
	preamble += c.convertVariables(allDecls, token.CONST)
	preamble += c.convertTypeDecls(allDecls)
	preamble += c.convertFuncDecls(allDecls)
	preamble += c.convertVariables(allDecls, token.VAR)

	return preamble
}

func (c *Context) ConvertBody(file *ast.File) string {
	c.walkDecls(file.Decls)
	return c.convertFuncBodies(file.Decls)
}
