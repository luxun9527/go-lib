// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.4
// source: grpcdemo/folder/imported.proto

//当别人导入这个protobuf文件，使用的包名 如 google/protobuf/empty.proto 定义的就是 package google.protobuf;

package folder

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

type ImportedMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImportedMessage string `protobuf:"bytes,1,opt,name=imported_message,json=importedMessage,proto3" json:"imported_message,omitempty"`
}

func (x *ImportedMessage) Reset() {
	*x = ImportedMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpcdemo_folder_imported_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImportedMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImportedMessage) ProtoMessage() {}

func (x *ImportedMessage) ProtoReflect() protoreflect.Message {
	mi := &file_grpcdemo_folder_imported_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImportedMessage.ProtoReflect.Descriptor instead.
func (*ImportedMessage) Descriptor() ([]byte, []int) {
	return file_grpcdemo_folder_imported_proto_rawDescGZIP(), []int{0}
}

func (x *ImportedMessage) GetImportedMessage() string {
	if x != nil {
		return x.ImportedMessage
	}
	return ""
}

var File_grpcdemo_folder_imported_proto protoreflect.FileDescriptor

var file_grpcdemo_folder_imported_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x67, 0x72, 0x70, 0x63, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x66, 0x6f, 0x6c, 0x64, 0x65,
	0x72, 0x2f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0f, 0x67, 0x72, 0x70, 0x63, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x66, 0x6f, 0x6c, 0x64, 0x65,
	0x72, 0x22, 0x3c, 0x0a, 0x0f, 0x49, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42,
	0x2b, 0x5a, 0x29, 0x67, 0x6f, 0x2d, 0x6c, 0x69, 0x62, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x66,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x3b, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpcdemo_folder_imported_proto_rawDescOnce sync.Once
	file_grpcdemo_folder_imported_proto_rawDescData = file_grpcdemo_folder_imported_proto_rawDesc
)

func file_grpcdemo_folder_imported_proto_rawDescGZIP() []byte {
	file_grpcdemo_folder_imported_proto_rawDescOnce.Do(func() {
		file_grpcdemo_folder_imported_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpcdemo_folder_imported_proto_rawDescData)
	})
	return file_grpcdemo_folder_imported_proto_rawDescData
}

var file_grpcdemo_folder_imported_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_grpcdemo_folder_imported_proto_goTypes = []interface{}{
	(*ImportedMessage)(nil), // 0: grpcdemo.folder.ImportedMessage
}
var file_grpcdemo_folder_imported_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpcdemo_folder_imported_proto_init() }
func file_grpcdemo_folder_imported_proto_init() {
	if File_grpcdemo_folder_imported_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpcdemo_folder_imported_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImportedMessage); i {
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
			RawDescriptor: file_grpcdemo_folder_imported_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_grpcdemo_folder_imported_proto_goTypes,
		DependencyIndexes: file_grpcdemo_folder_imported_proto_depIdxs,
		MessageInfos:      file_grpcdemo_folder_imported_proto_msgTypes,
	}.Build()
	File_grpcdemo_folder_imported_proto = out.File
	file_grpcdemo_folder_imported_proto_rawDesc = nil
	file_grpcdemo_folder_imported_proto_goTypes = nil
	file_grpcdemo_folder_imported_proto_depIdxs = nil
}