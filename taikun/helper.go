package taikun

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/itera-io/taikungoclient/models"
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

func stringIsInt(i interface{}, path cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.FromErr(path.NewErrorf("expected type to be string"))
	}

	_, err := strconv.Atoi(v)
	if err != nil {
		return diag.FromErr(path.NewErrorf("expected an int inside a string"))
	}

	return nil
}

func randomTestName() string {
	return randomName(testNamePrefix, 15)
}

func randomName(prefix string, length int) string {
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(length))
}

func randomString() string {
	rand.Seed(time.Now().UnixNano())
	return acctest.RandString(rand.Int()%10 + 10)
}

func randomURL() string {
	return fmt.Sprintf("https://%s.%s.example", randomString(), randomString())
}

func randomEmail() string {
	return fmt.Sprintf("%s@%s.example", randomString(), randomString())
}

func randomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()%2 == 0
}

// Return an integer in the range [0; maxInt[
func randomInt(maxInt int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int() % maxInt
}

func getLockMode(locked bool) string {
	if locked {
		return "lock"
	}
	return "unlock"
}

func getPrometheusType(prometheusType string) models.PrometheusType {
	if prometheusType == "Count" {
		return 100
	}
	return 200
}

func getShowbackType(showbackType string) models.ShowbackType {
	if showbackType == "General" {
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

func getUserRole(role string) models.UserRole {
	if role == "User" {
		return 400
	}
	// Manager
	return 200
}

func parseLoadBalancingSolution(loadBalancingSolution string) (bool, bool) {
	if loadBalancingSolution == "Octavia" {
		return true, false
	} else if loadBalancingSolution == "Taikun" {
		return false, true
	}
	return false, false
}

func getSlackConfigurationType(configType string) models.SlackType {
	if configType == "Alert" {
		return 100
	}
	return 200 // General
}

func getAlertingProfileReminder(reminder string) models.AlertingReminder {
	switch reminder {
	case "HalfHour":
		return 100
	case "Hourly":
		return 200
	case "Daily":
		return 300
	default: // "None"
		return -1
	}
}
