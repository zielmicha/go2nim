package go2nim

import "go/ast"
import "go/token"
import "strings"

func (c *Context) convertStmt(stmt ast.Stmt) string {
	var node interface{} = stmt
	switch node := node.(type) {
	case *ast.AssignStmt:
		return c.convertAssign(node)
	case *ast.BlockStmt:
		panic("not implemented")
	case *ast.BranchStmt:
		// TODO: labels, goto
		switch node.Tok {
		case token.BREAK:
			return "break"
		case token.CONTINUE:
			return "continue"
		default:
			panic("invalid branch keyword")
		}
	case *ast.DeclStmt:
		return c.convertDecl(node.Decl)
	case *ast.DeferStmt:
		// TODO: recover (use `return` from defer)
		return "defer: " + c.convertExpr(node.Call)
	case *ast.EmptyStmt:
		panic("not implemented")
	case *ast.ExprStmt:
		return c.convertExpr(node.X)
	case *ast.ForStmt:
		var cond string
		if node.Cond == nil {
			cond = "true"
		} else {
			cond = c.convertExpr(node.Cond)
		}
		body := c.convertBlockStmt(node.Body)
		if node.Post != nil {
			body += "\n" + c.convertStmt(node.Post)
		}
		result := "while " + cond + ":\n" + indent(body)
		if node.Init != nil {
			init := c.convertStmt(node.Init) + "\n"
			result = "block:\n" + indent(init + result)
		}
		return result
	case *ast.GoStmt:
		panic("not implemented")
	case *ast.IfStmt:
		cond := c.convertExpr(node.Cond)
		if node.Init != nil {
			cond = "(" + c.convertStmt(node.Init) + "; " + cond + ")"
		}
		result := "if " + cond + ":\n"
		result += indent(c.convertBlockStmt(node.Body))
		if node.Else != nil {
			switch node.Else.(type) {
			case *ast.BlockStmt:
				result += "\nelse:\n" + indent(c.convertBlockStmt(node.Else.(*ast.BlockStmt)))
			case *ast.IfStmt:
				result += "\nel" + c.convertStmt(node.Else)
			default:
				panic("unknown else")
			}
		}
		return result
	case *ast.IncDecStmt:
		result := c.convertExpr(node.X)
		if node.Tok == token.INC {
			result += " += 1"
		} else {
			result += " -= 1"
		}
		return result
	case *ast.LabeledStmt:
		ast.Print(c.Fset, node)
		return "LABELED"
	case *ast.RangeStmt:
		key := c.convertExpr(node.Value)
		if node.Key != nil {
			key = c.convertExpr(node.Key) + ", " + key
		}
		result := "for " + key + " in " + c.convertExpr(node.X) + ":\n"
		result += indent(c.convertBlockStmt(node.Body))
		return result
	case *ast.ReturnStmt:
		result := []string{}
		for _, expr := range node.Results {
			result = append(result, c.convertExpr(expr))
		}
		switch len(result) {
		case 0:
			return "return"
		case 1:
			return "return " + result[0]
		default:
			return "return (" + strings.Join(result, ", ") + ")"
		}
	case *ast.SelectStmt:
		panic("not implemented")
	case *ast.SendStmt:
		panic("not implemented")
	case *ast.SwitchStmt:
		return c.convertSwitch(node.Body, node.Init, node.Tag, false, "")
	case *ast.TypeSwitchStmt:
		//ast.Print(c.Fset, node.Assign)
		assign := node.Assign.(*ast.AssignStmt)
		if len(assign.Rhs) != 1 || len(assign.Lhs) > 1 {
			panic("invalid type switch")
		}
		var typeSwitchVar string
		if len(assign.Lhs) > 0 {
			typeSwitchVar = assign.Lhs[0].(*ast.Ident).Name
		}
		tag := assign.Rhs[0].(*ast.TypeAssertExpr).X
		return c.convertSwitch(node.Body, node.Init, tag, true, typeSwitchVar)
	default:
		panic("bad stmt")
	}
}

