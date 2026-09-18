package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil && (os.IsNotExist(err) || errors.Is(err, os.ErrNotExist)) {
		return nil
	}
	return err
}
