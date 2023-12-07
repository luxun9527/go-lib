package main

import (
	"os"
	"testing"
	"text/template"
)

func TestTemplate(t *testing.T) {
	p, err := template.New("").
		Funcs(template.FuncMap{
			"add": func(x, y int) int {

				return x + y
			},
		}).Parse(tpl)

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

	if err := p.Execute(os.Stdout, d); err != nil {
		t.Errorf("%v", err)
		return
	}
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
