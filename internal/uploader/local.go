package uploader

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

func (l *LocalUpload) Upload(r *http.Request, grp string, maxMem int64, allowMime []string) ([]string, error) {

	r.ParseMultipartForm(maxMem)

	FileFields := r.MultipartForm.File

	if r.Method == "PUT" && FileFields == nil {
		return nil, nil
	}

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

	path := filepath.Join(l.baseDir, grp)

	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		_ = os.Mkdir(path, 0755)
	}

	// save file
	for _, fileHeader := range fileHeader {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		name, ex := strings.Split(fileHeader.Filename, ".")[0], strings.Split(fileHeader.Filename, ".")[1]

		// file-uuid.ext
		fp := fmt.Sprintf("%s-%s.%s", name, uuid.New().String(), ex)

		dst, err := os.Create(filepath.Join(l.baseDir, grp, fp))
		if err != nil {
			return nil, err
		}

		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			return nil, err
		}

		filenames = append(filenames, fp)

		log.Info("Image uploaded successfully")
	}

	return filenames, nil
}
