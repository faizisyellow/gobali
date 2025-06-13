package uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

type LocalUpload struct {
	baseDir string
}

func NewLocalUpload(baseDir string) *LocalUpload {
	return &LocalUpload{baseDir}
}

func (l *LocalUpload) Upload(r *http.Request, dst string, maxMem int64, allowMime []string) ([]string, error) {

	r.ParseMultipartForm(maxMem)

	FileFields := r.MultipartForm.File

	fileHeader := []*multipart.FileHeader{}

	filenames := []string{}

	for _, headerFiles := range FileFields {
		headerFile := headerFiles[0]

		contentType := headerFile.Header["Content-Type"][0]

		fmt.Println("Detected file:", contentType)
		fmt.Println("Want  file:", allowMime)

		err := ValidateFile(allowMime, contentType)
		if err != nil {
			return nil, err
		}

		err = ValidateSize(headerFile.Size, maxMem)
		if err != nil {
			return nil, err
		}

		fileHeader = append(fileHeader, headerFile)

	}

	for _, fileHeader := range fileHeader {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		dst, err := os.Create(l.baseDir + fileHeader.Filename)
		if err != nil {
			return nil, err
		}

		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			return nil, err
		}

		filenames = append(filenames, fileHeader.Filename)

		log.Info("Success image uploaded")
	}

	return filenames, nil
}
