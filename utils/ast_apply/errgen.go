package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	statusPkg = `"google.golang.org/grpc/status"`
)

// 根据错误码生成对于的错误。
func main1() {
	codePath := flag.String("p", "", "错误码路径")
	flag.Parse()
	pathList, err := filepath.Glob(*codePath)
	if err != nil {
		log.Panicf("文件路径错误:%v", err)
	}
	if len(pathList) == 0 {
		log.Printf("该目录下无匹配的文件%v", *codePath)
		return
	}
	const (
		filedSuffix = "Code"
	)

	for _, path := range pathList {
		fileFullPath, err := filepath.Abs(path)
		if err != nil {
			log.Panicf("文件路径错误:%v path=%v", err, fileFullPath)
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, fileFullPath, nil, parser.ParseComments)
		if err != nil {
			log.Panicf("解析文件失败:%v", err)
		}
		var (
			varName string

			varDecl = &ast.GenDecl{
				Doc:    nil,
				TokPos: 0,
				Tok:    token.VAR,
				Lparen: 0,
				Specs:  nil,
				Rparen: 0,
			}
			hasVar = false
		)

		oldDecl := make([]ast.Decl, 0, 10)
		for _, v := range f.Decls {

			decl, ok := v.(*ast.GenDecl)
			if !ok {
				continue
			}
			oldDecl = append(oldDecl, decl)
			if decl.Tok == token.VAR {
				hasVar = true
			}
			if hasVar {
				return
			}
			if decl.Tok != token.CONST {
				continue
			}

			for _, v := range decl.Specs {
				var comment = ``
				valueSpec := v.(*ast.ValueSpec)
				varName = strings.TrimSuffix(valueSpec.Names[0].Name, filedSuffix)
				codeName := valueSpec.Names[0].Name

				comments := valueSpec.Names[0].Obj.Decl.(*ast.ValueSpec).Comment
				if comments != nil {
					comment = strings.TrimPrefix(comments.List[0].Text, "//")
				}
				varValueSpec := &ast.ValueSpec{}

				varIdent := ast.NewIdent(varName)
				varIdent.Obj = &ast.Object{
					Kind: ast.Var,
					Name: varName,
					Decl: &ast.ValueSpec{
						Doc:   nil,
						Names: []*ast.Ident{ast.NewIdent(varName)},
						Type:  nil,
						Values: []ast.Expr{&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent(varName),
								Sel: ast.NewIdent("Error"),
							},
							Lparen: 0,
							Args: []ast.Expr{
								&ast.BasicLit{
									ValuePos: 0,
									Kind:     token.STRING,
									Value:    comment,
								},
							},
							Ellipsis: 0,
							Rparen:   0,
						}},
						Comment: nil,
					},
					Data: nil,
					Type: nil,
				}
				varValueSpec.Values = []ast.Expr{&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent(codeName),
						Sel: ast.NewIdent("Error"),
					},
					Lparen: 0,
					Args: []ast.Expr{
						&ast.BasicLit{
							ValuePos: 0,
							Kind:     token.STRING,
							Value:    fmt.Sprintf("\"%v\"", comment),
						},
					},
					Ellipsis: 0,
					Rparen:   0,
				}}
				varValueSpec.Names = append(varValueSpec.Names, varIdent)
				varDecl.Specs = append(varDecl.Specs, varValueSpec)
			}
			oldDecl = append(oldDecl, varDecl)

		}
		f.Decls = oldDecl
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
