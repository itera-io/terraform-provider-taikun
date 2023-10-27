package taikun

import (
	"fmt"
	tkcore "github.com/itera-io/taikungoclient/client"
	tkshowback "github.com/itera-io/taikungoclient/showbackclient"
	"math/rand"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/robfig/cron/v3"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

const testNamePrefix = "tf-acc-test-"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func setResourceDataFromMap(d *schema.ResourceData, m map[string]interface{}) error {
	for key, value := range m {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("unable to set `%s` attribute: %s", key, err)
		}
	}
	return nil
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

func gibiByteToMebiByte(x int32) int32 {
	return x * 1024
}

func mebiByteToGibiByte(x int64) int32 {
	var kibi int64 = 1024
	return int32(x / kibi)
}

func gibiByteToByte(x int) int64 {
	return int64(1073741824 * x)
}

func byteToGibiByte(x int64) int64 {
	return x / 1073741824
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

func stringIsFilePath(i interface{}, path cty.Path) diag.Diagnostics {
	filePath, ok := i.(string)
	if !ok {
		return diag.FromErr(path.NewErrorf("expected type to be string"))
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return diag.FromErr(path.NewError(err))
	}

	if fileInfo.IsDir() {
		return diag.FromErr(path.NewErrorf("expected a file, not a directory"))
	}

	return nil
}

func stringIsCron(i interface{}, path cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.FromErr(path.NewErrorf("expected type to be string"))
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := parser.Parse(v); err != nil {
		return diag.FromErr(path.NewErrorf("expected a valid cron period"))
	}

	return nil
}

func dateToDateTime(date string) strfmt.DateTime {
	time, _ := time.Parse(time.RFC3339, dateToRfc3339DateTime(date))
	return strfmt.DateTime(time)
}

func dateToRfc3339DateTime(date string) string {
	return date[6:10] + "-" + date[3:5] + "-" + date[0:2] + "T00:00:00Z"
}

func rfc3339DateTimeToDate(date string) string {
	if date == "" {
		return date
	}
	return date[8:10] + "/" + date[5:7] + "/" + date[0:4]
}

func stringIsDate(i interface{}, path cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.FromErr(path.NewErrorf("expected type to be string"))
	}

	if len(v) != 10 {
		return diag.FromErr(path.NewErrorf("expected a valid date in the format: 'dd/mm/yyyy'"))
	}

	if _, err := time.Parse(time.RFC3339, dateToRfc3339DateTime(v)); len(v) != 10 || err != nil {
		return diag.FromErr(path.NewErrorf("expected a valid date in the format: 'dd/mm/yyyy'"))
	}

	return nil
}

func stringIsUUID(i interface{}, path cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.FromErr(path.NewErrorf("expected type to be string"))
	}

	if _, err := uuid.ParseUUID(v); err != nil {
		return diag.FromErr(path.NewErrorf("expected a valid UUID, got %v", v))
	}

	return nil
}

func resourceGetStringList(data interface{}) []string {
	rawList := data.([]interface{})
	result := make([]string, 0)
	for _, e := range rawList {
		result = append(result, e.(string))
	}
	return result
}

func randomTestName() string {
	return randomName(testNamePrefix, 15)
}

func shortRandomTestName() string {
	return randomName(testNamePrefix, 5)
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
	return fmt.Sprintf("%s@mailinator.com", randomString())
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

func getEPrometheusType(prometheusType string) tkshowback.EPrometheusType {
	return getPrometheusTypeInt(prometheusType)
}

func getPrometheusType(prometheusType string) tkcore.PrometheusType {
	if prometheusType == "Count" {
		return tkcore.PROMETHEUSTYPE_COUNT
	}
	return tkcore.PROMETHEUSTYPE_SUM
}

func getPrometheusTypeInt(prometheusType string) tkshowback.EPrometheusType {
	if prometheusType == "Count" {
		return tkshowback.EPROMETHEUSTYPE_COUNT
	}
	return tkshowback.EPROMETHEUSTYPE_SUM // Sum
}

func getShowbackType(showbackType string) tkshowback.EShowbackType {
	if showbackType == "General" {
		return tkshowback.ESHOWBACKTYPE_GENERAL
	}
	return tkshowback.ESHOWBACKTYPE_EXTERNAL // External
}

const (
	loadBalancerOctavia = "Octavia"
	loadBalancerTaikun  = "Taikun"
	loadBalancerNone    = "None"
)

func getLoadBalancingSolution(octaviaEnabled bool, taikunLBEnabled bool) string {
	if octaviaEnabled {
		return loadBalancerOctavia
	} else if taikunLBEnabled {
		return loadBalancerTaikun
	}
	return loadBalancerNone
}

func parseLoadBalancingSolution(loadBalancingSolution string) (octaviaEnabled bool, taikunLBEnabled bool) {
	if loadBalancingSolution == loadBalancerOctavia {
		return true, false
	} else if loadBalancingSolution == loadBalancerTaikun {
		return false, true
	}
	return false, false
}

func getAlertingIntegrationType(integrationType string) tkcore.AlertingIntegrationType {
	switch integrationType {
	case "Opsgenie":
		return tkcore.ALERTINGINTEGRATIONTYPE_OPSGENIE
	case "Pagerduty":
		return tkcore.ALERTINGINTEGRATIONTYPE_PAGERDUTY
	case "Splunk":
		return tkcore.ALERTINGINTEGRATIONTYPE_SPLUNK
	default: // "MicrosoftTeams"
		return tkcore.ALERTINGINTEGRATIONTYPE_MICROSOFT_TEAMS
	}
}

func getKubeconfigRoleID(role string) int32 {
	switch role {
	case "cluster-admin":
		return 1
	case "admin":
		return 2
	case "edit":
		return 3
	default: // view
		return 4
	}
}

const (
	cloudTypeAWS       = "AWS"
	cloudTypeAzure     = "Azure"
	cloudTypeOpenStack = "OpenStack"
	cloudTypeGCP       = "GCP"
)

func getSecurityGroupProtocol(protocol string) tkcore.SecurityGroupProtocol {
	switch strings.ToUpper(protocol) {
	case "ICMP":
		return tkcore.SECURITYGROUPPROTOCOL_ICMP
	case "TCP":
		return tkcore.SECURITYGROUPPROTOCOL_TCP
	default: // UDP
		return tkcore.SECURITYGROUPPROTOCOL_UDP
	}
}

func setResourceDataId(d *schema.ResourceData, id int32) {
	idAsString := strconv.FormatInt(int64(id), 10)
	d.SetId(idAsString)
}

func continentShorthand(continent string) string {
	continent = strings.ToLower(continent)
	if continent == "europe" {
		return "eu"
	}
	if continent == "asia" {
		return "as"
	}
	if continent == "america" {
		return "us"
	}
	return ""
}
