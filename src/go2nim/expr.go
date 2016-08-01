package go2nim

import "go/ast"
import "go/token"
import "strings"
import "strconv"

func (c *Context) convertBinaryToken(tok token.Token) string {
	switch tok {
	case token.ADD:
		return "+"
	case token.SUB:
		return "-"
	case token.MUL:
		return "*"
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
		ast.Print(c.Fset, tok)
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
	case token.XOR:
		return "not "

    default:
		ast.Print(c.Fset, tok)
		panic("unknown unary token")
	}
}

func (c *Context) convertCompositeLiteral(expr *ast.CompositeLit, isPointer bool) string {
	result := []string{}
	isKv := false
	itemSuffix := ""
	itemPrefix := ""

	switch typ := expr.Type.(type) {
	case *ast.ArrayType:
		// fill item type if needed
		for _, node := range expr.Elts {
			if sublit, ok := node.(*ast.CompositeLit); ok {
				if sublit.Type == nil {
					sublit.Type = typ.Elt
				}
			}
		}

		itemType := c.convertType(typ.Elt)
		switch itemType {
		case "float32", "float64", "int64", "uint64":
			itemPrefix = "" + itemType + "("
			itemSuffix = ")"
		}
	}

	for _, node := range expr.Elts {
		if kvexpr, ok := node.(*ast.KeyValueExpr); ok {
			isKv = true
			result = append(result, c.convertExpr(kvexpr.Key) + ": " + c.convertExpr(kvexpr.Value))
		} else {
			result = append(result, itemPrefix + c.convertExpr(node) + itemSuffix)
		}
	}

	val := strings.Join(result, ", ")
	var wrapped string
	if isKv {
		wrapped = "{" + val + "}"
	} else {
		wrapped = "[" + val + "]"
	}

	if expr.Type == nil {
		return wrapped
	}

	typeName := c.convertType(expr.Type)

	switch expr.Type.(type) {
	case *ast.MapType, *ast.ArrayType:
		if len(result) == 0 {
			return "make(" + typeName + ")"
		} else {
			return wrapped + ".make(" + typeName + ")"
		}
	default:
		if isPointer {
			typeName = "(ref " + typeName + ")"
		}
		if isKv {
			return typeName + "(" + val + ")"
		} else {
			return "make((" + val + "), " + typeName + ")"
		}
	}
}

func (c *Context) convertExprMulti(expr ast.Expr) string {
	var node interface{} = expr

	switch node := node.(type) {
	case *ast.TypeAssertExpr:
		return "maybeCastInterface(" + c.convertExpr(node.X) + ", to=" + c.convertType(node.Type) + ")"
	default:
		return c.convertExpr(expr)
	}
}