func (c *Context) convertSwitch(body *ast.BlockStmt, init ast.Stmt, tag ast.Expr, isTypeSwitch bool, typeSwitchVar string) string {
	// TODO: special case constant-only switches
	switchOn := "true"
	if tag != nil {
		switchOn = "condition"
	}

	var defaultCase *ast.CaseClause
	result := ""

	for _, clause := range body.List {
		clause := clause.(*ast.CaseClause)
		if clause.List == nil {
			defaultCase = clause
		} else {
			if result == "" {
				result += "if "
			} else {
				result += "elif "
			}
			cond := []string{}
			for _, item := range clause.List {
				var expr string
				if isTypeSwitch {
					expr = "typeId(" + c.convertType(item) + ")"
				} else {
					expr = c.convertExpr(item)
				}
				cond = append(cond, switchOn + " == " + expr)
			}
			result += strings.Join(cond, ", ") + ":\n"
			result += indent(c.convertStmtList(clause.Body)) + "\n"
		}
	}

	if defaultCase != nil {
		result += "else:\n"
		result += indent(c.convertStmtList(defaultCase.Body))
	}

	if tag != nil || init != nil {
		pre := "block:\n"
		if init != nil {
			pre += indent(c.convertStmt(init)) + "\n"
		}
		var convertedTag string
		if isTypeSwitch {
			convertedTag = "typeId(" + c.convertExpr(tag) + ")"
		} else {
			convertedTag = c.convertExpr(tag)
		}
		pre += indent("let condition = " + convertedTag + "\n")
		result = pre + indent(result)
	}
	return result
}

func (c *Context) convertAssign(node *ast.AssignStmt) string {
	lhsL := []string{}
	for _, expr := range node.Lhs {
		lhsL = append(lhsL, c.convertExpr(expr))
	}

	rhsL := []string{}
	for _, expr := range node.Rhs {
		rhsL = append(rhsL, c.convertExpr(expr))
	}

	lhs := strings.Join(lhsL, ", ")
	rhs := strings.Join(rhsL, ", ")

	if len(lhsL) > 1 {
		lhs = "(" + lhs + ")"
	}

	if len(rhsL) > 1 {
		rhs = "(" + rhs + ")"
	}

	switch(node.Tok) {
	case token.ASSIGN:
		return lhs + " = " + rhs
	case token.DEFINE:
		return "var " + lhs + " = " + rhs
	default:
		if len(lhsL) != 1 {
			panic("bad assigment token")
		}
		var convertedToken string
		switch(node.Tok) {
		case token.ADD_ASSIGN:
			convertedToken = "+="
		case token.SUB_ASSIGN:
			convertedToken = "-="
		case token.MUL_ASSIGN:
			convertedToken = "*="
		default:
			panic("bad assigment token")
		}
		return lhs + " " + convertedToken + " " + rhs
	}
}

func (c *Context) convertBlockStmt(stmts *ast.BlockStmt) string {
	return c.convertStmtList(stmts.List)
}

func (c *Context) convertStmtList(stmts []ast.Stmt) string {
	if len(stmts) == 0 {
		return "discard\n"
	}

	result := []string{}
	for _, stmt := range stmts {
		result = append(result, c.convertStmt(stmt))
	}
	return strings.Join(result, "\n")
}

func (c *Context) convertFuncBody(fdecl *ast.FuncDecl) string {
	paramList := c.convertParamList(fdecl.Type.Params, fdecl.Recv)
	name := c.convertFuncName(fdecl.Name.Name)

	// ast.Print(c.Fset, fdecl.Body)
	body := c.convertBlockStmt(fdecl.Body)

	return "proc " + name + "(" + paramList + "): " + c.convertReturnType(fdecl.Type.Results) + " =\n" + indent(body) + "\n"
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
