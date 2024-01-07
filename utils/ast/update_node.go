package main

import (
	"fmt"
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

type structDesc struct {
	Name   string
	Fields []fieldDesc
}
type fieldDesc struct {
	FieldName string
	FieldType string
}

var (
	_typeMap = map[string]string{
		"field.Int8":   "int8",
		"field.Uint8":  "uint8",
		"field.Uint16": "uint16",
		"field.int16":  "int16",
		"field.Int32":  "int32",
		"field.Uint32": "uint32",
		"field.Int64":  "int64",
		"field.UInt64": "uint64",
		"field.String": "string",
		"field.Bool":   "bool",
	}
)

// 给gorm gen 生成的代码增加一些我们自己的方法。
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
		suffix = "Do"
		all    = "ALL"
	)

	var (
		s           structDesc
		recv        *ast.FieldList
		whereIdent  = ast.NewIdent("Where")
		eqIdent     = ast.NewIdent("Eq")
		callerIdent = ast.NewIdent("caller")
		takeIdent   = ast.NewIdent("Take")
		paramIdent  = ast.NewIdent("param")
	)

	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.TypeSpec:

			name := node.Name.Name
			if strings.HasSuffix(name, suffix) {
				return true
			}
			s.Name = name
			structType, ok := node.Type.(*ast.StructType)
			if !ok {
				return true
			}
			fd := fieldDesc{}
			for _, v := range structType.Fields.List {
				if len(v.Names) == 0 {
					continue
				}
				fd.FieldName = v.Names[0].Name
				if fd.FieldName == all {
					continue
				}
				selectorExpr, ok := v.Type.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				i, ok := selectorExpr.X.(*ast.Ident)
				if !ok {
					continue
				}

				fd.FieldType = fmt.Sprintf("%s.%s", i.Name, selectorExpr.Sel.Name)
				s.Fields = append(s.Fields, fd)
			}
		case *ast.FuncDecl:

			if node.Recv == nil || recv != nil {
				return true
			}

			for _, v := range node.Recv.List {
				identNode, ok := v.Type.(*ast.Ident)
				if !ok {
					return true
				}
				if strings.HasSuffix(identNode.Name, suffix) {
					return true
				}

				if len(node.Recv.List) == 0 {
					log.Panicf("invalid recv !!")
				}
				//node.Recv.List[0].Type.(*ast.Ident).Name = "*" + node.Recv.List[0].Type.(*ast.Ident).Name
				recv = node.Recv
				callerIdent = node.Recv.List[0].Names[0]

			}
		}

		return true
	})

	var (
		s1 = &ast.SelectorExpr{
			X:   callerIdent,                   //c
			Sel: ast.NewIdent(s.Name + suffix), //cardDo
		}

		whereFunc = &ast.SelectorExpr{
			X:   s1,
			Sel: whereIdent,
		}
		withContext = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   callerIdent,
				Sel: ast.NewIdent("WithContext"),
			},
			Lparen:   0,
			Args:     []ast.Expr{ast.NewIdent("ctx")},
			Ellipsis: 0,
			Rparen:   0,
		}
		withContextWhere = &ast.SelectorExpr{
			X:   withContext,
			Sel: whereIdent,
		}
	)

	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.File:
			for _, v := range s.Fields {

				eqFunc := &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   callerIdent,
						Sel: ast.NewIdent(v.FieldName),
					},
					Sel: eqIdent,
				}

				funcName := fmt.Sprintf("Find%sBy%s", firstUpper(s.Name), v.FieldName)
				funcNameCtx := fmt.Sprintf("Find%sBy%sCtx", firstUpper(s.Name), v.FieldName)
				paramIdent = ast.NewIdent(v.FieldName)
				paramIdent.Obj = &ast.Object{
					Kind: ast.Var,
					Name: v.FieldName,
					Decl: &ast.Field{
						Doc:     nil,
						Names:   []*ast.Ident{ast.NewIdent(v.FieldName)},
						Type:    ast.NewIdent(_typeMap[v.FieldType]),
						Tag:     nil,
						Comment: nil,
					},
					Data: nil,
					Type: nil,
				}
				//新增func节点
				newSpec := &ast.FuncDecl{
					Doc:  nil,
					Recv: recv,
					Name: ast.NewIdent(funcName),
					Type: &ast.FuncType{
						Func:       0,
						TypeParams: nil,
						//参数
						Params: &ast.FieldList{
							Opening: 0,
							List: []*ast.Field{{
								Doc:     nil,
								Names:   []*ast.Ident{paramIdent},
								Type:    ast.NewIdent(_typeMap[v.FieldType]),
								Tag:     nil,
								Comment: nil,
							}},
							Closing: 0,
						},
						Results: &ast.FieldList{
							Opening: 0,
							List: []*ast.Field{
								{
									Doc:     nil,
									Names:   []*ast.Ident{ast.NewIdent("result")},
									Type:    ast.NewIdent("*model." + firstUpper(s.Name)),
									Tag:     nil,
									Comment: nil,
								}, {
									Doc:     nil,
									Names:   []*ast.Ident{ast.NewIdent("err")},
									Type:    ast.NewIdent("error"),
									Tag:     nil,
									Comment: nil,
								}},
							Closing: 0,
						},
					},
					Body: &ast.BlockStmt{
						Lbrace: 0,
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Return: 0,
								Results: []ast.Expr{&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun:    whereFunc,
											Lparen: 0,
											Args: []ast.Expr{&ast.CallExpr{
												Fun:      eqFunc,
												Lparen:   0,
												Args:     []ast.Expr{paramIdent},
												Ellipsis: 0,
												Rparen:   0,
											}},
											Ellipsis: 0,
											Rparen:   0,
										},
										Sel: takeIdent,
									},
									Lparen:   0,
									Args:     nil,
									Ellipsis: 0,
									Rparen:   0,
								}},
							},
						},
						Rbrace: 0,
					},
				}
				newSpecCtx := &ast.FuncDecl{
					Doc:  nil,
					Recv: recv,
					Name: ast.NewIdent(funcNameCtx),
					Type: &ast.FuncType{
						Func:       0,
						TypeParams: nil,
						//参数
						Params: &ast.FieldList{
							Opening: 0,
							List: []*ast.Field{{
								Doc:     nil,
								Names:   []*ast.Ident{ast.NewIdent("ctx")},
								Type:    ast.NewIdent("context.Context"),
								Tag:     nil,
								Comment: nil,
							}, {
								Doc:     nil,
								Names:   []*ast.Ident{paramIdent},
								Type:    ast.NewIdent(_typeMap[v.FieldType]),
								Tag:     nil,
								Comment: nil,
							}},
							Closing: 0,
						},
						Results: &ast.FieldList{
							Opening: 0,
							List: []*ast.Field{
								{
									Doc:     nil,
									Names:   []*ast.Ident{ast.NewIdent("result")},
									Type:    ast.NewIdent("*model." + firstUpper(s.Name)),
									Tag:     nil,
									Comment: nil,
								}, {
									Doc:     nil,
									Names:   []*ast.Ident{ast.NewIdent("err")},
									Type:    ast.NewIdent("error"),
									Tag:     nil,
									Comment: nil,
								}},
							Closing: 0,
						},
					},
					Body: &ast.BlockStmt{
						Lbrace: 0,
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Return: 0,
								Results: []ast.Expr{&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun:    withContextWhere,
											Lparen: 0,
											Args: []ast.Expr{&ast.CallExpr{
												Fun:      eqFunc,
												Lparen:   0,
												Args:     []ast.Expr{paramIdent},
												Ellipsis: 0,
												Rparen:   0,
											}},
											Ellipsis: 0,
											Rparen:   0,
										},
										Sel: takeIdent,
									},
									Lparen:   0,
									Args:     nil,
									Ellipsis: 0,
									Rparen:   0,
								}},
							},
						},
						Rbrace: 0,
					},
				}
				node.Decls = append(node.Decls, newSpec, newSpecCtx)
			}

			return false

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

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
