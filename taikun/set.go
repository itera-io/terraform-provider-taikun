package taikun

import (
	"hash/crc32"
	"strconv"
)

func hashString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func hashAttributes(keys ...string) func(v interface{}) int {
	return func(v interface{}) int {
		stringToHash := ""
		set := v.(map[string]interface{})

		for _, key := range keys {

			if v, ok := set[key].(string); ok {
				stringToHash += v
			}

			if v, ok := set[key].(int); ok {
				stringToHash += strconv.Itoa(v)
			}

			if labelsList, ok := set[key].([]interface{}); ok {
				for _, e := range labelsList {
					if label, ok2 := e.(map[string]interface{}); ok2 {
						stringToHash += label["key"].(string)
						stringToHash += label["value"].(string)
					}
				}
			}
		}
		return hashString(stringToHash)
	}
}
