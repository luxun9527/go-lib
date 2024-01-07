#  AST语法树

refer

https://juejin.cn/post/6844903982683389960#heading-2



本文的代码地址

https://github.com/luxun9527/go-lib/tree/master/utils/ast  您的star就是我更新的动力

## 核心概念

go文件，可以视为一颗由方法，包，字段，变量，注释组成的语法树，可以使用go提供的api，可以获取go文件中的方法，包，字段，变量，注释等信息，

一些代码生成工具，mock，注入都是使用的类似方法如 [protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag) 

在go提供的api中所有的，包，字段，变量，注释等元素视为node，主要是三种，Expressions and type nodes, statement nodes, and declaration nodes. 

```go
// There are 3 main classes of nodes: Expressions and type nodes,
// statement nodes, and declaration nodes. The node names usually
// match the corresponding Go spec production names to which they
// correspond. The node fields correspond to the individual parts
// of the respective productions.
//
// All nodes contain position information marking the beginning of
// the corresponding source text segment; it is accessible via the
// Pos accessor method. Nodes may contain additional position info
// for language constructs where comments may be found between parts
// of the construct (typically any larger, parenthesized subpart).
// That position information is needed to properly position comments
// when printing the construct.

// All node types implement the Node interface.
type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode()
}

// All declaration nodes implement the Decl interface.
type Decl interface {
	Node
	declNode()
}

```





### 1 Expression and Type

标识符和类型。

标识符：变量名，函数名，包名，类型名等。

```go
	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
		Obj     *Object   // denoted object; or nil
	}
```

类型：结构体类型，interface类型，方法类型等。

```go

// A type is represented by a tree consisting of one
// or more of the following type-specific expression
// nodes.
type (
	// An ArrayType node represents an array or slice type.
	ArrayType struct {
		Lbrack token.Pos // position of "["
		Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
		Elt    Expr      // element type
	}

	// A StructType node represents a struct type.
	StructType struct {
		Struct     token.Pos  // position of "struct" keyword
		Fields     *FieldList // list of field declarations
		Incomplete bool       // true if (source) fields are missing in the Fields list
	}

	// Pointer types are represented via StarExpr nodes.

	// A FuncType node represents a function type.
	FuncType struct {
		Func       token.Pos  // position of "func" keyword (token.NoPos if there is no "func")
		TypeParams *FieldList // type parameters; or nil
		Params     *FieldList // (incoming) parameters; non-nil
		Results    *FieldList // (outgoing) results; or nil
	}

	// An InterfaceType node represents an interface type.
	InterfaceType struct {
		Interface  token.Pos  // position of "interface" keyword
		Methods    *FieldList // list of embedded interfaces, methods, or types
		Incomplete bool       // true if (source) methods or types are missing in the Methods list
	}

	// A MapType node represents a map type.
	MapType struct {
		Map   token.Pos // position of "map" keyword
		Key   Expr
		Value Expr
	}

	// A ChanType node represents a channel type.
	ChanType struct {
		Begin token.Pos // position of "chan" keyword or "<-" (whichever comes first)
		Arrow token.Pos // position of "<-" (token.NoPos if there is no "<-")
		Dir   ChanDir   // channel direction
		Value Expr      // value type
	}
)

```



### 2 Statement

赋值语句，控制语句（if，else,for，select...）等均属于statement node。

