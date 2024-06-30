package stream

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

/*
实现功能
1. 创建一个multipart/form-data格式的流，在读取文件的时候不将整个文件读入内存。
使用chunkded模式传输数据
*/
type MultipartReader struct {
	buf            *bytes.Buffer
	closeData      *bytes.Buffer
	r              io.Reader
	hasStartedFile bool
	boundary       string
	contentType    string
}

func (fr *MultipartReader) Read(p []byte) (int, error) {
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
		if errors.Is(err, io.EOF) {
			return fr.closeData.Read(p)
		}
		return 0, err
	}
	return n, err
}

func NewMultiparkReader() (*MultipartReader, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 500))
	closeData := bytes.NewBuffer(make([]byte, 0, 500))

	writer := multipart.NewWriter(buffer)
	file, err := os.Open("example.txt")
	if err != nil {
		return nil, err
	}
	if err := writer.WriteField("name", "zhangsan"); err != nil {
		return nil, err
	}

	reader1, err := NewFileFieldReader(file, "example.txt", "file", writer.Boundary(), true)
	if err != nil {
		return nil, err
	}
	file1, err := os.Open("example.txt")

	reader2, err := NewFileFieldReader(file1, "example.txt", "file1", writer.Boundary(), true)
	if err != nil {
		return nil, err
	}
	if _, err := fmt.Fprintf(closeData, "\r\n--%s--\r\n", writer.Boundary()); err != nil {
		return nil, err
	}
	multiReader := io.MultiReader(reader1, reader2)
	multipartReader := &MultipartReader{
		buf:         buffer,
		r:           multiReader,
		closeData:   closeData,
		contentType: writer.FormDataContentType(),
	}

	return multipartReader, nil
}

type FileFieldReader struct {
	r              io.Reader
	buf            *bytes.Buffer
	hasStartedFile bool
	hasPrev        bool
}

func NewFileFieldReader(r io.Reader, fileName, fieldName, boundary string, hasPrev bool) (*FileFieldReader, error) {

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

	return &FileFieldReader{
		r:              r,
		buf:            buffer,
		hasStartedFile: false,
	}, nil
}

func (fr *FileFieldReader) Read(p []byte) (int, error) {
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
