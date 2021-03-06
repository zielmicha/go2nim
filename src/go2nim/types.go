package go2nim

import "strings"
import "strconv"
import "go/ast"

func (c *Context) convertStruct(node *ast.StructType) string {
	fields := []string{}
	for _, field := range node.Fields.List {
		if len(field.Names) == 0 {
			// anonymous field, must be ident
			typeName := field.Type.(*ast.Ident).Name
			fields = append(fields, "inline " + c.convertFieldName(typeName) + ": " + c.convertType(field.Type))
		} else {
			for _, name := range field.Names {
				fields = append(fields, c.convertFieldName(name.Name) + ": " + c.convertType(field.Type))
			}
		}
	}
	return "struct((" + strings.Join(fields, ", ") + "))"
}

func (c *Context) convertInterface(node *ast.InterfaceType) string {
	fields := make([]string, len(node.Methods.List))
	// ast.Print(c.Fset, node)
	for i, method := range node.Methods.List {
		if len(method.Names) == 0 {
			fields[i] = "extends " + c.convertType(method.Type)
		} else {
			name := method.Names[0].Name
			var typeBase interface{} = method.Type
			returnc := c.copy()
			returnc.anonymousTuples = true
			fields[i] = c.convertFuncName(name) + "(" + c.convertParamList(typeBase.(*ast.FuncType).Params, nil) + "): " + returnc.convertReturnType(typeBase.(*ast.FuncType).Results)
		}
	}
	if len(fields) == 0 {
		return "EmptyInterface"
	} else {
		return "iface((" + strings.Join(fields, ", ") + "))"
	}
}

func (c *Context) convertTypeDecl(expr ast.Expr) string {
	var node interface{} = expr
	switch node := node.(type) {
	case *ast.StructType:
		return c.convertStruct(node)
	case *ast.InterfaceType:
		return c.convertInterface(node)
	default:
		return c.convertType(expr)
	}
}

func (c *Context) convertType(expr ast.Expr) string {
	var node interface{} = expr
	switch node := node.(type) {
	case *ast.ParenExpr:
		return c.convertType(node.X)
	case *ast.Ident:
		switch node.Name {
		case "bool", "byte", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "string", "float", "float64", "float32", "complex128", "complex64", "uintptr":
			return node.Name
		default:
			return c.convertTypeName(node.Name)
		}
	case *ast.FuncType:
		return "(proc(" + c.convertParamList(node.Params, nil) + "): " + c.convertReturnType(node.Results) + ")"
	case *ast.StructType:
		return c.convertStruct(node)
	case *ast.InterfaceType:
		return c.convertInterface(node)
	case *ast.MapType:
		return "Map[" + c.convertType(node.Key) + ", " + c.convertType(node.Value) + "]"
	case *ast.ArrayType:
		if node.Len == nil {
			return "GoSlice[" + c.convertType(node.Elt) + "]"
		} else {
			if _, ok := node.Len.(*ast.Ellipsis); ok {
				return "GoAutoArray[" + c.convertType(node.Elt) + "]"
			} else {
				return "GoArray[" + c.convertType(node.Elt) + ", " + c.convertExpr(node.Len) + "]"
			}
		}
	case *ast.ChanType:
		if node.Dir == ast.SEND {
			return "SyncProvider[" + c.convertType(node.Value) + "]"
		} else {
			return "SyncStream[" + c.convertType(node.Value) + "]"
		}
	case *ast.StarExpr:
		return "gcptr[" + c.convertType(node.X) + "]"
	case *ast.Ellipsis:
		return "GoVarArgs[" + c.convertType(node.Elt) + "]"
	case *ast.SelectorExpr:
		return c.convertExpr(node.X) + "." + node.Sel.Name
	default:
		ast.Print(c.Fset, expr)
		panic("bad type expression")
	}
}

func (c *Context) looksLikeType(expr ast.Node) int {
	var node interface{} = expr
	switch node := node.(type) {
	case *ast.Ident:
		switch node.Name {
		case "rune", "bool", "byte", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "string", "complex128", "complex64", "uintptr":
			return 1
		}
		if node.Obj != nil {
			if _, ok := node.Obj.Decl.(*ast.TypeSpec); ok {
				return 1
			} else {
				return -1
			}
		}
		return 0
	case *ast.ParenExpr:
		return c.looksLikeType(node.X)
	case *ast.ChanType, *ast.MapType, *ast.StructType, *ast.ArrayType, *ast.InterfaceType:
		return 1
	case *ast.StarExpr:
		r := c.looksLikeType(node.X)
		if r == 0 {
			return 1
		}
		return r
	default:
		return 0
	}
}

func (c *Context) convertParamList(fields *ast.FieldList, recv *ast.FieldList) string {
	if fields == nil {
		return ""
	}

	result := []string{}
	for i, field := range fields.List {
		if len(field.Names) == 0 {
			result = append(result, "arg" + strconv.Itoa(i) + ": " + c.convertType(field.Type))
		} else {
			for _, name := range field.Names {
				result = append(result, c.convertFuncName(name.Name) + ": " + c.convertType(field.Type))
			}
		}
	}

	if recv != nil {
		receiver := recv.List[0]
		recvName := "_"
		if len(receiver.Names) > 0 {
			recvName = receiver.Names[0].Name
		}
		recvParam := recvName + ": " + c.convertType(receiver.Type)
		result = append([]string{recvParam}, result...)
	}

	return strings.Join(result, ", ")
}

func (c *Context) convertReturnType(fields *ast.FieldList) string {
	result := []string{}

	if fields == nil || len(fields.List) == 0 {
		return "void"
	}

	for i, field := range fields.List {
		if len(field.Names) == 0 {
			if c.anonymousTuples {
				result = append(result, c.convertType(field.Type))
			} else {
				result = append(result, "arg" + strconv.Itoa(i) + ": " + c.convertType(field.Type))
			}
		}
		for _, name := range field.Names {
			if c.anonymousTuples {
				result = append(result, c.convertType(field.Type))
			} else {
				result = append(result, name.Name + ": " + c.convertType(field.Type))
			}
		}
	}

	if len(result) == 1 {
		return c.convertType(fields.List[0].Type)
	} else {
		if c.anonymousTuples {
			return "(" + strings.Join(result, ", ") + ")"
		} else {
			return "tuple[" + strings.Join(result, ", ") + "]"
		}
	}
}

func (c *Context) getReturnTypeResultVariables(fields *ast.FieldList) []string {
	result := []string{}

	if fields == nil || len(fields.List) == 0 {
		return nil
	}

	for i, field := range fields.List {
		if len(field.Names) == 0 {
			result = append(result, "arg" + strconv.Itoa(i))
		}
		for _, name := range field.Names {
			result = append(result, name.Name)
		}
	}
	if len(result) == 1 {
		return nil
	} else {
		return result
	}
}
