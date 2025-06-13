package uploader

import (
	"os"
	"reflect"
	"testing"
)

// TODO: test upload
func TestLocalUploaded(t *testing.T) {
	dst := "./assets-test/output-assets/"
	os.MkdirAll(dst, 0755)
	defer os.RemoveAll(dst)

	lu := NewLocalUpload(dst)

	t.Run("should success upload images", func(t *testing.T) {

		var maxMemo int64 = 3 * 1024 * 1024 // 1 mb
		allowMime := []string{"image/png", "image/jpeg", "image/jpg"}

		want := []string{"image1.jpeg"}

		result, err := lu.Upload(nil, dst, maxMemo, allowMime)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(result, want) {
			t.Errorf("expected: %v but got: %v", want, result)
		}
	})
}
