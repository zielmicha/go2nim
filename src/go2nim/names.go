package go2nim

import "strings"

func (c *Context) downcaseFirstLetter(name string) string {
	firstLetter := strings.ToLower(string(name[0]))
	if firstLetter != string(name[0]) {
		if len(name) > 1 {
			secondLetter := strings.ToLower(string(name[1]))
			if secondLetter != string(name[1]) {
				return name
			}
		}
		return firstLetter + name[1:]
	} else {
		return name
	}
}

func (c *Context) upcaseFirstLetter(name string) string {
	firstLetter := strings.ToUpper(string(name[0]))
	if firstLetter != string(name[0]) {
		return firstLetter + name[1:]
	} else {
		return name
	}
}


func (c *Context) isPublic(name string) bool {
	firstLetter := strings.ToLower(string(name[0]))
	return firstLetter != string(name[0])
}

func (c *Context) quoteLabel(name string) string {
	switch name {
	case "is", "in", "addr", "and", "as", "asm", "atomic", "bind", "block", "break", "case", "cast", "concept", "const", "continue", "converter", "defer", "discard", "distinct", "div", "do", "elif", "else", "end", "enum", "except", "export", "finally", "for", "from", "func", "generic", "if", "import", "include", "interface", "isnot", "iterator", "let", "macro", "method", "mixin", "mod", "nil", "not", "notin", "object", "of", "or", "out", "proc", "ptr", "raise", "ref", "return", "shl", "shr", "static", "template", "try", "tuple", "type", "using", "var", "when", "while", "with", "without", "xor", "yield":
		return name + "Kw"
	default:
		return name
	}
}

func (c *Context) quoteKeywords(name string) string {
	if name == "_" {
		return "_"
	}
	if name[0] == '_' {
		return "underscore" + name[1:]
	}
	if name[len(name) - 1] == '_' {
		return name[:len(name) - 1] + "underscore"
	}

	switch name {
	case "is", "in":
		// Nim is very confused if we override these
		return name + "A"
	case "addr", "and", "as", "asm", "atomic", "bind", "block", "break", "case", "cast", "concept", "const", "continue", "converter", "defer", "discard", "distinct", "div", "do", "elif", "else", "end", "enum", "except", "export", "finally", "for", "from", "func", "generic", "if", "import", "include", "interface", "isnot", "iterator", "let", "macro", "method", "mixin", "mod", "nil", "not", "notin", "object", "of", "or", "out", "proc", "ptr", "raise", "ref", "return", "shl", "shr", "static", "template", "try", "tuple", "type", "using", "var", "when", "while", "with", "without", "xor", "yield":
		return "`" + name + "`"
	case "pP": // special case for stdlib
		return "`p^P`"
	default:
		return name
	}
}

func (c *Context) convertGenericName(name string) string {
	if !c.isPublic(name) {
		if _, ok := c.publicNames[name]; ok {
			return c.quoteKeywords(c.downcaseFirstLetter(name) + "Internal")
		}
		return c.quoteKeywords(c.downcaseFirstLetter(name))
	}
	return c.quoteKeywords(c.downcaseFirstLetter(name))
}

func (c *Context) convertFuncName(name string) string {
	switch name {
	case "make", "copy", "len":
		return name
	case "String", "string":
		return "`$`"
	}
	return c.convertGenericName(name)
}

func (c *Context) convertFieldName(name string) string {
	return c.convertGenericName(name)
}

func (c *Context) convertTypeName(name string) string {
	if !c.isPublic(name) {
		name = c.upcaseFirstLetter(name)
		if _, ok := c.publicNames[name]; ok {
			return c.quoteKeywords(name + "Internal")
		}
	}
	return c.quoteKeywords(c.upcaseFirstLetter(name))
}