```go
// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.
type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	//
	BadStmt struct {
		From, To token.Pos // position range of bad statement
	}

	// A DeclStmt node represents a declaration in a statement list.
	DeclStmt struct {
		Decl Decl // *GenDecl with CONST, TYPE, or VAR token
	}

	// An EmptyStmt node represents an empty statement.
	// The "position" of the empty statement is the position
	// of the immediately following (explicit or implicit) semicolon.
	//
	EmptyStmt struct {
		Semicolon token.Pos // position of following ";"
		Implicit  bool      // if set, ";" was omitted in the source
	}

	// A LabeledStmt node represents a labeled statement.
	LabeledStmt struct {
		Label *Ident
		Colon token.Pos // position of ":"
		Stmt  Stmt
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	//
	ExprStmt struct {
		X Expr // expression
	}

	// A SendStmt node represents a send statement.
	SendStmt struct {
		Chan  Expr
		Arrow token.Pos // position of "<-"
		Value Expr
	}

	// An IncDecStmt node represents an increment or decrement statement.
	IncDecStmt struct {
		X      Expr
		TokPos token.Pos   // position of Tok
		Tok    token.Token // INC or DEC
	}

	// An AssignStmt node represents an assignment or
	// a short variable declaration.
	//
	AssignStmt struct {
		Lhs    []Expr
		TokPos token.Pos   // position of Tok
		Tok    token.Token // assignment token, DEFINE
		Rhs    []Expr
	}

	// A GoStmt node represents a go statement.
	GoStmt struct {
		Go   token.Pos // position of "go" keyword
		Call *CallExpr
	}

	// A DeferStmt node represents a defer statement.
	DeferStmt struct {
		Defer token.Pos // position of "defer" keyword
		Call  *CallExpr
	}

	// A ReturnStmt node represents a return statement.
	ReturnStmt struct {
		Return  token.Pos // position of "return" keyword
		Results []Expr    // result expressions; or nil
	}

	// A BranchStmt node represents a break, continue, goto,
	// or fallthrough statement.
	//
	BranchStmt struct {
		TokPos token.Pos   // position of Tok
		Tok    token.Token // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
		Label  *Ident      // label name; or nil
	}

	// A BlockStmt node represents a braced statement list.
	BlockStmt struct {
		Lbrace token.Pos // position of "{"
		List   []Stmt
		Rbrace token.Pos // position of "}", if any (may be absent due to syntax error)
	}

	// An IfStmt node represents an if statement.
	IfStmt struct {
		If   token.Pos // position of "if" keyword
		Init Stmt      // initialization statement; or nil
		Cond Expr      // condition
		Body *BlockStmt
		Else Stmt // else branch; or nil
	}

	// A CaseClause represents a case of an expression or type switch statement.
	CaseClause struct {
		Case  token.Pos // position of "case" or "default" keyword
		List  []Expr    // list of expressions or types; nil means default case
		Colon token.Pos // position of ":"
		Body  []Stmt    // statement list; or nil
	}

	// A SwitchStmt node represents an expression switch statement.
	SwitchStmt struct {
		Switch token.Pos  // position of "switch" keyword
		Init   Stmt       // initialization statement; or nil
		Tag    Expr       // tag expression; or nil
		Body   *BlockStmt // CaseClauses only
	}

	// A TypeSwitchStmt node represents a type switch statement.
	TypeSwitchStmt struct {
		Switch token.Pos  // position of "switch" keyword
		Init   Stmt       // initialization statement; or nil
		Assign Stmt       // x := y.(type) or y.(type)
		Body   *BlockStmt // CaseClauses only
	}

	// A CommClause node represents a case of a select statement.
	CommClause struct {
		Case  token.Pos // position of "case" or "default" keyword
		Comm  Stmt      // send or receive statement; nil means default case
		Colon token.Pos // position of ":"
		Body  []Stmt    // statement list; or nil
	}

	// A SelectStmt node represents a select statement.
	SelectStmt struct {
		Select token.Pos  // position of "select" keyword
		Body   *BlockStmt // CommClauses only
	}

	// A ForStmt represents a for statement.
	ForStmt struct {
		For  token.Pos // position of "for" keyword
		Init Stmt      // initialization statement; or nil
		Cond Expr      // condition; or nil
		Post Stmt      // post iteration statement; or nil
		Body *BlockStmt
	}

	// A RangeStmt represents a for statement with a range clause.
	RangeStmt struct {
		For        token.Pos   // position of "for" keyword
		Key, Value Expr        // Key, Value may be nil
		TokPos     token.Pos   // position of Tok; invalid if Key == nil
		Tok        token.Token // ILLEGAL if Key == nil, ASSIGN, DEFINE
		Range      token.Pos   // position of "range" keyword
		X          Expr        // value to range over
		Body       *BlockStmt
	}
)

```



### 3 Spec Node

Spec节点表示单个（非括号括起的）导入、常量、类型或变量声明。

Spec node只有3种，分别是`ImportSpec`，`ValueSpec`和`TypeSpec`：

ImportSpec表示一个单独的import，ValueSpec一个常量或变量的声明，TypeSpec 则表示一个type声明。

