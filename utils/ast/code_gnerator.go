package main

import (
	"github.com/spf13/cast"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type constantDesc struct {
	Name    string
	Value   int32
	Comment string
}

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "E:\\demoproject\\go-lib\\utils\\ast\\code.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	const filedSuffix = "Code"
	var (
		currentVal = int32(0)
	)
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ValueSpec:
			d := &constantDesc{}
			names := node.Names
			if len(names) == 0 {
				return true
			}
			//name
			name := names[0].Name
			if !strings.HasSuffix(name, filedSuffix) {
				return true
			}
			name = strings.TrimSuffix(name, filedSuffix)
			d.Name = name
			//value
			for _, v := range node.Values {
				e, ok := v.(*ast.BinaryExpr)
				if ok {
					l, ok := e.Y.(*ast.BasicLit)
					if ok {
						d.Value = cast.ToInt32(l.Value)
						currentVal = d.Value
					}
				}
			}
			if d.Value == 0 {
				currentVal++
				d.Value = currentVal
			}
			//comment
			var comment string
			if node.Comment != nil {
				comment = node.Comment.Text()
			}
			if node.Doc != nil {
				comment = node.Doc.Text()
			}
			d.Comment = comment

			log.Printf("%+v", d)
		}
		return true
	})

}

const tmp = ``

func castToGrpcError(c *constantDesc) string {

	return ""
}
