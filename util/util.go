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

func GetEnvAndRaise(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Variável de ambiente %s não definida!", key)
	}

	return value
}
