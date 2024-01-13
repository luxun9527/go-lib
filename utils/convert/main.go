package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

const (
	inputPath = "E:\\demoproject\\go-lib\\utils\\convert\\card.gen.go"
	src       = "E:\\demoproject\\go-lib\\utils\\convert\\card.proto"
)

type fieldNameMode int8

const (
	snakeCase fieldNameMode = iota + 1
	upperCamelCase
	lowerCamelCase
)

var (
	fn      = flag.Int("fieldNameMode", 1, "")
	addHttp = flag.Bool("addHttp", true, "")
)

// type goType int32
var (
	_typeMap = map[string]string{
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

type Data struct {
	Msg     []*Message
	AddHttp bool
}
type Message struct {
	MessageName string
	Comment     string
	Fields      []*Field
}
type Field struct {
	FieldName string
	FieldType string
	Comment   string
}

func main() {
	codePath := flag.String("p", "", "path")
	flag.Parse()
	pathList, err := filepath.Glob(*codePath)
	if err!=nil{
		log.Panicf("err %v", err)
	}
	for _,v := range pathList {
		fset := token.NewFileSet()
		path, _ := filepath.Abs(v)
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Printf("init parse file failed %v", err)
			return
		}
		var d = Data{
			Msg:     nil,
			AddHttp: *addHttp,
		}
		messages := make([]*Message, 0, 10)
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
			msg := &Message{
				MessageName: ts.Name.Name,
				Comment:     "",
				Fields:      nil,
			}
			fields := make([]*Field, 0, 10)
			for _, field := range structDecl.Fields.List {
				fieldName := field.Names[0].Name
				switch fieldNameMode(*fn) {
				case upperCamelCase:

				case lowerCamelCase:
					fieldName = toLowerCase(fieldName)
				case snakeCase:
					fieldName = camelToSnake(fieldName)
				}

				fieldType := field.Type.(*ast.Ident).Name
				fieldType = _typeMap[fieldType]
				var comment string
				//字段右方的注释
				if field.Comment != nil {
					for _, c := range field.Comment.List {
						comment += c.Text
					}
				}
				f := &Field{
					FieldName: fieldName,
					FieldType: fieldType,
					Comment:   comment,
				}
				fields = append(fields, f)
			}
			msg.Fields = fields
			messages = append(messages, msg)
		}
		d.Msg = messages
		for _, v := range d.Msg {
			log.Printf("message %v", v.MessageName)
			for _, v := range v.Fields {
				log.Printf("fieldname=%v,fieldtype=%v,comment=%v", v.FieldName, v.FieldType, v.Comment)
			}
		}
		p, err := template.New("proto.tpl").
			Funcs(template.FuncMap{
				"add": func(x, y int) int {
					return x + y
				},
			}).ParseFiles("E:\\demoproject\\go-lib\\utils\\convert\\proto.tpl")

		if err != nil {
			log.Panicf("parse failed %v", err)
		}
		fs, err := os.OpenFile("test.proto", os.O_CREATE|os.O_RDWR, 0644)

		if err := p.Execute(fs, d); err != nil {
			log.Panicf("Execute failed %v", err)
			return
		}
	}

}

func parseTag(tag string, target string) string {
	tags := strings.Split(tag, " ")
	for _, v := range tags {
		d := strings.Split(v, ":")
		if len(d) == 2 {
			if d[0] == target {
				return strings.Trim(d[1], "\"`")
			}
		}
	}
	return ""
}
func camelToSnake(input string) string {
	if len(input) == 0 {
		return input
	}

	fieldName := []rune(input)
	var result []rune

	for i := 0; i < len(fieldName); i++ {
		c := fieldName[i]
		// 处理ID字符
		if c == 'I' && i+1 < len(fieldName) && fieldName[i+1] == 'D' {
			if len(fieldName) == 2 {
				result = append(result, []rune("id")...)
			} else {
				result = append(result, []rune("_id")...)
			}
			i++
			continue
		}

		// 处理其他大写字母
		if i > 0 && unicode.IsUpper(c) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(c))

	}

	return string(result)
}

func toLowerCase(s string) string {
	if len(s) == 0 {
		return s // 空字符串无需处理
	}

	// 将字符串转为 rune 切片，以处理 Unicode 字符
	runes := []rune(s)

	// 将第一个字符转为小写
	runes[0] = unicode.ToLower(runes[0])

	// 返回转换后的字符串
	return string(runes)
}
