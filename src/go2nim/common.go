package go2nim

import "go/ast"
import "go/token"

type ListIdentifiers struct {
	idents []string
}

func (l *ListIdentifiers) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.Ident:
		l.idents = append(l.idents, node.Name)
	}
	return l
}

func (c *Context) getIdentifiers(node ast.Node) []string {
	visitor := &ListIdentifiers{idents: []string{}}
	ast.Walk(visitor, node)
	return visitor.idents
}

func (c *Context) getAssignedIdentifiers(decl ast.Decl) []string {
	result := []string{}
	gdecl := decl.(*ast.GenDecl)
	for _, gspec := range gdecl.Specs {
		spec := gspec.(*ast.ValueSpec)
		for _, name := range spec.Names {
			result = append(result, name.Name)
		}
	}
	return result
}


type GotoVisitor struct {
	containsGoto bool
	lookFor string
}

func (l *GotoVisitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.BranchStmt:
		if node.Tok == token.GOTO && (node.Label.Name == l.lookFor || "" == l.lookFor) {
			l.containsGoto = true
		}
	case *ast.FuncLit:
		return nil
	}
	return l
}

func (c *Context) containsGoto(node ast.Node, label string) bool {
	l := &GotoVisitor{lookFor: label, containsGoto: false}
	ast.Walk(l, node)
	return l.containsGoto
}
