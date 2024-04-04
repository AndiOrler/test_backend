package env

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadLocalFile Loads the local .env file
func LoadLocalFile(filenames ...string) error {
	// load local env file
	err := godotenv.Load(filenames...)
	if err != nil {
		return err
	}
	return nil
}

// GetStr Returns env string
func GetStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, errors.New("getenv: environment variable empty")
	}
	return v, nil
}

// GetInt Returns env int
func GetInt(key string) (int, error) {
	s, err := GetStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// GetFloat Returns env float
func GetFloat(key string) (float64, error) {
	s, err := GetStr(key)
	if err != nil {
		return 0, err
	}
	bitSize := 64
	v, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// GetBool Returns env bool
func GetBool(key string) (bool, error) {
	s, err := GetStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}
