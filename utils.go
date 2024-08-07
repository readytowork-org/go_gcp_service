package go_gcp_service

import (
	"errors"
	"path/filepath"
	"strconv"
	"time"
)

func GetFileName(filename string) (string, error) {
	if filename == "" {
		return "", errors.New("filename cannot be empty")
	}

	fileExtension := filepath.Ext(filename)
	return GenerateRandomFileName() + fileExtension, nil
}

// GenerateRandomFileName genrates the fileName with unique time
func GenerateRandomFileName() string {
	_time := time.Now().UnixNano()
	return strconv.FormatInt(_time, 10)
}
