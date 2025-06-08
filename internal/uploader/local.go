package uploader

import (
	"net/http"
	"net/textproto"
)

type LocalUpload struct {
	baseDir string
}

func NewLocalUpload(baseDir string) *LocalUpload {
	return &LocalUpload{baseDir}
}

// TODO: upload the file
func (l *LocalUpload) Upload(r *http.Request, dst string, maxMem *int64, allowType textproto.MIMEHeader) ([]string, error) {

	r.ParseMultipartForm(*maxMem)

	return []string{""}, nil
}
