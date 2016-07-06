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
			result += "  " + tspec.Name.Name + "* = " + c.convertType(tspec.Type) + "\n"
		}
	}
	if result != "" {
		result = "\ntype\n" + result
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

func (c *Context) convertVariables(decls []ast.Decl) string {
	result := ""
	for _, decl := range decls {
		if gdecl, ok := decl.(*ast.GenDecl); ok && (gdecl.Tok == token.CONST || gdecl.Tok == token.VAR) {
			result += c.convertVariableSection(gdecl.Tok, gdecl.Specs) + "\n"
		}
	}
	return result
}
