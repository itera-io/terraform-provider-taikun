package taikun

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

const testNamePrefix = "tf-acc-test-"

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

func stringIsLowercase(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}
	if strings.ToLower(v) != v {
		return nil, []error{fmt.Errorf("expected %q to be lowercase", k)}
	}
	return nil, nil
}

func stringIsInt(i interface{}, k string) ([]string, []error) {
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

func randomTestName() string {
	return randomName(testNamePrefix, 10)
}

func randomName(prefix string, length int) string {
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(length))
}

func randomString() string {
	return acctest.RandString(rand.Int()%10 + 10)
}

func getLockMode(locked bool) string {
	if locked {
		return "lock"
	}
	return "unlock"
}

func getPrometheusType(prometheusType string) int32 {
	if prometheusType == "Count" {
		return 100
	}
	return 200
}

func getLoadBalancingSolution(octaviaEnabled bool, taikunLBEnabled bool) string {
	if octaviaEnabled {
		return "Octavia"
	} else if taikunLBEnabled {
		return "Taikun"
	}
	return "None"
}

func parseLoadBalancingSolution(loadBalancingSolution string) (bool, bool) {
	if loadBalancingSolution == "Octavia" {
		return true, false
	} else if loadBalancingSolution == "Taikun" {
		return false, true
	}
	return false, false
}

func getSlackConfigurationType(configType string) int32 {
	if configType == "Alert" {
		return 100
	}
	return 200 // General
}
