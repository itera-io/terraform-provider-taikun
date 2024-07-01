package utils

import (
	"context"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
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

//func init() {
//	rand.Seed(time.Now().UnixNano())
//}

func SetResourceDataFromMap(d *schema.ResourceData, m map[string]interface{}) error {
	for key, value := range m {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("unable to set `%s` attribute: %s", key, err)
		}
	}
	return nil
}

func Atoi32(str string) (int32, error) {
	res, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(res), nil
}

func I32toa(x int32) string {
	return strconv.FormatInt(int64(x), 10)
}

func GibiByteToByte(gibiBytes int) float64 {
	return float64(1073741824) * float64(gibiBytes)
}

func ByteToGibiByte(x float64) int {
	return int(x / 1073741824)
}

func StringIsInt(i interface{}, path cty.Path) diag.Diagnostics {
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

func StringIsEmail(i interface{}, path cty.Path) diag.Diagnostics {
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

func StringIsFilePath(i interface{}, path cty.Path) diag.Diagnostics {
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

func StringIsCron(i interface{}, path cty.Path) diag.Diagnostics {
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

func DateToDateTime(date string) strfmt.DateTime {
	time, _ := time.Parse(time.RFC3339, dateToRfc3339DateTime(date))
	return strfmt.DateTime(time)
}

func dateToRfc3339DateTime(date string) string {
	return date[6:10] + "-" + date[3:5] + "-" + date[0:2] + "T00:00:00Z"
}

func Rfc3339DateTimeToDate(date string) string {
	if date == "" {
		return date
	}
	return date[8:10] + "/" + date[5:7] + "/" + date[0:4]
}

func StringIsDate(i interface{}, path cty.Path) diag.Diagnostics {
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

func StringIsUUID(i interface{}, path cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.FromErr(path.NewErrorf("expected type to be string"))
	}

	if _, err := uuid.ParseUUID(v); err != nil {
		return diag.FromErr(path.NewErrorf("expected a valid UUID, got %v", v))
	}

	return nil
}

func StringLenBetween(min int, max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.FromErr(path.NewErrorf("expected type to be string"))
		}

		if len(v) < min || len(v) > max {
			return diag.FromErr(fmt.Errorf("expected length of %d to be in the range (%d - %d), got %s", len(v), min, max, v))
		}

		return nil
	}
}

func ResourceGetStringList(data interface{}) []string {
	rawList := data.([]interface{})
	result := make([]string, 0)
	for _, e := range rawList {
		result = append(result, e.(string))
	}
	return result
}

func RandomTestName() string {
	return randomName(testNamePrefix, 15)
}

func ShortRandomTestName() string {
	return randomName(testNamePrefix, 5)
}

func randomName(prefix string, length int) string {
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(length))
}

func RandomString() string {
	return fmt.Sprintf("%s%s", "tf", acctest.RandString(rand.Int()%6+7)) // Taikun can have problems with strings starting with numbers
}

func RandomURL() string {
	return fmt.Sprintf("https://%s.%s.example", RandomString(), RandomString())
}

func RandomEmail() string {
	return fmt.Sprintf("%s@mailinator.com", RandomString())
}

func RandomBool() bool {
	return rand.Int()%2 == 0
}

// Return an integer in the range [0; maxInt[
func RandomInt(maxInt int) int {
	return rand.Int() % maxInt
}

func GetLockMode(locked bool) string {
	if locked {
		return "lock"
	}
	return "unlock"
}

func GetEPrometheusType(prometheusType string) tkshowback.EPrometheusType {
	return getPrometheusTypeInt(prometheusType)
}

func GetPrometheusType(prometheusType string) tkcore.PrometheusType {
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

func GetShowbackType(showbackType string) tkshowback.EShowbackType {
	if showbackType == "General" {
		return tkshowback.ESHOWBACKTYPE_GENERAL
	}
	return tkshowback.ESHOWBACKTYPE_EXTERNAL // External
}

const (
	loadBalancerOctavia = "Octavia"
	LoadBalancerTaikun  = "Taikun"
	loadBalancerNone    = "None"
)

func GetLoadBalancingSolution(octaviaEnabled bool, taikunLBEnabled bool) string {
	if octaviaEnabled {
		return loadBalancerOctavia
	} else if taikunLBEnabled {
		return LoadBalancerTaikun
	}
	return loadBalancerNone
}

func ParseLoadBalancingSolution(loadBalancingSolution string) (octaviaEnabled bool, taikunLBEnabled bool) {
	if loadBalancingSolution == loadBalancerOctavia {
		return true, false
	} else if loadBalancingSolution == LoadBalancerTaikun {
		return false, true
	}
	return false, false
}

func GetAlertingIntegrationType(integrationType string) tkcore.AlertingIntegrationType {
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

func GetKubeconfigRoleID(role string) int32 {
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
	CloudTypeOpenStack = "OpenStack"
	cloudTypeGCP       = "GCP"
	cloudTypeProxmox   = "Proxmox"
)

func GetSecurityGroupProtocol(protocol string) tkcore.SecurityGroupProtocol {
	switch strings.ToUpper(protocol) {
	case "ICMP":
		return tkcore.SECURITYGROUPPROTOCOL_ICMP
	case "TCP":
		return tkcore.SECURITYGROUPPROTOCOL_TCP
	default: // UDP
		return tkcore.SECURITYGROUPPROTOCOL_UDP
	}
}

func SetResourceDataId(d *schema.ResourceData, id int32) {
	idAsString := strconv.FormatInt(int64(id), 10)
	d.SetId(idAsString)
}

func ContinentShorthand(continent string) string {
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

// This function is used when we have an API parameter PARAM with the following behavior
// - We set PARAM to "x" and API returns PARAM with "x" (there is no change)
// - We do not set PARAM. API figures out some value "y" and returns PARAM with value "y" (Terraform must ignore the change)
func IgnoreChangeFromEmpty(k string, old string, new string, d *schema.ResourceData) bool {
	// First apply, we did not specify the PARAM. Don't supress diff.
	if old == "" && new == "" {
		return false
	}

	// Second apply, we did not specify a PARAM. Supress diff.
	// There used to be a PARAM, but now we specified none. Supress diff.
	if old != "" && new == "" {
		return true
	}

	// The PARAM was changed in .tf file. Don't supress diff.
	if old != "" && new != "" {
		return false
	}

	// Else, Don't supress diff.
	return false
}

func GetLastCharacter(zoneString string) string {
	if len(zoneString) == 0 {
		return ""
	}
	return zoneString[len(zoneString)-1:]
}

// Proxmox storage options in kubernetes profile and Project are different from what we send while creating k8s servers
func GetProxmoxStorageStringForServer(projectID int32, apiClient *tk.Client) (string, error) {
	data, response, err := apiClient.Client.ServersAPI.ServersDetails(context.TODO(), projectID).Execute()
	if err != nil {
		return "", tk.CreateError(response, err)
	}

	kubernetesProfile := data.GetProject()
	proxmoxStorageString := kubernetesProfile.GetProxmoxStorage()
	switch proxmoxStorageString {
	case "NFS":
		return "NFS", nil
	case "OpenEBS":
		return "STORAGE", nil
	default:
		return "", fmt.Errorf("Parsed an unrecognised Proxmox Storage type from Kubernetes Profile for this project.")
	}
}
