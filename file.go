package main

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log"
)

type Package struct {
	name         string
	defs         map[*ast.Ident]types.Object
	files        []*File
	hasTestFiles bool
}

type File struct {
	pkg  *Package
	file *ast.File

	typeName string
	values   []Value
}

type Value struct {
	name  string
	info  types.Basic
	value constant.Value
}

func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		return true
	}

	typ := ""
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec)
		if vspec.Type == nil && len(vspec.Values) > 0 {
			typ = ""

			ce, ok := vspec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}

			id, ok := ce.Fun.(*ast.Ident)
			if !ok {
				continue
			}

			typ = id.Name
		}

		if vspec.Type != nil {
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}

		if typ != f.typeName {
			continue
		}

		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}

			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("no value for constant %s", name)
			}

			info := obj.Type().Underlying().(*types.Basic)
			value := obj.(*types.Const).Val()
			v := Value{
				name:  name.Name,
				info:  *info,
				value: value,
			}
			f.values = append(f.values, v)
		}
	}

	return false
}
