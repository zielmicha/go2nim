package go2nim

import "go/ast"
import "go/token"
import "strings"
import "strconv"

func (c *Context) walkStmt(stmt ast.Stmt) {
	var node interface{} = stmt
	switch node := node.(type) {
	case *ast.AssignStmt:
		for _, expr := range node.Lhs {
			if ident, ok := expr.(*ast.Ident); ok {
				c.assignedTo[ident.Name] = true
			}
		}
	case *ast.BlockStmt:
		if node == nil {
			return
		}
		for _, child := range node.List {
			c.walkStmt(child)
		}
	case *ast.CaseClause:
		for _, child := range node.Body {
			c.walkStmt(child)
		}
	case *ast.EmptyStmt:
		panic("not implemented")
	case *ast.ForStmt:
		if node.Init != nil {
			c.walkStmt(node.Init)
		}
		c.walkStmt(node.Body)
		c.walkStmt(node.Post)
	case *ast.IfStmt:

		if node.Else != nil {
			c.walkStmt(node.Else)
		}
		if node.Init != nil {
			c.walkStmt(node.Init)
		}
		c.walkStmt(node.Body)
	case *ast.SelectStmt:
		panic("not implemented")
	case *ast.SwitchStmt:
		if node.Init != nil {
			c.walkStmt(node.Init)
		}
		c.walkStmt(node.Body)
	case *ast.TypeSwitchStmt:
		c.walkStmt(node.Body)
	}
}

