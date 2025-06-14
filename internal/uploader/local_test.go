package uploader

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func MultipartRequest(t *testing.T, fieldName, filePath string) *http.Request {
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Detect MIME type from file extension
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Set part headers manually
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filepath.Base(filePath)))
	h.Set("Content-Type", mimeType)

	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatalf("Failed to create part: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Failed to write file content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", "/", &body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func CheckEqual(t *testing.T, result, want []string) {
	if !reflect.DeepEqual(result, want) {
		t.Errorf("expected: %v but got: %v", want, result)
	}
}

func TestLocalUploaded(t *testing.T) {
	var maxMemo int64 = 3 * 1024 * 1024 // 3 mb

	dst := "./assets-test/output-assets/"
	os.MkdirAll(dst, 0755)

	// Remove after uploading
	defer os.RemoveAll(dst)

	lu := NewLocalUpload(dst)

	t.Run("should success upload images", func(t *testing.T) {
		req := MultipartRequest(t, "file", "./assets-test/input-assets/image1.jpeg")

		allowMime := []string{"image/png", "image/jpeg", "image/jpg"}

		want := []string{"image1.jpeg"}

		result, err := lu.Upload(req, "test", maxMemo, allowMime)

		if err != nil {
			t.Error(err)
		}

		CheckEqual(t, result, want)
	})

	t.Run("should fail upload file if the type not allowed", func(t *testing.T) {
		req := MultipartRequest(t, "file", "./assets-test/input-assets/excel.xlsx")

		allowMime := []string{"image/png", "image/jpeg", "image/jpg"}

		want := ErrExtNotAllowed

		_, err := lu.Upload(req, "test", maxMemo, allowMime)
		if err != want {
			t.Errorf("expected: %v but got: %v", want, err)
		}
	})

	t.Run("should fail upload file if the size is too large", func(t *testing.T) {

		maxMemo = 1 * 1024 * 1024 // 1 mb

		req := MultipartRequest(t, "file", "./assets-test/input-assets/image-large.png")

		allowMime := []string{"image/png"}

		want := ErrSizeLarger

		_, err := lu.Upload(req, "test", maxMemo, allowMime)
		if err != want {
			t.Errorf("expected: %v but got: %v", want, err)
		}
	})
}
