package stream

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
)

/*
实现功能
1. 创建一个multipart/form-data格式的流，在读取文件的时候不将整个文件读入内存。
使用chunkded模式传输数据
*/
type MultipartReaderWriter struct {
	buf            *bytes.Buffer
	closeData      *bytes.Buffer
	r              io.Reader
	hasStartedFile bool
	boundary       string
	contentType    string
	writer         *multipart.Writer
	FileFiledList  []*FileField
}

func (fr *MultipartReaderWriter) Read(p []byte) (int, error) {
	// 如果文件读取还未开始，先从 buffer 中读取数据
	if !fr.hasStartedFile {
		n, err := fr.buf.Read(p)
		if n > 0 || err != io.EOF {
			return n, err
		}

		fr.hasStartedFile = true
		if _, err := fmt.Fprintf(fr.buf, "\r\n--%s--\r\n", fr.boundary); err != nil {
			return 0, err
		}
	}
	n, err := fr.r.Read(p)
	if err != nil {
		//读取文件数据后，将close数据读取出来
		if errors.Is(err, io.EOF) {
			return fr.closeData.Read(p)
		}
		return 0, err
	}
	return n, err
}

func (fr *MultipartReaderWriter) WriteFiled(key, value string) error {
	return fr.writer.WriteField(key, value)
}

func (fr *MultipartReaderWriter) WriteFileField(filedName, fileFiledName string, data io.Reader) error {
	if data == nil {
		return errors.New("data is nil")
	}

	field, err := NewFileField(data, filedName, fileFiledName, fr.boundary, fr.buf.Len() > 0)
	if err != nil {
		return err
	}
	fr.FileFiledList = append(fr.FileFiledList, field)
	return err
}

func (fr *MultipartReaderWriter) Close() error {
	fields := make([]io.Reader, 0, len(fr.FileFiledList))
	for _, v := range fr.FileFiledList {
		fields = append(fields, v)
	}
	fr.r = io.MultiReader(fields...)

	return nil
}

func NewMultipartReaderWriter() (*MultipartReaderWriter, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 500))
	closeData := bytes.NewBuffer(make([]byte, 0, 500))

	writer := multipart.NewWriter(buffer)

	if _, err := fmt.Fprintf(closeData, "\r\n--%s--\r\n", writer.Boundary()); err != nil {
		return nil, err
	}
	multipartReader := &MultipartReaderWriter{
		buf:         buffer,
		closeData:   closeData,
		contentType: writer.FormDataContentType(),
		boundary:    writer.Boundary(),
		writer:      writer,
	}

	return multipartReader, nil
}

type FileField struct {
	r              io.Reader
	buf            *bytes.Buffer
	hasStartedFile bool
	hasPrev        bool
}

func NewFileField(r io.Reader, fileName, fieldName, boundary string, hasPrev bool) (*FileField, error) {

	buffer := bytes.NewBuffer(make([]byte, 0, 500))
	if hasPrev {
		if _, err := fmt.Fprintf(buffer, "\r\n"); err != nil {
			return nil, err
		}
	}
	writer := multipart.NewWriter(buffer)
	if err := writer.SetBoundary(boundary); err != nil {
		return nil, err
	}
	if _, err := writer.CreateFormFile(fieldName, fileName); err != nil {
		return nil, err
	}
	return &FileField{
		r:              r,
		buf:            buffer,
		hasStartedFile: false,
	}, nil
}

func (fr *FileField) Read(p []byte) (int, error) {
	// 如果文件读取还未开始，先从 buffer 中读取数据
	if !fr.hasStartedFile {
		n, err := fr.buf.Read(p)
		if n > 0 || err != io.EOF {
			return n, err
		}

		fr.hasStartedFile = true
	}

	return fr.r.Read(p)
}
