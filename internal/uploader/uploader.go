package uploader

import (
	"errors"
	"net/http"
	"slices"
)

var (
	ErrExtNotAllowed = errors.New("extension not allowed")
	ErrSizeLarger    = errors.New("size more than limit")
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
		return ErrSizeLarger
	}

	return nil
}
