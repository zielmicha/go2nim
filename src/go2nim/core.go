package go2nim

import "strings"
import "go/token"
import "go/ast"

type Context struct {
	Fset *token.FileSet
	resultVariables map[string]bool
	resultNames []string
	publicFunctions map[string]bool
	assignedTo map[string]bool
	iotaValue int
}

func NewContext() *Context {
	return &Context{
		Fset: token.NewFileSet(),
		publicFunctions: map[string]bool{},
		assignedTo: map[string]bool{}}
}

func (c *Context) copy() *Context {
	newc := *c
	return &newc
}

func (c *Context) downcaseFirstLetter(name string) string {
	firstLetter := strings.ToLower(string(name[0]))
	if firstLetter != string(name[0]) {
		if len(name) > 1 {
			secondLetter := strings.ToLower(string(name[1]))
			if secondLetter != string(name[1]) {
				return name
			}
		}
		return firstLetter + name[1:]
	} else {
		return name
	}
}

func (c *Context) upcaseFirstLetter(name string) string {
	firstLetter := strings.ToUpper(string(name[0]))
	if firstLetter != string(name[0]) {
		return firstLetter + name[1:]
	} else {
		return name
	}
}


func (c *Context) isPublic(name string) bool {
	firstLetter := strings.ToLower(string(name[0]))
	return firstLetter != string(name[0])
}

func (c *Context) convertFuncName(name string) string {
	if !c.isPublic(name) {
		if _, ok := c.publicFunctions[name]; ok {
			return c.downcaseFirstLetter(name) + "Internal"
		}
		return c.downcaseFirstLetter(name)
	}
	return c.downcaseFirstLetter(name)
}

func (c *Context) convertFieldName(name string) string {
	return c.downcaseFirstLetter(name)
}

func (c *Context) convertTypeName(name string) string {
	return c.upcaseFirstLetter(name)
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
	c.walkFuncDecls(allDecls)
	preamble += c.convertImports(allDecls)
	preamble += c.convertTypeDecls(allDecls)
	preamble += c.convertFuncDecls(allDecls)
	preamble += c.convertVariables(allDecls)

	return preamble
}

func (c *Context) ConvertBody(file *ast.File) string {
	c.walkFuncDecls(file.Decls)
	return c.convertFuncBodies(file.Decls)
}
