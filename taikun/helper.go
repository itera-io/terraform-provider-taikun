package taikun

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/itera-io/taikungoclient/models"
)

const testNamePrefix = "tf-acc-test-"

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

func jsonNumberAsFloatToInt32(value json.Number) int32 {
	x, _ := strconv.ParseFloat(string(value), 32)
	return int32(x)
}

func gibiByteToMebiByte(x int32) int32 {
	return x * 1024
}

func mebiByteToGibiByte(x int32) int32 {
	return x / 1024
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

func stringIsEmail(i interface{}, path cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.FromErr(path.NewErrorf("expected type to be string"))
	}

	_, err := mail.ParseAddress(v)
	if err != nil {
		return diag.FromErr(path.NewErrorf("expected an email"))
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
	return acctest.RandString(rand.Int()%10 + 10)
}

func randomURL() string {
	return fmt.Sprintf("https://%s.%s.example", randomString(), randomString())
}

func randomEmail() string {
	return fmt.Sprintf("%s@%s.example", randomString(), randomString())
}

func randomBool() bool {
	return rand.Int()%2 == 0
}

// Return an integer in the range [0; maxInt[
func randomInt(maxInt int) int {
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

func getAlertingIntegrationType(integrationType string) models.AlertingIntegrationType {
	switch integrationType {
	case "Opsgenie":
		return 100
	case "Pagerduty":
		return 200
	case "Splunk":
		return 300
	default: // "MicrosoftTeams"
		return 400
	}
}

func getAWSRegion(region string) models.AwsRegion {
	switch region {
	case "us-east-1":
		return 1
	case "us-east-2":
		return 2
	case "us-west-1":
		return 3
	case "us-west-2":
		return 4
	case "eu-north-1":
		return 5
	case "eu-west-1":
		return 6
	case "eu-west-2":
		return 7
	case "eu-west-3":
		return 8
	case "eu-central-1":
		return 9
	case "eu-south-1":
		return 10
	case "ap-east-1":
		return 11
	case "ap-northeast-1":
		return 12
	case "ap-northeast-2":
		return 13
	case "ap-northeast-3":
		return 14
	case "ap-south-1":
		return 15
	case "ap-southeast-1":
		return 16
	case "ap-southeast-2":
		return 17
	case "sa-east-1":
		return 18
	case "us-gov-east-1":
		return 19
	case "us-gov-west-1":
		return 20
	case "cn-north-1":
		return 21
	case "cn-northwest-1":
		return 22
	case "ca-central-1":
		return 23
	case "me-south-1":
		return 24
	default: // af-south-1
		return 25
	}
}
