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

func main() {
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
		const (
			filedSuffix = "Code"
		)

		var (
			newVarGenDecl = &ast.GenDecl{
				Tok: token.VAR,
			}
			newDecls =make([]ast.Decl,0,4)
			hasImportNode = false
			newImportGenSpec = &ast.GenDecl{
				Doc:    nil,
				TokPos: 30,
				Tok: token.IMPORT,
				Lparen: 0,
				Specs:  nil,
				Rparen: 0,
			}
			newImpSpec = &ast.ImportSpec{
				Doc:     nil,
				Name:    nil,
				Path:    &ast.BasicLit{
					ValuePos: 0,
					Kind:     0,
					Value:    statusPkg,
				},
				Comment: nil,
				EndPos:  0,
			}
		)
		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if genDecl.Tok == token.IMPORT {
				hasImportNode = true
				hasStatusPkg := false
				for _,v := range genDecl.Specs {
					imSpec,ok := v.(*ast.ImportSpec)
					if !ok{
						continue
					}
					if statusPkg == imSpec.Path.Value{
						hasStatusPkg = true
					}
				}
				if !hasStatusPkg{
					genDecl.Specs = append(genDecl.Specs, newImpSpec)
				}
			}

			//删除掉原来的变量node
			if genDecl.Tok == token.VAR {
				continue
			}
			//复制原来的
			newDecls = append(newDecls, decl)
			if  genDecl.Tok != token.CONST {
				continue
			}
			//2、遍历const,构建新的node
			isValidConstNode := true
			for _, sp := range genDecl.Specs {
				valueSpec, ok := sp.(*ast.ValueSpec)
				if !ok {
					continue
				}
				if len(valueSpec.Names) == 0 {
					continue
				}

				//1、取code变量名
				codeFieldName := valueSpec.Names[0].Name
				//不是以code 结尾就跳过
				if !strings.HasSuffix(codeFieldName, filedSuffix) {
					isValidConstNode = false
					break
				}

				//2、取注释
				var comment string
				if valueSpec.Comment != nil {
					comment = valueSpec.Comment.Text()
				}
				if valueSpec.Doc != nil {
					comment = valueSpec.Doc.Text()
				}
				comment = strings.TrimSpace(comment)
				//3 取变量ident
				codeIndent := valueSpec.Names[0]

				newValueSpec := &ast.ValueSpec{
					Doc:     nil,
					Names:   nil,
					Type:    nil,
					Values:  nil,
					Comment: nil,
				}
				newValueSpec.Names = append(newValueSpec.Names, ast.NewIdent(strings.TrimSuffix(codeFieldName, filedSuffix)))

				newCallExpr := &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("status"),
						Sel: ast.NewIdent("Error"),
					},
					Lparen: 0,
					Args: []ast.Expr{
						codeIndent,
						ast.NewIdent(fmt.Sprintf(`"%s"`,comment)),
					},
					Ellipsis: 0,
					Rparen:   0,
				}

				newValueSpec.Values = append(newValueSpec.Values, newCallExpr)

				newVarGenDecl.Specs = append(newVarGenDecl.Specs, newValueSpec)
			}
			if isValidConstNode{
				newDecls = append(newDecls, newVarGenDecl)
			}
		}
		f.Decls = newDecls
		if !hasImportNode{
			newImportGenSpec.Specs = append(newImportGenSpec.Specs, newImpSpec)
			// 添加import 语句
			s := make([]ast.Decl,0, len(f.Decls)+1)
			s = append(s, newImportGenSpec)
			s = append(s, f.Decls...)
			f.Decls =s
		}

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
