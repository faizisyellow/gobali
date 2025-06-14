package uploader

import (
	"errors"
	"net/http"
	"slices"
)

var (
	ErrExtNotAllowed = errors.New("this file type not allowed")
	ErrSizeTooLarger = errors.New("size is too large")
)

type Uploader interface {
	Upload(r *http.Request, dst string, maxMem int64, allowMime []string) ([]string, error)
}

func ValidateFile(allowMime []string, contentType string) error {
	if !slices.Contains(allowMime, contentType) {
		return ErrExtNotAllowed
	}

	return nil
}

func ValidateSize(size, limit int64) error {
	if size >= limit {
		return ErrSizeTooLarger
	}

	return nil
}
