package go2nim

import "go/ast"
import "go/token"
import "strings"

func (c *Context) convertBinaryToken(tok token.Token) string {
	switch tok {
	case token.ADD:
		return "+"
	case token.SUB:
		return "-"
	case token.MUL:
		return "*"
	case token.QUO:
		return "/"
	case token.REM:
		return "mod"
	case token.AND, token.LAND:
        return "and"
    case token.OR, token.LOR:
        return "or"
    case token.XOR:
        return "xor"
    case token.SHL:
        return "shl"
    case token.SHR:
        return "shr"
	case token.EQL:
		return "=="
	case token.LSS:
		return "<"
	case token.GTR:
		return ">"
	case token.NEQ:
		return "!="
	case token.LEQ:
		return "<="
	case token.GEQ:
		return ">="
    default:
		panic("unknown binary token")
	}
}

func (c *Context) convertUnaryToken(tok token.Token) string {
	switch tok {
	case token.ADD:
		return "+"
	case token.SUB:
		return "-"
	case token.AND:
		return "gcaddr "
	case token.NOT:
		return "not "

    default:
		panic("unknown unary token ")
	}
}

func (c *Context) convertCompositeLiteral(expr *ast.CompositeLit, pointer bool) string {
	result := []string{}
	isKv := false
	for _, node := range expr.Elts {
		if kvexpr, ok := node.(*ast.KeyValueExpr); ok {
			isKv = true
			result = append(result, c.convertExpr(kvexpr.Key) + ": " + c.convertExpr(kvexpr.Value))
		} else {
			result = append(result, c.convertExpr(node))
		}
	}
	val := strings.Join(result, ", ")
	if isKv {
		val = "{" + val + "}"
	} else {
		val = "[" + val + "]"
	}
	return "make(" + c.convertType(expr.Type) + ", " + val + ")"
}

func (c *Context) convertExpr(expr ast.Expr) string {
	var node interface{} = expr

	switch node := node.(type) {
	case *ast.CallExpr:
		// TODO: ellipsis
		args := make([]string, len(node.Args))
		for i, arg := range node.Args {
			args[i] = c.convertExpr(arg)
		}
		return c.convertExpr(node.Fun) + "(" + strings.Join(args, ", ") + ")"
	case *ast.BinaryExpr:
		return c.convertExpr(node.X) + " " + c.convertBinaryToken(node.Op) + " " + c.convertExpr(node.Y)
	case *ast.IndexExpr:
		return c.convertExpr(node.X) + "[" + c.convertExpr(node.Index) + "]"
	case *ast.KeyValueExpr:
		panic("key value expr not expected here")
	case *ast.ParenExpr:
		return "(" + c.convertExpr(node.X) + ")"
	case *ast.SelectorExpr:
		return c.convertExpr(node.X) + "." + c.convertFieldName(node.Sel.Name)
	case *ast.SliceExpr:
		result := []string{c.convertExpr(node.X)}
		if node.Low != nil {
			result = append(result, "low=" + c.convertExpr(node.Low))
		}
		if node.High != nil {
			result = append(result, "high=" + c.convertExpr(node.Low))
		}
		if node.Max != nil {
			result = append(result, "max=" + c.convertExpr(node.Max))
		}
		return "slice(" + strings.Join(result, ", ") + ")"
	case *ast.StarExpr:
		return c.convertExpr(node.X) + "[]"
	case *ast.TypeAssertExpr:
		return "castInterface(" + c.convertExpr(node.X) + ", to=" + c.convertType(node.Type) + ")"
	case *ast.UnaryExpr:
		if node.Op == token.AND {
			if lit, ok := node.X.(*ast.CompositeLit); ok {
				return c.convertCompositeLiteral(lit, true)
			}
		}
		return c.convertUnaryToken(node.Op) + c.convertExpr(node.X)

	case *ast.BasicLit:
		// TODO: convert string literals
		return node.Value
	case *ast.CompositeLit:
		return c.convertCompositeLiteral(node, false)
	case *ast.FuncLit:
		paramList := c.convertParamList(node.Type.Params, nil)
		body := c.convertBlockStmt(node.Body)
		return "(proc(" + paramList + "): " + c.convertReturnType(node.Type.Results) + " =\n" + indent(body) + ")"
	case *ast.Ident:
		return c.convertFuncName(node.Name)

	default:
		return c.convertType(expr)
	}
}
