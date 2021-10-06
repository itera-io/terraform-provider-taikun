package taikun

import (
	"fmt"
	"strconv"
)

func atoi32(str string) (int32, error) {
	res, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(res), nil
}

func i32toa(x int32) string {
	return strconv.FormatInt(int64(x), 10)
}

func StringIsInt(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	_, err := strconv.Atoi(v)
	if err != nil {
		return nil, []error{fmt.Errorf("expected %q to be an int inside a string", k)}
	}

	return nil, nil
}
