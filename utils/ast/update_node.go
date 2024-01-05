package main

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 遍历节点
func main() {
	fileFullPath, err := filepath.Abs("utils\\ast\\query\\card.gen.go")
	if err != nil {
		panic(errors.WithMessage(err, "获取文件路径失败"))
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileFullPath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	const (
		suffix = "DO"
	)
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:

		case *ast.File:
			for i,v := range node.Decls {

				f,ok := v.(*ast.FuncDecl)
				if !ok{
					continue
				}
				newFuncDecl := *f
				//newSpec :=&ast.FuncDecl{
				//	Doc:  nil,
				//	Recv: nil,
				//	Name: nil,
				//	Type: nil,
				//	Body: nil,
				//}
				log.Println(i)
				if strings.HasSuffix(f.Name.Name, "DO") {
					newNameIdent :=*f.Name
					newNameIdent.Name = "test1"
					newFuncDecl.Name = &newNameIdent
					node.Decls = append(node.Decls, &newFuncDecl)
					return false
				}
			}
			return false
		case *ast.GenDecl:
			color.Red("GenDecl node %+v", node)
		case *ast.TypeSpec: //类型定义 ((重点)) 当我们要解析一个结构体用到
			//获取所有的结构体类型
			color.Yellow("TypeSpec node %+v", node)
			if node.Name == nil {
				log.Panicf("Name is nil")
			}
			if !strings.HasSuffix(node.Name.Name, suffix) {
				return true
			}

			_, ok := node.Type.(*ast.StructType)
			if !ok {
				return true
			}


		}
		return true
	})
	base := filepath.Base(fileFullPath)
	dir := filepath.Dir(fileFullPath)
	p := filepath.Join(dir,	strings.TrimSuffix(base,".go")+".new.go")
	fs, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err!=nil{
		panic(err)
	}

	err = format.Node(fs, fset, f)
}
