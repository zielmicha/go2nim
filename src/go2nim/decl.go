package go2nim

import "go/ast"
import "go/token"
import "strings"

func (c *Context) convertDecl(decl ast.Decl) string {
	var node interface{} = decl
	switch node := node.(type) {
	case *ast.GenDecl:
		switch node.Tok {
		case token.CONST, token.VAR:
			return c.convertVariableSection(node.Tok, node.Specs)
		case token.TYPE:
			panic("not implemented")
		default:
			panic("bad declaration token")
		}
	case *ast.FuncDecl:
		panic("not implemented")
	default:
		ast.Print(c.Fset, decl)
		panic("bad declaration")
	}
}

func (c *Context) convertVariableSection(sectionToken token.Token, specs []ast.Spec) string {
	result := ""
	if sectionToken == token.VAR {
		result += "var "
	} else {
		result += "const "
	}
	parts := []string{}
	for _, gspec := range specs {
		spec := gspec.(*ast.ValueSpec)
		for i, name := range spec.Names {
			var value ast.Expr
			if spec.Values == nil {
				value = nil
			} else {
				value = spec.Values[i]
			}
			part := c.convertFuncName(name.Name)
			if spec.Type != nil {
				part += ": " + c.convertType(spec.Type)
			}
			if value != nil {
				part += " = " + c.convertExpr(value)
			}
			parts = append(parts, part)
		}
	}
	if len(parts) == 1 {
		return result + parts[0]
	} else {
		return result + "\n" + indent(strings.Join(parts, "\n"))
	}
}
