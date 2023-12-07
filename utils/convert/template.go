package main

const tpl = `
syntax = "proto3";

package pb;

option go_package = "./pb";

{{ range .Msg}}
message {{.MessageName}}{
	{{  range $idx, $field :=  .Fields}}
	{{$field.FieldType}} {{$field.FieldName}} = {{ add $idx 1}};	
	{{- end }}
}

{{- end }}


`
