// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.4
// source: grpcdemo/custom.proto

//当别人导入这个protobuf文件，使用的包名 如 google/protobuf/empty.proto 定义的就是 package google.protobuf;

package grpcdemo

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CustomMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomMessage string `protobuf:"bytes,1,opt,name=custom_message,json=customMessage,proto3" json:"custom_message,omitempty"`
}

func (x *CustomMessage) Reset() {
	*x = CustomMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpcdemo_custom_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CustomMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CustomMessage) ProtoMessage() {}

func (x *CustomMessage) ProtoReflect() protoreflect.Message {
	mi := &file_grpcdemo_custom_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CustomMessage.ProtoReflect.Descriptor instead.
func (*CustomMessage) Descriptor() ([]byte, []int) {
	return file_grpcdemo_custom_proto_rawDescGZIP(), []int{0}
}

func (x *CustomMessage) GetCustomMessage() string {
	if x != nil {
		return x.CustomMessage
	}
	return ""
}

var File_grpcdemo_custom_proto protoreflect.FileDescriptor

var file_grpcdemo_custom_proto_rawDesc = []byte{
	0x0a, 0x15, 0x67, 0x72, 0x70, 0x63, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x63, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x67, 0x72, 0x70, 0x63, 0x64, 0x65, 0x6d,
	0x6f, 0x22, 0x36, 0x0a, 0x0d, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x75, 0x73, 0x74,
	0x6f, 0x6d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x15, 0x5a, 0x13, 0x2e, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x64, 0x65, 0x6d, 0x6f, 0x3b, 0x67, 0x72, 0x70, 0x63, 0x64, 0x65, 0x6d, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpcdemo_custom_proto_rawDescOnce sync.Once
	file_grpcdemo_custom_proto_rawDescData = file_grpcdemo_custom_proto_rawDesc
)

func file_grpcdemo_custom_proto_rawDescGZIP() []byte {
	file_grpcdemo_custom_proto_rawDescOnce.Do(func() {
		file_grpcdemo_custom_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpcdemo_custom_proto_rawDescData)
	})
	return file_grpcdemo_custom_proto_rawDescData
}

var file_grpcdemo_custom_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_grpcdemo_custom_proto_goTypes = []interface{}{
	(*CustomMessage)(nil), // 0: grpcdemo.CustomMessage
}
var file_grpcdemo_custom_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpcdemo_custom_proto_init() }
func file_grpcdemo_custom_proto_init() {
	if File_grpcdemo_custom_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpcdemo_custom_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CustomMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpcdemo_custom_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_grpcdemo_custom_proto_goTypes,
		DependencyIndexes: file_grpcdemo_custom_proto_depIdxs,
		MessageInfos:      file_grpcdemo_custom_proto_msgTypes,
	}.Build()
	File_grpcdemo_custom_proto = out.File
	file_grpcdemo_custom_proto_rawDesc = nil
	file_grpcdemo_custom_proto_goTypes = nil
	file_grpcdemo_custom_proto_depIdxs = nil
}
