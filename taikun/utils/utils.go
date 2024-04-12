package utils

import (
	tkcore "github.com/itera-io/taikungoclient/client"
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
