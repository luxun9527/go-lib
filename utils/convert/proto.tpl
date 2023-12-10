syntax = "proto3";
package pb;
option go_package = "./pb";

{{- if .AddHttp }}
import "google/api/annotations.proto";
{{ end }}

{{- range .Msg}}
message {{.MessageName}}{
	{{- range $idx, $field :=  .Fields}}
	{{$field.FieldType}} {{$field.FieldName}} = {{ add $idx 1}}; {{$field.Comment}}
	{{- end }}
}

message Get{{.MessageName}}ListReq{}

message Get{{.MessageName}}ListResp{
    repeated {{.MessageName}} list=4;
}



message Update{{.MessageName}}Resp{}

message Create{{.MessageName}}Resp{}

message Delete{{.MessageName}}Req{
    int64 id=1;
}

message Delete{{.MessageName}}Resp{
    int64 id=1;
}

{{- if $.AddHttp }}

service {{.MessageName}}Service{
	//获取{{.MessageName}}列表
	rpc Get{{.MessageName}}List(Get{{.MessageName}}ListReq)returns(Get{{.MessageName}}ListResp){
	    option (google.api.http) = {
           get: "/api/v1/get{{.MessageName}}List"
        };
	};
	//修改{{.MessageName}}
	rpc Update{{.MessageName}}({{.MessageName}})returns(Update{{.MessageName}}Resp){
	    option (google.api.http) = {
            post: "/api/v1/update{{.MessageName}}"
            body: "*"
        };
	};
    //创建{{.MessageName}}
	rpc Create{{.MessageName}}({{.MessageName}})returns(Create{{.MessageName}}Resp){
	    option (google.api.http) = {
            post: "/api/v1/create{{.MessageName}}"
            body: "*"
        };
	};
    //删除{{.MessageName}}
	rpc Delete{{.MessageName}}(Delete{{.MessageName}}Req)returns(Delete{{.MessageName}}Resp){
	    option (google.api.http) = {
             post: "/api/v1/del{{.MessageName}}"
             body: "*"
        };
	};
}
{{- else}}

service {{.MessageName}}Service{
	//获取{{.MessageName}}列表
	rpc Get{{.MessageName}}List(Get{{.MessageName}}ListReq)returns(Get{{.MessageName}}ListResp);
	//修改{{.MessageName}}
	rpc Update{{.MessageName}}({{.MessageName}})returns(Update{{.MessageName}}Resp);
    //创建{{.MessageName}}
	rpc Create{{.MessageName}}({{.MessageName}})returns(Create{{.MessageName}}Resp);
    //删除{{.MessageName}}
	rpc Delete{{.MessageName}}({{.MessageName}})returns(Create{{.MessageName}}Resp);
}

{{- end }}




{{- end }}