```go
/ ----------------------------------------------------------------------------
// Declarations

// A Spec node represents a single (non-parenthesized) import,
// constant, type, or variable declaration.
type (
	// The Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec.
	Spec interface {
		Node
		specNode()
	}

	// An ImportSpec node represents a single package import.
	ImportSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Name    *Ident        // local package name (including "."); or nil
		Path    *BasicLit     // import path
		Comment *CommentGroup // line comments; or nil
		EndPos  token.Pos     // end of spec (overrides Path.Pos if nonzero)
	}

	// A ValueSpec node represents a constant or variable declaration
	// (ConstSpec or VarSpec production).
	//
	ValueSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Names   []*Ident      // value names (len(Names) > 0)
		Type    Expr          // value type; or nil
		Values  []Expr        // initial values; or nil
		Comment *CommentGroup // line comments; or nil
	}

	// A TypeSpec node represents a type declaration (TypeSpec production).
	TypeSpec struct {
		Doc        *CommentGroup // associated documentation; or nil
		Name       *Ident        // type name
		TypeParams *FieldList    // type parameters; or nil
		Assign     token.Pos     // position of '=', if any
		Type       Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
		Comment    *CommentGroup // line comments; or nil
	}
)

```



### 4 others

```
// Comment 注释节点，代表单行的 //-格式 或 /*-格式的注释.
type Comment struct {
    ...
}
...
// CommentGroup 注释块节点，包含多个连续的Comment
type CommentGroup struct {
    ...
}

// Field 字段节点, 可以代表结构体定义中的字段，接口定义中的方法列表，函数前面中的入参和返回值字段
type Field struct {
    ...
}
...
// FieldList 包含多个Field
type FieldList struct {
    ...
}

// File 表示一个文件节点
type File struct {
	...
}

// Package 表示一个包节点
type Package struct {
    ...
}

```

## 基本用法

### 类型定义

我们常用的node,一般是有一些比较基础的类型组成的，**ast.File**，**ast.ValueSpec:**  ， **ast.TypeSpec** ，**ast.Field** 类型，比如TypeSpec 类型下 会doc(CommentGroup) 和filed(FieldList) 字段。ValueSpec 类型下有 ident，doc类型的字段， 具体使用还是要debug 去看 分析这个node下有什么字段。

```go
	// A TypeSpec node represents a type declaration (TypeSpec production).
	TypeSpec struct {
		Doc        *CommentGroup // associated documentation; or nil
		Name       *Ident        // type name
		TypeParams *FieldList    // type parameters; or nil
		Assign     token.Pos     // position of '=', if any
		Type       Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
		Comment    *CommentGroup // line comments; or nil
	}
```

```go
	ValueSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Names   []*Ident      // value names (len(Names) > 0)
		Type    Expr          // value type; or nil
		Values  []Expr        // initial values; or nil
		Comment *CommentGroup // line comments; or nil
	}
```



