package uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
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

	// save file
	for _, fileHeader := range fileHeader {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		name, ex := strings.Split(fileHeader.Filename, ".")[0], strings.Split(fileHeader.Filename, ".")[1]

		// dir/file-id.ext
		filedes := fmt.Sprintf("%v%v-%v.%v", l.baseDir, name, uuid.New().String(), ex)

		dst, err := os.Create(filedes)
		if err != nil {
			return nil, err
		}

		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			return nil, err
		}

		filenames = append(filenames, fileHeader.Filename)

		log.Info("Image uploaded successfully")
	}

	return filenames, nil
}
