package taikun

import (
	tkcore "github.com/chnyda/taikungoclient/client"
)

func stringPtr(s string) *string {
	return &s
}

func boolPtr(s bool) *bool {
	return &s
}

func newNullableFloat64(s float64) tkcore.NullableFloat64 {
	res := tkcore.NullableFloat64{}
	res.Set(&s)
	return res
}
