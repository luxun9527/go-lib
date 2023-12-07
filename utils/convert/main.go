package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
)

const (
	inputPath = "E:\\demoproject\\go-lib\\utils\\convert\\card.gen.go"
	src       = "E:\\demoproject\\go-lib\\utils\\convert\\card.proto"
)

// type goType int32
var (
	typeMap = map[string]string{
		"int8":   "int32",
		"uint8":  "int32",
		"int16":  "int32",
		"uint16": "uint32",
		"int32":  "int32",
		"uint32": "uint32",
		"int64":  "int64",
		"uint64": "uint64",
		"string": "string",
		"[]byte": "bytes",
		"bool":   "bool",
	}
)

func main() {

	fset := token.NewFileSet()
	path, _ := filepath.Abs("./utils/convert/card.gen.go")
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Printf("init parse file failed %v", err)
		return
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		var ts *ast.TypeSpec
		for _, spec := range genDecl.Specs {
			if t, ok := spec.(*ast.TypeSpec); ok {
				ts = t
				break
			}

		}

		if ts == nil {
			continue
		}
		structDecl, ok := ts.Type.(*ast.StructType)
		if !ok {
			continue
		}
		for _, field := range structDecl.Fields.List {
			//	fieldName := field.Type.(*ast.Ident).MessageName
			if field.Doc != nil {
				log.Printf("doc start")
				for _, comment := range field.Comment.List {
					log.Printf("field doc %v", comment.Text)
				}
			}
			if field.Comment != nil {
				log.Printf("comment start")
				for _, comment := range field.Comment.List {
					log.Printf("field Comment %v", comment.Text)
				}
			}
		}
	}
	//ast.ObjKind()

}
