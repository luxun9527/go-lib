package main

import (
	"github.com/spf13/cast"
	"os"
	"testing"
	"text/template"
)

func TestTemplate(t *testing.T) {
	p, err := template.New("").Funcs(template.FuncMap{"add": func(x, y int64) string {
		return cast.ToString(x + y)
	}}).Parse(tpl)

	if err != nil {
		t.Fatal("parse failed", err)
	}

	d := Data{
		Msg: []*Message{{
			MessageName: "m1",
			Fields: []*Field{{
				FieldName: "username",
				FieldType: "string",
			}},
		}},
	}

	p.Execute(os.Stdout, d)
}

type Data struct {
	Msg []*Message
}
type Message struct {
	MessageName string
	Fields      []*Field
}
type Field struct {
	FieldName string
	FieldType string
}
