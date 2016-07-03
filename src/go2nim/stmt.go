package go2nim

import "go/ast"

func (c *Context) convertStmt(stmt ast.Stmt) string {
	var node interface{} = stmt
	switch node := node.(type) {
	case *ast.AssignStmt:
		ast.Print(c.Fset, node)
		return "assign stmt"
	case *ast.BlockStmt:
		panic("not implemented")
	case *ast.BranchStmt:
		panic("not implemented")
	case *ast.DeclStmt:
		return c.convertDecl(node.Decl)
	case *ast.DeferStmt:
		panic("not implemented")
	case *ast.EmptyStmt:
		panic("not implemented")
	case *ast.ExprStmt:
		return c.convertExpr(node.X)
	case *ast.ForStmt:
		panic("not implemented")
	case *ast.GoStmt:
		panic("not implemented")
	case *ast.IfStmt:
		panic("not implemented")
	case *ast.IncDecStmt:
		panic("not implemented")
	case *ast.LabeledStmt:
		ast.Print(c.Fset, node)
		return "LABELED"
	case *ast.RangeStmt:
		panic("not implemented")
	case *ast.ReturnStmt:
		ast.Print(c.Fset, node)
		return "return something"
	case *ast.SelectStmt:
		panic("not implemented")
	case *ast.SendStmt:
		panic("not implemented")
	case *ast.SwitchStmt:
		panic("not implemented")
	case *ast.TypeSwitchStmt:
		panic("not implemented")
	default:
		panic("bad stmt")
	}
}

func (c *Context) convertBlockStmt(stmts *ast.BlockStmt) string {
	if len(stmts.List) == 0 {
		return "discard\n"
	}

	result := ""
	for _, stmt := range stmts.List {
		result += c.convertStmt(stmt) + "\n"
	}
	return result
}

func (c *Context) convertFuncBody(fdecl *ast.FuncDecl) string {
	paramList := c.convertParamList(fdecl.Type.Params, fdecl.Recv)
	name := c.convertFuncName(fdecl.Name.Name)

	body := c.convertBlockStmt(fdecl.Body)

	return "proc " + name + "(" + paramList + "): " + c.convertReturnType(fdecl.Type.Results) + " =\n" + indent(body)
}

func (c *Context) convertFuncBodies(decls []ast.Decl) string {
	result := "\n"
	for _, decl := range decls {
		if fdecl, ok := decl.(*ast.FuncDecl); ok {
			result += c.convertFuncBody(fdecl) + "\n"
		}
	}
	return result
}
