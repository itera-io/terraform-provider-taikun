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
	if filePath == "" {
		return "", nil
	}
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	fileContents = bytes.TrimSuffix(fileContents, []byte("\n")) // Remove trailing newline character. https://unix.stackexchange.com/questions/18743/whats-the-point-in-adding-a-new-line-to-the-end-of-a-file
	fileEncoded := b64.StdEncoding.EncodeToString(fileContents)
	return fileEncoded, nil
}

// Slice of strings to slice of int32
func SliceOfSTringsToSliceOfInt32(listOfStrings []interface{}) ([]int32, error) {
	var sliceOfInt32 []int32
	for _, elementString := range listOfStrings {
		elementInt32, err := Atoi32(elementString.(string))
		if err != nil {
			return nil, err
		}
		sliceOfInt32 = append(sliceOfInt32, elementInt32)
	}
	return sliceOfInt32, nil
}
