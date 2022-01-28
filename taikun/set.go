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

			if v, ok := set[key].(bool); ok {
				stringToHash += strconv.FormatBool(v)
			}

			if list, ok := set[key].([]interface{}); ok {
				for _, e := range list {
					if str, ok2 := e.(string); ok2 {
						stringToHash += str
					}

					if label, ok2 := e.(map[string]interface{}); ok2 {
						for _, value := range label {
							if str, ok3 := value.(string); ok3 {
								stringToHash += str
							}
						}
					}
				}
			}

			if amap, ok := set[key].(map[string]interface{}); ok {
				for _, value := range amap {
					if str, ok2 := value.(string); ok2 {
						stringToHash += str
					}
				}
			}
		}
		return hashString(stringToHash)
	}
}
