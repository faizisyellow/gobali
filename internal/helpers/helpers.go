package helpers

import (
	"os"
	"reflect"

	"github.com/charmbracelet/log"
)

func RemoveFile(filepath string) error {
	if err := os.Remove(filepath); err != nil {
		log.Error(err)
		return err
	}

	log.Info("Remove image successfully")
	return nil
}

func CreateNewStructByReflect[T any](data *T, fields ...string) *T {
	origVal := reflect.ValueOf(data).Elem()
	origType := origVal.Type()

	newVal := reflect.New(origType).Elem()

	fieldSet := make(map[string]bool)
	for _, f := range fields {
		fieldSet[f] = true
	}

	for i := 0; i < origVal.NumField(); i++ {
		field := origType.Field(i)
		if fieldSet[field.Name] {
			newVal.Field(i).Set(origVal.Field(i))
		}
	}

	newInstance := newVal.Addr().Interface().(*T)

	return newInstance
}
