package taikun

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			defaultString := fmt.Sprint(s.Default)
			if len(defaultString) == 0 {
				defaultString = " "
			}

			desc += fmt.Sprintf(" Defaults to `%s`.", defaultString)
		}
		if s.Deprecated != "" {
			desc += " " + s.Deprecated
		}
		if len(s.ConflictsWith) != 0 {
			desc += " Conflicts with: "
			for index, element := range s.ConflictsWith {
				desc += fmt.Sprintf("`%s`", element)
				if index != len(s.ConflictsWith)-1 {
					desc += ", "
				}
			}
			desc += "."
		}
		if len(s.RequiredWith) != 0 {
			desc += " Required with: "
			for index, element := range s.RequiredWith {
				desc += fmt.Sprintf("`%s`", element)
				if index != len(s.RequiredWith)-1 {
					desc += ", "
				}
			}
			desc += "."
		}
		return strings.TrimSpace(desc)
	}
}

var ApiVersion = "1"

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"taikun_access_profile":              dataSourceTaikunAccessProfile(),
			"taikun_access_profiles":             dataSourceTaikunAccessProfiles(),
			"taikun_alerting_profile":            dataSourceTaikunAlertingProfile(),
			"taikun_alerting_profiles":           dataSourceTaikunAlertingProfiles(),
			"taikun_backup_credential":           dataSourceTaikunBackupCredential(),
			"taikun_backup_credentials":          dataSourceTaikunBackupCredentials(),
			"taikun_billing_credential":          dataSourceTaikunBillingCredential(),
			"taikun_billing_credentials":         dataSourceTaikunBillingCredentials(),
			"taikun_billing_rule":                dataSourceTaikunBillingRule(),
			"taikun_billing_rules":               dataSourceTaikunBillingRules(),
			"taikun_cloud_credential_aws":        dataSourceTaikunCloudCredentialAWS(),
			"taikun_cloud_credential_azure":      dataSourceTaikunCloudCredentialAzure(),
			"taikun_cloud_credential_openstack":  dataSourceTaikunCloudCredentialOpenStack(),
			"taikun_cloud_credentials_aws":       dataSourceTaikunCloudCredentialsAWS(),
			"taikun_cloud_credentials_azure":     dataSourceTaikunCloudCredentialsAzure(),
			"taikun_cloud_credentials_openstack": dataSourceTaikunCloudCredentialsOpenStack(),
			"taikun_flavors":                     dataSourceTaikunFlavors(),
			"taikun_kubeconfig":                  dataSourceTaikunKubeconfig(),
			"taikun_kubeconfigs":                 dataSourceTaikunKubeconfigs(),
			"taikun_kubernetes_profile":          dataSourceTaikunKubernetesProfile(),
			"taikun_kubernetes_profiles":         dataSourceTaikunKubernetesProfiles(),
			"taikun_opa_profile":                 dataSourceTaikunOPAProfile(),
			"taikun_opa_profiles":                dataSourceTaikunOPAProfiles(),
			"taikun_organization":                dataSourceTaikunOrganization(),
			"taikun_organizations":               dataSourceTaikunOrganizations(),
			"taikun_project":                     dataSourceTaikunProject(),
			"taikun_projects":                    dataSourceTaikunProjects(),
			"taikun_showback_credential":         dataSourceTaikunShowbackCredential(),
			"taikun_showback_credentials":        dataSourceTaikunShowbackCredentials(),
			"taikun_showback_rule":               dataSourceTaikunShowbackRule(),
			"taikun_showback_rules":              dataSourceTaikunShowbackRules(),
			"taikun_slack_configuration":         dataSourceTaikunSlackConfiguration(),
			"taikun_slack_configurations":        dataSourceTaikunSlackConfigurations(),
			"taikun_user":                        dataSourceTaikunUser(),
			"taikun_users":                       dataSourceTaikunUsers(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"taikun_access_profile":                       resourceTaikunAccessProfile(),
			"taikun_alerting_profile":                     resourceTaikunAlertingProfile(),
			"taikun_backup_credential":                    resourceTaikunBackupCredential(),
			"taikun_backup_policy":                        resourceTaikunBackupPolicy(),
			"taikun_billing_credential":                   resourceTaikunBillingCredential(),
			"taikun_billing_rule":                         resourceTaikunBillingRule(),
			"taikun_cloud_credential_aws":                 resourceTaikunCloudCredentialAWS(),
			"taikun_cloud_credential_azure":               resourceTaikunCloudCredentialAzure(),
			"taikun_cloud_credential_openstack":           resourceTaikunCloudCredentialOpenStack(),
			"taikun_kubeconfig":                           resourceTaikunKubeconfig(),
			"taikun_kubernetes_profile":                   resourceTaikunKubernetesProfile(),
			"taikun_opa_profile":                          resourceTaikunOPAProfile(),
			"taikun_organization_billing_rule_attachment": resourceTaikunOrganizationBillingRuleAttachment(),
			"taikun_organization":                         resourceTaikunOrganization(),
			"taikun_project":                              resourceTaikunProject(),
			"taikun_project_user_attachment":              resourceTaikunProjectUserAttachment(),
			"taikun_showback_credential":                  resourceTaikunShowbackCredential(),
			"taikun_showback_rule":                        resourceTaikunShowbackRule(),
			"taikun_slack_configuration":                  resourceTaikunSlackConfiguration(),
			"taikun_user":                                 resourceTaikunUser(),
		},
		Schema: map[string]*schema.Schema{
			"api_host": {
				Type:         schema.TypeString,
				Description:  "Custom Taikun API host.",
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("TAIKUN_API_HOST", nil),
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"email": {
				Type:          schema.TypeString,
				Description:   "Taikun email.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_EMAIL", nil),
				ConflictsWith: []string{"keycloak_email"},
				RequiredWith:  []string{"password"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"password": {
				Type:          schema.TypeString,
				Description:   "Taikun password.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_PASSWORD", nil),
				ConflictsWith: []string{"keycloak_password"},
				RequiredWith:  []string{"email"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"keycloak_email": {
				Type:          schema.TypeString,
				Description:   "Taikun Keycloak email.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_KEYCLOAK_EMAIL", nil),
				ConflictsWith: []string{"email"},
				RequiredWith:  []string{"keycloak_password"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"keycloak_password": {
				Type:          schema.TypeString,
				Description:   "Taikun Keycloak password.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_KEYCLOAK_PASSWORD", nil),
				ConflictsWith: []string{"password"},
				RequiredWith:  []string{"keycloak_email"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
		},
		ConfigureContextFunc: configureContextFunc,
	}
}

func configureContextFunc(_ context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {

	email, keycloakEnabled := data.GetOk("keycloak_email")
	password := data.Get("keycloak_password")

	if !keycloakEnabled {
		email = data.Get("email")
		password = data.Get("password")
	}

	if email == "" || password == "" {
		return nil, diag.Errorf("You must define an email and a password")
	}

	transportConfig := client.DefaultTransportConfig()

	if apiHost, apiHostIsSet := data.GetOk("api_host"); apiHostIsSet {
		transportConfig = transportConfig.WithHost(apiHost.(string))
	}

	return &apiClient{
		client:              client.NewHTTPClientWithConfig(nil, transportConfig),
		email:               email.(string),
		password:            password.(string),
		useKeycloakEndpoint: keycloakEnabled,
	}, nil
}
