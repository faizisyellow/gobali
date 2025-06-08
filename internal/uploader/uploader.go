package uploader

import (
	"net/http"
	"net/textproto"
)

type Uploader interface {
	Upload(r *http.Request, dst string, maxMem *int64, allowType textproto.MIMEHeader) ([]string, error)
}
