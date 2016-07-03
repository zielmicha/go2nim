package go2nim

import "strings"
import "go/ast"

func (c *Context) convertInlineStruct(node *ast.StructType) string {
	fields := []string{}
	for _, field := range node.Fields.List {
		for _, name := range field.Names {
			fields = append(fields, c.convertFieldName(name.Name) + ": " + c.convertType(field.Type))
		}
	}
	return "struct({" + strings.Join(fields, ", ") + "})"
}

func (c *Context) convertInlineInterface(node *ast.InterfaceType) string {
	fields := make([]string, len(node.Methods.List))
	for i, method := range node.Methods.List {
		name := method.Names[0].Name
		var typeBase interface{} = method.Type
		fields[i] = c.convertFieldName(name) + "(" + c.convertParamList(typeBase.(*ast.FuncType).Params, nil) + "): " + c.convertReturnType(typeBase.(*ast.FuncType).Results)
	}
	return "interface(" + strings.Join(fields, ", ") + ")"
}

func (c *Context) convertType(expr ast.Expr) string {
	var node interface{} = expr
	switch node := node.(type) {
	case *ast.Ident:
		// TODO: translate builtin identifiers
		return node.Name
	case *ast.FuncType:
		return "(proc(" + c.convertParamList(node.Params, nil) + "): " + c.convertReturnType(node.Results) + ")"
	case *ast.StructType:
		return c.convertInlineStruct(node)
	case *ast.InterfaceType:
		return c.convertInlineInterface(node)
	case *ast.MapType:
		return "Map[" + c.convertType(node.Key) + ", " + c.convertType(node.Value) + "]"
	case *ast.ArrayType:
		if node.Len == nil {
			return "GoSlice[" + c.convertType(node.Elt) + "]"
		} else {
			panic("todo: array type")
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
		return "varargs[" + c.convertType(node.Elt) + "]"
	case *ast.SelectorExpr:
		return c.convertType(node.X) + "." + node.Sel.Name
	default:
		ast.Print(c.Fset, expr)
		panic("bad type expression")
	}
}

func (c *Context) convertParamList(fields *ast.FieldList, recv *ast.FieldList) string {
	if fields == nil {
		return ""
	}

	result := []string{}
	for _, field := range fields.List {
		for _, name := range field.Names {
			result = append(result, name.Name + ": " + c.convertType(field.Type))
		}
	}

	if recv != nil {
		receiver := recv.List[0]
		recvParam := receiver.Names[0].Name + ": " + c.convertType(receiver.Type)
		result = append([]string{recvParam}, result...)
	}

	return strings.Join(result, ", ")
}

func (c *Context) convertReturnType(fields *ast.FieldList) string {
	if fields == nil || len(fields.List) == 0 {
		return "void"
	} else if len(fields.List) == 1 {
		return c.convertType(fields.List[0].Type)
	} else {
		return "composite return type"
	}
}
