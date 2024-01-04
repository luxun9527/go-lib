package main

import (
	"github.com/fatih/color"
	"go/ast"
	"go/parser"
	"go/token"
)

// 遍历节点
func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "D:\\project\\go-lib\\utils\\ast\\card.gen.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			color.Red("GenDecl node %+v",node)
		case  *ast.TypeSpec: //类型定义 ((重点))
			color.Yellow("TypeSpec node %+v",node)
		case *ast.StructType:
			color.Yellow("StructType node %+v", node)
		case *ast.Field: //字段定义 ((重点))
			/*
			// Expressions and types

			// A Field represents a Field declaration list in a struct type,
			// a method list in an interface type, or a parameter/result declaration
			// in a signature.
			// Field.Names is nil for unnamed parameters (parameter lists which only contain types)
			// and embedded struct fields. In the latter case, the field name is the type name.
			type Field struct {
				Doc     *CommentGroup // associated documentation; or nil
				Names   []*Ident      // field/method/(type) parameter names; or nil
				Type    Expr          // field/method/parameter type; or nil
				Tag     *BasicLit     // field tag; or nil
				Comment *CommentGroup // line comments; or nil
			}
			*/
			color.Blue("Field node %+v",node)
		case *ast.Ident://标识符，包括类型和字段名，变量名，结构体名等
			color.Green("ident node %+v",node)
		case *ast.SwitchStmt:
			color.Magenta("SwitchStmt node %+v",node)
		case *ast.Comment:
			color.Black("comment node %+v",node)
		case *ast.CommentGroup:
			color.Magenta("commentGroup node %+v",node)
		case *ast.InterfaceType:
			color.Magenta("interfaceType node %+v",node)
		}
		return true
	})
}
