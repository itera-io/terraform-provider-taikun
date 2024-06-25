package utils

import (
	"bytes"
	b64 "encoding/base64"
	tkcore "github.com/itera-io/taikungoclient/client"
	"os"
)

func StringPtr(s string) *string {
	return &s
}

func BoolPtr(s bool) *bool {
	return &s
}

func NewNullableFloat64(s float64) tkcore.NullableFloat64 {
	res := tkcore.NullableFloat64{}
	res.Set(&s)
	return res
}

// Read file contents, encode in base64 and return
func FilePathToBase64String(filePath string) (string, error) {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	fileContents = bytes.TrimSuffix(fileContents, []byte("\n")) // Remove trailing newline character. https://unix.stackexchange.com/questions/18743/whats-the-point-in-adding-a-new-line-to-the-end-of-a-file
	fileEncoded := b64.StdEncoding.EncodeToString(fileContents)
	return fileEncoded, nil
}
