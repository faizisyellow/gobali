package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct{}

func (e *Env) Set() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}

func (e *Env) GetString(key, fallback string) string {

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func (e *Env) GetInt(key string, fallback int) int {

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}

func (e *Env) GetBool(key string, fallback bool) bool {

	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}
