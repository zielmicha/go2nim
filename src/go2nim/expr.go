package go2nim

import "go/ast"
import "strings"

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
		panic("not implemented")
	case *ast.IndexExpr:
		panic("not implemented")
	case *ast.KeyValueExpr:
		panic("not implemented")
	case *ast.ParenExpr:
		panic("not implemented")
	case *ast.SelectorExpr:
		return c.convertExpr(node.X) + "." + c.convertFieldName(node.Sel.Name)
	case *ast.SliceExpr:
		panic("not implemented")
	case *ast.StarExpr:
		return c.convertExpr(node.X) + "[]"
	case *ast.TypeAssertExpr:
		panic("not implemented")
	case *ast.UnaryExpr:
		panic("not implemented")

	case *ast.BasicLit:
		// TODO: convert string literals
		return node.Value
	case *ast.CompositeLit:
		panic("not implemented")
	case *ast.FuncLit:
		panic("not implemented")

	case *ast.Ident:
		return c.convertFuncName(node.Name)

	default:
		panic("bad expr type")
	}
}
