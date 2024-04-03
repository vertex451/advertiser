package s3storage

import (
	cfg "advertiser/shared/config/config"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestService_Load(t *testing.T) {
	filePath := "/Users/vertex451/workplace/silverspase/tg-bot/golang/agency/example.jpg"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the contents of the file into a buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		return
	}

	cfg.Load()
	s := New()
	res, err := s.Store("test.jpg", bytes.NewReader(buf.Bytes()))
	assert.Nil(t, err)
	fmt.Println("### res", res)
}