func (c *Context) convertExpr(expr ast.Expr) string {
	var node interface{} = expr

	switch node := node.(type) {
	case *ast.CallExpr:
		// One of a few places where Go is context sensitive - this may be either conversion or call
		if len(node.Args) == 1 && c.looksLikeType(node.Fun) == 1 {
			return "convert(" + c.convertType(node.Fun) + ", " + c.convertExpr(node.Args[0]) + ")"
		}

		// TODO: ellipsis
		args := make([]string, len(node.Args))
		for i, arg := range node.Args {
			args[i] = c.convertExpr(arg)
		}

		if ident, ok := node.Fun.(*ast.Ident); ok {
			// check for builtin functions
			switch ident.Name {
			case "new":
				if len(node.Args) != 1 { panic("bad number of arguments") }
				return "gcnew(" + c.convertType(node.Args[0]) + ")"
			case "make":
				args[0] = c.convertType(node.Args[0])
			}
		}

		return c.convertExpr(node.Fun) + "(" + strings.Join(args, ", ") + ")"
	case *ast.BinaryExpr:
		switch node.Op {
		case token.AND_NOT:
			return "(" + c.convertExpr(node.X) + " and (not " + c.convertExpr(node.Y) + ")" + ")"
		case token.QUO:
			return "`go/`(" + c.convertExpr(node.X) + ", " + c.convertExpr(node.Y) + ")"
		default:
			return "(" + c.convertExpr(node.X) + " " + c.convertBinaryToken(node.Op) + " " + c.convertExpr(node.Y) + ")"
		}
	case *ast.IndexExpr:
		return c.convertExpr(node.X) + "[" + c.convertExpr(node.Index) + "]"
	case *ast.KeyValueExpr:
		panic("key value expr not expected here")
	case *ast.ParenExpr:
		return "(" + c.convertExpr(node.X) + ")"
	case *ast.SelectorExpr:
		left := c.convertExpr(node.X)
		right := c.convertFieldName(node.Sel.Name)
		// TODO: configurable renames
		switch left {
		case "reflect":
			switch node.Sel.Name {
			case "Invalid", "Bool", "Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Uintptr", "Float32", "Float64", "Complex64", "Complex128", "Array", "Chan", "Func", "Interface", "Map", "Ptr", "Slice", "String", "Struct", "UnsafePointer":
				return "reflect.Kind." + c.quoteKeywords(c.downcaseFirstLetter(node.Sel.Name))
			}
		}
		return left + "." + right
	case *ast.SliceExpr:
		result := []string{c.convertExpr(node.X)}
		if node.Low != nil {
			result = append(result, "low=" + c.convertExpr(node.Low))
		}
		if node.High != nil {
			result = append(result, "high=" + c.convertExpr(node.High))
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
		switch node.Kind {
		case token.CHAR:
			if strings.HasPrefix(node.Value, "'\\u") || strings.HasPrefix(node.Value, "'\\U") {
				return "Rune(0x" + node.Value[3:len(node.Value)-1] + ")"
			}
			if node.Value == "'\\n'" {
				return "'\\L'"
			}
		case token.STRING:
			// TODO
			if strings.HasPrefix(node.Value, "`") {
				return "r\"" + node.Value[1:len(node.Value)-1] + "\""
			}
		case token.FLOAT:
			if !strings.HasPrefix(node.Value, "0x") && !strings.Contains(node.Value, ".") && !strings.Contains(node.Value, "-") && strings.Contains(node.Value, "e") {
				// FIXME: potential overflows
				split := strings.Split(node.Value, "e")
				val, err := strconv.Atoi(split[0])
				if err != nil { panic(err) }
				exp, err := strconv.Atoi(split[1])
				if err != nil { panic(err) }
				var uval uint64
				var sgn int
				if val >= 0 {
					uval = uint64(val)
					sgn = 1
				} else {
					uval = uint64(-val)
					sgn = -1
				}
				if exp == 19 && uval == 1 {
					// Nim can't represent 1e19...
					return "1000000000000000000 * 10"
				}
				if exp > 18 {
					return node.Value
				}
				for exp > 0 {
					uval *= 10
					exp --
				}
				result := strconv.FormatUint(uval, 10)
				if sgn == -1 {
					result = "-" + result
				}
				return result
			}
			return node.Value
		}

		return node.Value
	case *ast.CompositeLit:
		return c.convertCompositeLiteral(node, false)
	case *ast.FuncLit:
		paramList := c.convertParamList(node.Type.Params, nil)
		body := c.convertFuncCode(node.Type, node.Body)
		return "(proc(" + paramList + "): " + c.convertReturnType(node.Type.Results) + " =\n" + indent(body) + ")"
	case *ast.Ident:
		if node.Name == "nil" {
			return "null"
		}
		if node.Name == "iota" {
			return strconv.Itoa(c.iotaValue)
		}
		// Is this result variable?
		//ast.Print(c.Fset, node)
		if node.Obj != nil {
			if _, ok := node.Obj.Decl.(*ast.Field); ok {
				if _, ok := c.resultVariables[node.Name]; ok {
					if len(c.resultVariables) == 1 {
						return "result"
					} else {
						return "result." + c.quoteKeywords(node.Name)
					}
				}
			}
		}
		return c.convertFuncName(node.Name)
	case *ast.Ellipsis:
		//ast.Print(c.Fset, node)
		// TODO
		return "..."
	default:
		return c.convertType(expr)
	}
}
