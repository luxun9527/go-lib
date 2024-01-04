package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "D:\\project\\go-lib\\utils\\ast\\code.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.TypeSpec:
			log.Println(node)
		case *ast.StructType:

		}
		return true
	})

}