func (c *Context) convertStmt(stmt ast.Stmt) string {
	if stmt == nil {
		return ""
	}
	var node interface{} = stmt
	switch node := node.(type) {
	case *ast.AssignStmt:
		return c.convertAssign(node)
	case *ast.BlockStmt:
		return "block:" + indent(c.convertStmtList(node.List))
	case *ast.BranchStmt:
		// TODO: labels, goto
		switch node.Tok {
		case token.BREAK:
			return c.loops[len(c.loops)-1].breakStmt
		case token.CONTINUE:
			return c.loops[len(c.loops)-1].continueStmt
		case token.GOTO:
			return "break " + c.quoteLabel(node.Label.Name)
		default:
			panic("invalid branch keyword")
		}
	case *ast.DeclStmt:
		return c.convertDecl(node.Decl)
	case *ast.DeferStmt:
		// TODO: recover (use `return` from defer)
		// FIXME: args are evaluated imm
		return "godefer(" + c.convertExpr(node.Call) + ")"
	case *ast.EmptyStmt:
		panic("not implemented")
	case *ast.ExprStmt:
		return c.convertExpr(node.X)
	case *ast.ForStmt:
		return c.newScope().convertFor(node)
	case *ast.GoStmt:
		return "go " + c.convertExpr(node.Call)
	case *ast.IfStmt:
		newc := c.newScope()
		var initStr string
		if node.Init != nil {
			initStr = newc.convertStmt(node.Init)
		}
		cond := newc.convertExpr(node.Cond)
		if node.Init != nil {
			cond = "(" + initStr + "; " + cond + ")"
		}
		result := "if " + cond + ":\n"
		result += indent(newc.convertBlockStmt(node.Body))
		if node.Else != nil {
			switch node.Else.(type) {
			case *ast.BlockStmt:
				result += "\nelse:\n" + indent(c.newScope().convertBlockStmt(node.Else.(*ast.BlockStmt)))
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
		// ast.Print(c.Fset, node)
		var result string

		newc := c.newScope()
		if node.Value == nil {
			key := newc.convertExpr(node.Key)
			result = "for " + key + " in 0..<len(" + newc.convertExpr(node.X) + "):\n"
		} else {
			key := newc.convertExpr(node.Value)
			if node.Key != nil {
				key = newc.convertExpr(node.Key) + ", " + key
			}
			result = "for " + key + " in " + c.convertExpr(node.X) + ":\n"
		}
		newc.loops = append(newc.loops, Loop{continueStmt: "continue", breakStmt: "break"})
		result += indent(newc.convertBlockStmt(node.Body))
		return result
	case *ast.ReturnStmt:
		result := []string{}
		for _, expr := range node.Results {
			result = append(result, c.convertExpr(expr))
		}

		var resultStr string
		switch len(result) {
		case 0:
			resultStr = ""
		case 1:
			resultStr = result[0]
		default:
			resultStr = "(" + strings.Join(result, ", ") + ")"
		}

		if c.resultNames != nil && len(result) != 0 {
			resultVars := []string{}
			for _, name := range c.resultNames {
				resultVars = append(resultVars, "result." + name)
			}
			return "(" + strings.Join(resultVars, ", ") + ") = " + resultStr + "\nreturn"
		} else {
			return "return " + resultStr
		}
	case *ast.SelectStmt:
		panic("not implemented")
	case *ast.SendStmt:
		panic("not implemented")
	case *ast.SwitchStmt:
		return c.newScope().convertSwitch(node.Body, node.Init, node.Tag, false, "")
	case *ast.TypeSwitchStmt:
		//ast.Print(c.Fset, node.Assign)

		var typeSwitchVar string
		var tag ast.Expr

		if assign, ok := node.Assign.(*ast.AssignStmt); ok {
			if len(assign.Rhs) != 1 || len(assign.Lhs) > 1 {
				panic("invalid type switch")
			}
			if len(assign.Lhs) > 0 {
				typeSwitchVar = assign.Lhs[0].(*ast.Ident).Name
			}
			tag = assign.Rhs[0].(*ast.TypeAssertExpr).X
		} else {
			tag = node.Assign.(*ast.ExprStmt).X.(*ast.TypeAssertExpr).X
		}

		return c.convertSwitch(node.Body, node.Init, tag, true, typeSwitchVar)
	default:
		panic("bad stmt")
	}
}

func (c *Context) convertFor(node *ast.ForStmt) string {
	var cond string
	if node.Cond == nil {
		cond = "true"
	} else {
		cond = c.convertExpr(node.Cond)
	}
	loop := Loop{}
	loopName := "loop" + strconv.Itoa(c.loopCounter)
	c.loopCounter ++
	if node.Post != nil {
		loop.breakStmt = "break " + loopName
		loop.continueStmt = "break " + loopName + "Continue"
	} else {
		loop.breakStmt = "break"
		loop.continueStmt = "continue"
	}
	newc := c.copy()
	newc.loops = append(newc.loops, loop)

	var init string
	if node.Init != nil {
		init = c.convertStmt(node.Init) + "\n"
	}

	body := newc.convertBlockStmt(node.Body)
	if node.Post != nil {
		// FIXME: this is invalid together with continue/break
		body = "block " + loopName + "Continue:\n" + indent(body) + "\n" + c.convertStmt(node.Post)
	}
	result := "while " + cond + ":\n" + indent(body)
	if node.Init != nil {
		result = init + result
	}
	if node.Post != nil || node.Init != nil {
		result = "block " + loopName + ":\n" + indent(result)
	}
	return result

}

func (c *Context) convertSwitch(body *ast.BlockStmt, init ast.Stmt, tag ast.Expr, isTypeSwitch bool, typeSwitchVar string) string {
	// TODO: special case constant-only switches
	switchOn := "true"
	if tag != nil {
		switchOn = "condition"
	}

	var defaultCase *ast.CaseClause
	result := ""

	for i, clause := range body.List {
		clause := clause.(*ast.CaseClause)
		if clause.List == nil {
			defaultCase = clause
		} else {
			fallthroughIndex := i
			caseBody := clause.Body
			for len(clause.Body) > 0 {
				fallthroughStmt, ok := caseBody[len(caseBody)-1].(*ast.BranchStmt)
				if !ok { break }
				if fallthroughStmt.Tok != token.FALLTHROUGH { break }
				fallthroughIndex += 1
				caseBody = caseBody[:len(caseBody)-1]
				nestBody := body.List[fallthroughIndex].(*ast.CaseClause).Body
				caseBody = append(caseBody, nestBody...)
			}

			if result == "" {
				result += "if "
			} else {
				result += "elif "
			}
			cond := []string{}
			var commonType string
			for _, item := range clause.List {
				var expr string
				if isTypeSwitch {
					expr = "typeId(" + c.convertType(item) + ")"
					if commonType == "" {
						commonType = c.convertType(item)
					} else {
						commonType = "EmptyInterface"
					}
				} else {
					expr = c.convertExpr(item)
				}
				cond = append(cond, switchOn + " == " + expr)
			}
			result += strings.Join(cond, " or ") + ":\n"
			newc := c.newScope()
			newc.loops = append(newc.loops, Loop{continueStmt: "continue", breakStmt: "break"})
			if typeSwitchVar != "" {
				result += indent("let " + typeSwitchVar + " = castInterface(typeSwitchOn, to=" + commonType + ")") + "\n"
 			}
			result += indent(newc.convertStmtList(caseBody)) + "\n"
		}
	}

	if defaultCase != nil {
		result += "else:\n"
		if isTypeSwitch {
			result += indent("let " + typeSwitchVar + " = typeSwitchOn\n")
		}
		result += indent(c.newScope().convertStmtList(defaultCase.Body))
	}

	if tag != nil || init != nil {
		pre := "block:\n"
		if init != nil {
			pre += indent(c.convertStmt(init)) + "\n"
		}
		var convertedTag string
		if isTypeSwitch {
			pre += indent("let typeSwitchOn = " + c.convertExpr(tag)) + "\n"
			convertedTag = "typeId(typeSwitchOn)"
		} else {
			if tag != nil {
				convertedTag = c.convertExpr(tag)
			} else {
				convertedTag = "true"
			}
		}
		pre += indent("let condition = " + convertedTag + "\n")
		result = pre + indent(result)
	}
	return result
}

func (c *Context) convertAssign(node *ast.AssignStmt) string {
	lhsL := []string{}
	variablesRedeclared := []string{}
	variablesNotRedeclared := []string{}

	rhsL := []string{}
	if len(node.Rhs) == 1 && len(node.Lhs) > 1 {
		rhsL = []string{ c.convertExprMulti(node.Rhs[0]) }
	} else {
		for _, expr := range node.Rhs {
			rhsL = append(rhsL, c.convertExpr(expr))
		}
	}
	rhs := strings.Join(rhsL, ", ")

	for i, expr := range node.Lhs {
		var exprStr string
		if ident, ok := expr.(*ast.Ident); ok && ident.Obj != nil {
			if _, ok := ident.Obj.Decl.(*ast.Field); ok {
				exprStr = c.convertExpr(expr)
				variablesRedeclared = append(variablesRedeclared, "")
			} else {
				// not result.X
				exprStr = c.convertFieldName(ident.Name)

				if exprStr == "_" {
					exprStr = "unused"
				}

				varDef := "var " + exprStr + ": type((" + rhs + ")[" + strconv.Itoa(i) + "])"
				_, ok := c.currentScopeVariables[exprStr]
				if ok && (ident.Name != "_") {
					variablesRedeclared = append(variablesRedeclared, varDef)
				} else {
					variablesNotRedeclared = append(variablesNotRedeclared, varDef)
				}
				c.variableDeclared(ident.Name)
			}
		} else {
			exprStr = c.convertExpr(expr)
		}
		lhsL = append(lhsL, exprStr)
	}

	// Count...
	underscoreCount := 0
	nonunderscoreStmt := ""
	nonunderscoreStmtI := -1
	for i, val := range lhsL {
		if val == "_" {
			underscoreCount ++
		} else {
			nonunderscoreStmt = val
			nonunderscoreStmtI = i
		}
	}

	lhs := strings.Join(lhsL, ", ")

	if len(lhsL) > 1 {
		lhs = "(" + lhs + ")"
	}

	if len(rhsL) > 1 {
		rhs = "(" + rhs + ")"
	}

	if underscoreCount == len(lhsL) - 1 {
		lhs = nonunderscoreStmt
		if len(lhsL) != 1 {
			rhs = rhs + "[" + strconv.Itoa(nonunderscoreStmtI) + "]"
		}
	}
	if underscoreCount == len(lhsL) {
		return "discard " + rhs
	}

	switch(node.Tok) {
	case token.ASSIGN:
		return lhs + " = " + rhs
	case token.DEFINE:
		if len(variablesRedeclared) == 0 {
			return "var " + lhs + " = " + rhs
		} else if len(variablesNotRedeclared) == 0 {
			return lhs + " = " + rhs
		} else {
			return strings.Join(variablesNotRedeclared, "\n") + "\n" + lhs + " = " + rhs
		}
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
			switch(node.Tok) {
			case token.AND_ASSIGN:
				convertedToken = "and"
			case token.OR_ASSIGN:
				convertedToken = "or"
			case token.XOR_ASSIGN:
				convertedToken = "xor"
			case token.SHL_ASSIGN:
				convertedToken = "shl"
			case token.SHR_ASSIGN:
				convertedToken = "shr"
			case token.QUO_ASSIGN:
				// TODO: side effects?
				return lhs + " = `go/`(" + lhs + ", " + rhs + ")"
			default:
				ast.Print(c.Fset, node)
				panic("bad assigment token")
			}
			// TODO: side effects?
			return lhs + " = " + lhs + " " + convertedToken + " (" + rhs + ")"
		}
		return lhs + " " + convertedToken + " " + rhs
	}
}

func (c *Context) convertBlockStmt(stmts *ast.BlockStmt) string {
	if stmts == nil {
		return "discard"
	}
	return c.convertStmtList(stmts.List)
}

func (c *Context) convertStmtList(stmts []ast.Stmt) string {
	if len(stmts) == 0 {
		return "discard\n"
	}

	var label string = ""
	var labelPos int
	var labelStmt ast.Stmt
	for i, stmt := range stmts {
		if labeled, ok := stmt.(*ast.LabeledStmt); ok {
			label = labeled.Label.Name
			labelPos = i
			labelStmt = labeled.Stmt
		}
	}

	if label != "" {
		stmts[labelPos] = labelStmt

		var firstGoto int = -1
		for i, stmt := range stmts {
			if c.containsGoto(stmt, "") {
				firstGoto = i
				break
			}
		}

		if firstGoto != -1 {
			result := c.convertStmtList(stmts[:firstGoto])
			result += "\nblock " + c.quoteLabel(label) + ":\n" + indent(c.convertStmtList(stmts[firstGoto:labelPos])) + "\n"
			result += c.convertStmtList(stmts[labelPos:])
			return result
		}
	}

	result := []string{}
	for _, stmt := range stmts {
		result = append(result, c.convertStmt(stmt))
	}
	return strings.Join(result, "\n")
}

func (c *Context) convertFuncCode(funcType *ast.FuncType, body *ast.BlockStmt) string {
	newc := c.newScope()
	newc.resultVariables = map[string]bool{}
	for k, v := range c.resultVariables {
		newc.resultVariables[k] = v
	}
	newc.assignedTo = map[string]bool{}
	newc.walkStmt(body)
	newc.resultNames = c.getReturnTypeResultVariables(funcType.Results)

	if funcType.Results != nil {
		for _, field := range funcType.Results.List {
			for _, name := range field.Names {
				newc.resultVariables[name.Name] = true
			}
		}
	}

	reassignNames := []string{}
	for _, arg := range funcType.Params.List {
		for _, name := range arg.Names {
			if _, ok := newc.assignedTo[name.Name]; ok {
				reassignNames = append(reassignNames, name.Name)
			}
		}
	}

	head := ""
	if len(reassignNames) != 0 {
		namesStr := strings.Join(reassignNames, ",")
		if len(reassignNames) != 1 {
			namesStr = "(" + namesStr + ")"
		}
		head += "var " + namesStr + " = " + namesStr + "\n"
	}

	return head + newc.convertBlockStmt(body)
}

func (c *Context) convertFuncBody(fdecl *ast.FuncDecl) string {
	paramList := c.convertParamList(fdecl.Type.Params, fdecl.Recv)
	name := c.convertFuncName(fdecl.Name.Name)

	// ast.Print(c.Fset, fdecl.Body)
	body := c.convertFuncCode(fdecl.Type, fdecl.Body)

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
