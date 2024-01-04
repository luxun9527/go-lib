#  AST语法树

refer

https://juejin.cn/post/6844903982683389960#heading-2

## 核心概念

go文件，可以视为一颗由方法，包，字段，变量，注释组成的语法树，可以使用go提供的api，可以获取go文件中的方法，包，字段，变量，注释等信息。

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

Spec node只有3种，分别是`ImportSpec`，`ValueSpec`和`TypeSpec`：

ImportSpec表示一个单独的import，ValueSpec一个常量或变量的声明，TypeSpec 则表示一个type声明。

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

## 语法树遍历

各种类型互相嵌套，一般用的比较多的是获取结构体和字段的定义。

```go
package main

import (
	"github.com/fatih/color"
	"go/ast"
	"go/parser"
	"go/token"
)

// 遍历节点
func main() {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "D:\\project\\go-lib\\utils\\ast\\card.gen.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.GenDecl:
			genDecl := node.(*ast.GenDecl)
			color.Red("GenDecl node %+v",genDecl)
		case  *ast.TypeSpec: //类型定义 ((重点)) 
			typeSpec := node.(*ast.TypeSpec)
			color.Yellow("TypeSpec node %+v",typeSpec)
		case *ast.StructType:
			structType := node.(*ast.StructType)
			color.Yellow("StructType node %+v", structType)
		case *ast.Field: //字段定义 ((重点))
			field := node.(*ast.Field)
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
			color.Blue("Field node %+v",field)
		case *ast.Ident://标识符，包括类型和字段名，变量名，结构体名等
			ident := node.(*ast.Ident)
			color.Green("ident node %+v",ident)
		case *ast.SwitchStmt:
			switchStmt := node.(*ast.SwitchStmt)
			color.Magenta("SwitchStmt node %+v",switchStmt)
		case *ast.Comment:
			comment := node.(*ast.Comment)
			color.Black("comment node %+v",comment)
		case *ast.CommentGroup:
			comment := node.(*ast.CommentGroup)
			color.Magenta("commentGroup node %+v",comment)
		case *ast.InterfaceType:
			interfaceType := node.(*ast.InterfaceType)
			color.Magenta("interfaceType node %+v",interfaceType)
		}
		return true
	})
}

```



![image-20240104164144046](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20240104164144046.png)

![image-20240104164502529](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20240104164502529.png)