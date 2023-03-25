//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"regexp"
	"text/template"

	"github.com/VKCOM/php-parser/pkg/ast"
)

type templData struct {
	Types []nodeType
}

type nodeType struct {
	Name     string
	FuncName string
}

var fileTempl = template.Must(
	template.New("").
		Parse(`// Code generated by "go generate go run node_funcs_gen.go"; DO NOT EDIT.

package ast

import "github.com/VKCOM/php-parser/pkg/position"
{{range $typ := .Types}}
var _ Vertex = &{{$typ.Name}}{}

func (n *{{$typ.Name}}) Accept(v Visitor) {
    v.{{$typ.FuncName}}(n)
}

func (n *{{$typ.Name}}) GetPosition() *position.Position {
    return n.Position
}
{{end}}`),
)

var typesRgx = regexp.MustCompile(`type ([a-zA-Z]+) struct`)

func main() {
	content, err := os.ReadFile("node.go")
	if err != nil {
		panic(fmt.Errorf("reading 'ast.go': %w", err))
	}

	matches := typesRgx.FindAllSubmatch(content, -1)
	types := make([]nodeType, 0, len(matches))
	for _, match := range matches {
		name := string(match[1])
		funcName := name
		if overwrite, ok := ast.TypeToVisitorNameMap[name]; ok {
			funcName = overwrite
		}

		types = append(types, nodeType{
			Name:     name,
			FuncName: funcName,
		})
	}

	f, err := os.Create("node_funcs.go")
	if err != nil {
		panic(fmt.Errorf("creating/opening 'node_funcs.go': %w", err))
	}

	fileTempl.Execute(f, templData{types})
}