![](https://cdn.learnku.com/uploads/images/202401/05/51993/e5EqOvzYTY.png!large)



![](https://cdn.learnku.com/uploads/images/202401/05/51993/Rhx6mwZQWq.png!large)

### 遍历语法树

```go
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

```

### 修改语法树

 给gorm gen 生成的代码增加一些我们自己的方法。直接在ast.FIle node下增加即可

```go
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
			

				}
				node.Decls = append(node.Decls, newSpec, newSpecCtx)
			}

			return false

		}
		return true
	})
```

```go
func (c cardDo) Scan(result interface{}) (err error) {
    return c.DO.Scan(result)
}

func (c cardDo) Delete(models ...*model.Card) (result gen.ResultInfo, err error) {
    return c.DO.Delete(models)
}

func (c *cardDo) withDO(do gen.Dao) *cardDo {
    c.DO = *do.(*gen.DO)
    return c
}
//新增的方法
func (c card) FindCardById(id int32) (*model.Card, error) {
    return c.cardDo.Where(c.ID.Eq(id)).Take()
}
```

## 例子一枚

根据我们定义的常量错误码，生成具体的错误。

```go
package main

import (
	"google.golang.org/grpc/codes"
)

const (
	UnAuthorizedCode codes.Code = iota + 403 //403未授权
	NotFoundCode                             //404未找到
)




```

```go
//==============================target==========================================
// Code generated by tool. DO NOT EDIT.
// Code generated by tool. DO NOT EDIT.
// Code generated by tool. DO NOT EDIT.

package main

import (
	"google.golang.org/grpc/status"
)

var (
	UnAuthorized    =   status.Error(UnAuthorizedCode,"403未授权")
	NotFound        =   status.Error(NotFoundCode,"404未找到")
)
/*

403     403未授权    
404     404未找到    


*/
```



```go
package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)
// 定义文件描述结构体
type fileDesc struct {
	FileName string
	PkgName          string
	ConstantDescList []*constantDesc
	FileComments string
}
// 定义常量描述结构体
type constantDesc struct {
	CodeName string
	VarName string
	Value   int32
	Comment string
	ConcatStr string
}

func main() {
	fileFullPath, err := filepath.Abs("utils\\ast\\code.go")
	if err!=nil{
		panic(errors.WithMessage(err,"获取文件路径失败"))
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileFullPath, nil, parser.ParseComments)
	if err != nil {
		panic(errors.WithMessage(err,"解析文件失败"))
	}

	const (
		filedSuffix = "Code"
	)
	var (
		currentVal = int32(0)
		fd =  &fileDesc{}
	)
	fd.FileName = fileFullPath
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.File:
			// 获取文件包名
			if node.Name == nil {
				panic("invalid go file")
			}
			fd.PkgName = node.Name.Name
		case *ast.ValueSpec:
			// 获取文件中定义的常量
			d := &constantDesc{}
			names := node.Names
			if len(names) == 0 {
				return true
			}
			//name
			name := names[0].Name
			if !strings.HasSuffix(name, filedSuffix) {
				return true
			}
			d.CodeName = name
			d.VarName = strings.TrimSuffix(name, filedSuffix)
			//Value
			for _, v := range node.Values {
				e, ok := v.(*ast.BinaryExpr)
				if ok {
					l, ok := e.Y.(*ast.BasicLit)
					if ok {
						d.Value = cast.ToInt32(l.Value)
						currentVal = d.Value
					}
				}
			}
			if d.Value == 0 {
				currentVal++
				d.Value = currentVal
			}
			//Comment
			var comment string
			if node.Comment != nil {
				comment = node.Comment.Text()
			}
			if node.Doc != nil {
				comment = node.Doc.Text()
			}
			d.Comment = strings.TrimSpace(comment)
			fd.ConstantDescList = append(fd.ConstantDescList, d)

		}
		return true
	})
	log.Printf("%+v", fd)
	castToGrpcError(fd)
}

const tpl = `
// Code generated by tool. DO NOT EDIT.
// Code generated by tool. DO NOT EDIT.
// Code generated by tool. DO NOT EDIT.

package {{.PkgName}}

import (
	"google.golang.org/grpc/status"
)

var (
	{{- range $idx, $constDesc :=  .ConstantDescList}}
	{{ $constDesc.ConcatStr }}	
	{{- end }}
)
/*

{{ .FileComments }}

*/
`
//根据解析出来的数据转为grpc的错误定义
func castToGrpcError(c *fileDesc)  {
	const fileSuffix = ".err.go"
	p, err := template.New("").Funcs(template.FuncMap{}).Parse(tpl)

	if err !=nil{
		panic(errors.WithMessage(err,"parse template failed"))
	}
	base := filepath.Dir(c.FileName)
	n := strings.TrimSuffix(filepath.Base(c.FileName),".go")
	fileFullPath := filepath.Join(base,n+fileSuffix)
	f, err := os.OpenFile(fileFullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err !=nil{
		panic(errors.WithMessage(err,"open file failed"))
	}
	var (

		fileComments string
	)

	for _,v := range c.ConstantDescList {
		v.ConcatStr = fmt.Sprintf("%-15s %-3s %-8s",v.VarName,"=",fmt.Sprintf(`status.Error(%v,%v)`,v.CodeName,"\""+v.Comment+"\""))
		fileComments +=fmt.Sprintf("%-8d%-10s\n",v.Value,v.Comment)
	}

	c.FileComments = fileComments
	if err := p.Execute(f, c);err!=nil{
		panic(errors.WithMessage(err,"execute template failed"))
	}

}

```

