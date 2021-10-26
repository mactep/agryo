package util

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func GetTestdata(t *testing.T, file string) io.ReadCloser {
	t.Helper()

	f, err := os.Open(filepath.Join("testdata", file))
	if err != nil {
		t.Error(err)
	}

	return f
}

func GetEnvOrPanic(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s not defined", key)
	}

	return value
}
