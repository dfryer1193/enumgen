package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
)

type Generator struct {
	buf bytes.Buffer
	pkg *Package
}

func (g *Generator) Printf(format string, args ...any) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) generate(typeName string, values []Value) {
	g.Printf("\n")
	baseTypeName := values[0].info.Name() // All the values are the same type
	g.Printf("var _%sValues = map[%s]%s{", typeName, baseTypeName, typeName)
	for _, v := range values {
		actualVal := v.value.ExactString()
		g.Printf("\n\t%s: %s,", actualVal, v.name)
	}
	g.Printf("\n}")
	g.Printf("\n")
	g.Printf("\nfunc Get%s(x %s) (%s, bool) {", typeName, baseTypeName, typeName)
	g.Printf("\n\tval, ok := _%sValues[x]", typeName)
	g.Printf("\n\treturn val, ok")
	g.Printf("\n}\n")
}

func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}

	return src
}
