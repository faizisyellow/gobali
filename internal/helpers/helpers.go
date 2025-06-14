package helpers

import (
	"os"

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
