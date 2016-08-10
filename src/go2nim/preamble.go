package go2nim

import "go/token"
import "go/ast"

func (c *Context) convertImportPath(path string) string {
	return "golib/" + path
}

func (c *Context) convertImport(spec *ast.ImportSpec) string {
	var result string
	path := c.convertImportPath(fromLiteralString(spec.Path))
	var name string
	if spec.Name == nil {
		name = "_"
	} else {
		name = spec.Name.Name
	}

	switch name {
	case ".", "_":
		result = "import " + path
	default:
		result = "import " + path + " as " + spec.Name.Name
	}
	result += c.convertComment(spec.Comment) + "\n"
	return result
}

func (c *Context) convertImports(decls []ast.Decl) string {
	result := ""
	for _, decl := range decls {
		if gdecl, ok := decl.(*ast.GenDecl); ok && gdecl.Tok == token.IMPORT {
			for _, spec := range gdecl.Specs {
				result += c.convertImport(spec.(*ast.ImportSpec))
			}
		}
	}
	return result
}

func (c *Context) convertTypeDecls(decls []ast.Decl) string {
	result := ""
	for _, decl := range decls {
		if gdecl, ok := decl.(*ast.GenDecl); ok && gdecl.Tok == token.TYPE {
			tspec := gdecl.Specs[0].(*ast.TypeSpec)
			result += "  " + c.convertTypeName(tspec.Name.Name) + "* = " + c.convertTypeDecl(tspec.Type) + "\n"
		}
	}
	if result != "" {
		result = "\nexttypes:\n  type\n" + indent(result)
	}
	return result
}

func (c *Context) convertFuncDecl(fdecl *ast.FuncDecl) string {
	paramList := c.convertParamList(fdecl.Type.Params, fdecl.Recv)
	name := c.convertFuncName(fdecl.Name.Name)
	isPublic := name != fdecl.Name.Name
	if isPublic {
		name += "*"
	}
	result := "proc " + name + "(" + paramList + "): " + c.convertReturnType(fdecl.Type.Results) + " {.gofunc"
	if fdecl.Recv != nil {
		if _, ok := fdecl.Recv.List[0].Type.(*ast.StarExpr); ok {
			result += ", gomethod"
		}
	}
	return result + ".}"
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

func (c *Context) walkDecls(decls []ast.Decl) {
	for _, decl := range decls {
		if gdecl, ok := decl.(*ast.GenDecl); ok && (gdecl.Tok == token.CONST || gdecl.Tok == token.VAR) {
			for _, gspec := range gdecl.Specs {
				spec := gspec.(*ast.ValueSpec)
				for _, name := range spec.Names {
					if c.isPublic(name.Name) {
						c.publicNames[c.downcaseFirstLetter(name.Name)] = true
					}
				}
			}
		}

		if fdecl, ok := decl.(*ast.FuncDecl); ok {
			name := fdecl.Name.Name
			if c.isPublic(name) {
				c.publicNames[c.downcaseFirstLetter(name)] = true
			}
		}

		if gdecl, ok := decl.(*ast.GenDecl); ok && gdecl.Tok == token.TYPE {
			tspec := gdecl.Specs[0].(*ast.TypeSpec)
			name := tspec.Name.Name
			if c.isPublic(name) {
				c.publicNames[c.upcaseFirstLetter(name)] = true
			}
		}
	}
}

func (c *Context) sortDecls(decls []ast.Decl) []ast.Decl {
	// topological ordering of dependency graph
	result := []ast.Decl{}
	visited := map[*ast.Decl]bool{}
	nameToDecl := map[string]*ast.Decl{}
	for i, decl := range decls {
		for _, ident := range c.getAssignedIdentifiers(decl) {
			nameToDecl[ident] = &decls[i]
		}
	}

	var visit func(decl *ast.Decl)
	visit = func(decl *ast.Decl) {
		if _, ok := visited[decl]; ok {
			return
		}
		visited[decl] = true

		for _, ident := range c.getIdentifiers(*decl) {
			if childdecl, ok := nameToDecl[ident]; ok {
				visit(childdecl)
			}
		}
		result = append(result, *decl)
	}

	for i := range decls {
		visit(&decls[i])
	}

	return result
}

func (c *Context) convertVariables(decls []ast.Decl, which token.Token) string {
	result := ""
	gdecls := []ast.Decl{}
	for _, decl := range decls {
		if gdecl, ok := decl.(*ast.GenDecl); ok && (gdecl.Tok == which) {
			gdecls = append(gdecls, gdecl)
		}
	}

	gdecls = c.sortDecls(gdecls)
	for _, decl := range gdecls {
		gdecl := decl.(*ast.GenDecl)
		result += c.convertVariableSection(gdecl.Tok, gdecl.Specs, true) + "\n"
	}
	return result
}
