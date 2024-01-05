package main

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

// 遍历节点
func main() {
	fileFullPath, err := filepath.Abs("utils\\ast\\card.gen.go")
	if err!=nil{
		panic(errors.WithMessage(err,"获取文件路径失败"))
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileFullPath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			color.Red("GenDecl node %+v",node)
		case  *ast.TypeSpec: //类型定义 ((重点)) 当我们要解析一个结构体用到 
			color.Yellow("TypeSpec node %+v",node)
		case *ast.StructType:
			color.Yellow("StructType node %+v", node)
		case *ast.Field:
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
			color.Black("Comment node %+v",node)
		case *ast.CommentGroup:
			color.Magenta("commentGroup node %+v",node)
		case *ast.InterfaceType:
			color.Magenta("interfaceType node %+v",node)
		case *ast.ValueSpec: //((重点))当我们要解析全局变量，常量的时候用到
			color.Magenta("ValueSpec node %+v",node)

		}
		return true
	})
}
