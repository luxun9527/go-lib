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
	"strings"
)

// 给grom gen 增加RawWhere RawSelect方法 addgen.exe -p
func main() {
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
			panic(err)
		}
		const (
			suffix    = "Do"
			RawSelect = "RawSelect"
		)

		var (
			recv *ast.FieldList
		)

		var (
			modelName   string
			firstLetter string
			hasGen      bool
		)

		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {

			case *ast.FuncDecl:

				if node.Recv == nil {
					return true
				}
				if node.Name != nil && node.Name.Name == RawSelect {
					hasGen = true
				}
				for _, v := range node.Recv.List {
					if recv != nil {
						return true
					}
					identNode, ok := v.Type.(*ast.Ident)
					if !ok {
						return true
					}
					if strings.HasSuffix(identNode.Name, suffix) {
						if len(node.Recv.List) == 0 {
							log.Panicf("invalid recv !!")
						}
						recv = node.Recv
						modelName = identNode.Name
						firstLetter = modelName[0:1]
						return true
					}

				}

			}

			return true
		})
		//已经生成过
		if hasGen {
			continue
		}
		//新增func节点
		whereSpec := &ast.FuncDecl{
			Doc:  nil,
			Recv: recv,
			Name: ast.NewIdent("RawWhere"),
			Type: &ast.FuncType{
				Func:       0,
				TypeParams: nil,
				//参数
				Params: &ast.FieldList{
					Opening: 0,
					List: []*ast.Field{{
						Doc:   nil,
						Names: []*ast.Ident{ast.NewIdent("query")},
						Type: &ast.InterfaceType{
							Interface: token.Pos(8648),
							Methods: &ast.FieldList{
								Opening: 8657,
								List:    nil,
								Closing: 8658,
							},
							Incomplete: false,
						},
						Tag:     nil,
						Comment: nil,
					}, {
						Doc:   nil,
						Names: []*ast.Ident{ast.NewIdent("args")},
						Type: &ast.Ellipsis{
							Ellipsis: 0,
							Elt: &ast.InterfaceType{
								Interface: 8648,
								Methods: &ast.FieldList{
									Opening: 8657,
									List:    nil,
									Closing: 8658,
								},
								Incomplete: false,
							},
						},
						Tag:     nil,
						Comment: nil,
					}},
					Closing: 0,
				},
				Results: &ast.FieldList{
					Opening: 0,
					List: []*ast.Field{
						{
							Doc: nil,
							Type: &ast.StarExpr{
								Star: 0,
								X:    ast.NewIdent(modelName),
							},
							Tag:     nil,
							Comment: nil,
						}},
					Closing: 0,
				},
			},
			Body: &ast.BlockStmt{
				Lbrace: 0,
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs:    []ast.Expr{ast.NewIdent("db")},
						TokPos: 0,
						Tok:    token.DEFINE,
						Rhs: []ast.Expr{&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												NamePos: 0,
												Name:    firstLetter,
												Obj:     ast.NewObj(ast.Var, firstLetter),
											},
											Sel: ast.NewIdent("DO"),
										},
										Sel: ast.NewIdent("UnderlyingDB"),
									},
									Lparen:   0,
									Args:     nil,
									Ellipsis: 0,
									Rparen:   0,
								},
								Sel: ast.NewIdent("Where"),
							},
							Lparen:   0,
							Args:     []ast.Expr{ast.NewIdent("query"), ast.NewIdent("args")},
							Ellipsis: 0,
							Rparen:   0,
						}},
					},
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(firstLetter),
							Sel: ast.NewIdent("ReplaceDB"),
						},
						Lparen:   0,
						Args:     []ast.Expr{ast.NewIdent("db")},
						Ellipsis: 0,
						Rparen:   0,
					}},
					&ast.ReturnStmt{
						Return: 0,
						Results: []ast.Expr{&ast.UnaryExpr{
							OpPos: 0,
							Op:    token.AND,
							X:     ast.NewIdent(firstLetter),
						}},
					},
				},
				Rbrace: 0,
			},
		}
		selectSpec := &ast.FuncDecl{
			Doc:  nil,
			Recv: recv,
			Name: ast.NewIdent("RawSelect"),
			Type: &ast.FuncType{
				Func:       0,
				TypeParams: nil,
				//参数
				Params: &ast.FieldList{
					Opening: 0,
					List: []*ast.Field{{
						Doc:   nil,
						Names: []*ast.Ident{ast.NewIdent("query")},
						Type: &ast.InterfaceType{
							Interface: 0,
							Methods: &ast.FieldList{
								Opening: 8657,
								List:    nil,
								Closing: 8658,
							},
							Incomplete: false,
						},
						Tag:     nil,
						Comment: nil,
					}, {
						Doc:   nil,
						Names: []*ast.Ident{ast.NewIdent("args")},
						Type: &ast.Ellipsis{
							Ellipsis: 0,
							Elt: &ast.InterfaceType{
								Interface: 0,
								Methods: &ast.FieldList{
									Opening: 8657,
									List:    nil,
									Closing: 8658,
								},
								Incomplete: false,
							},
						},
						Tag:     nil,
						Comment: nil,
					}},
					Closing: 0,
				},
				Results: &ast.FieldList{
					Opening: 0,
					List: []*ast.Field{
						{
							Doc: nil,
							Type: &ast.StarExpr{
								Star: 0,
								X:    ast.NewIdent(modelName),
							},
							Tag:     nil,
							Comment: nil,
						}},
					Closing: 0,
				},
			},
			Body: &ast.BlockStmt{
				Lbrace: 0,
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs:    []ast.Expr{ast.NewIdent("db")},
						TokPos: 0,
						Tok:    token.DEFINE,
						Rhs: []ast.Expr{&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent(firstLetter),
											Sel: ast.NewIdent("DO"),
										},
										Sel: ast.NewIdent("UnderlyingDB"),
									},
									Lparen:   0,
									Args:     nil,
									Ellipsis: 0,
									Rparen:   0,
								},
								Sel: ast.NewIdent("Select"),
							},
							Lparen:   0,
							Args:     []ast.Expr{ast.NewIdent("query"), ast.NewIdent("args")},
							Ellipsis: 0,
							Rparen:   0,
						}},
					},
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(firstLetter),
							Sel: ast.NewIdent("ReplaceDB"),
						},
						Lparen:   0,
						Args:     []ast.Expr{ast.NewIdent("db")},
						Ellipsis: 0,
						Rparen:   0,
					}},
					&ast.ReturnStmt{
						Return: 0,
						Results: []ast.Expr{&ast.UnaryExpr{
							OpPos: 0,
							Op:    token.AND,
							X:     ast.NewIdent(firstLetter),
						}},
					},
				},
				Rbrace: 0,
			},
		}

		f.Decls = append(f.Decls, whereSpec)
		f.Decls = append(f.Decls, selectSpec)

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
