package main

import (
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

const (
	DeletedAt  = "DeletedAt"
	ImportPath = "\"gorm.io/plugin/soft_delete\""
	Tag        = "`gorm:\"softDelete:unix\" json:\"deleted_at\"`"
)

// 使用修改ast语法树，使用soft_deleted适配gen
func main11() {
	codePath := flag.String("p", "", "要替换的路径")
	flag.Parse()
	pathList, err := filepath.Glob(*codePath)
	if err != nil {
		log.Panicf("文件路径错误:%v", err)
	}
	if len(pathList) == 0 {
		log.Printf("该目录下无匹配的文件%v", *codePath)
		return
	}
	for _, path := range pathList {
		fileFullPath, err := filepath.Abs(path)
		if err != nil {
			log.Panicf("文件路径错误:%v path=%v", err, fileFullPath)
		}
		log.Println(fileFullPath)
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fileFullPath, nil, parser.ParseComments)
		if err != nil {
			log.Panicf("解析文件失败:%v", err)
		}
		decls := make([]ast.Decl, len(f.Decls)+1)
		var hasImported bool
		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.File:
				copy(decls[1:], node.Decls)
				for _, v := range node.Decls {
					decl, ok := v.(*ast.GenDecl)
					if !ok || decl.Tok != token.IMPORT {
						continue
					}
					if len(decl.Specs) == 0 {
						continue
					}
					importSpec := decl.Specs[0].(*ast.ImportSpec)
					if importSpec.Path.Value == ImportPath {
						hasImported = true
						return true
					}
				}
				if !hasImported {
					d := &ast.GenDecl{
						Doc:    nil,
						TokPos: 0,
						Tok:    token.IMPORT,
						Lparen: 0,
						Specs: []ast.Spec{&ast.ImportSpec{
							Doc:  nil,
							Name: nil,
							Path: &ast.BasicLit{
								ValuePos: 0,
								Kind:     token.STRING,
								Value:    ImportPath,
							},
							Comment: nil,
							EndPos:  0,
						}},
						Rparen: 0,
					}
					decls[0] = d
					node.Decls = decls
				}

			case *ast.Field:
				if len(node.Names) == 0 {
					return true
				}
				if node.Names[0].Name != DeletedAt {
					return true
				}
				node.Type = &ast.SelectorExpr{
					X:   ast.NewIdent("soft_delete"),
					Sel: ast.NewIdent("DeletedAt"),
				}
				node.Tag = &ast.BasicLit{
					ValuePos: 0,
					Kind:     token.STRING,
					Value:    Tag,
				}
			}
			return true

		})

		fs, err := os.OpenFile(fileFullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Panicf("open failed node = %v", err)
		}
		defer fs.Close()
		if err = format.Node(fs, fset, f); err != nil {
			log.Panicf("store  node = %v", err)
		}
	}

}
