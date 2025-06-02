package provider

import (
	"context"
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/access_profile"
	"github.com/itera-io/terraform-provider-taikun/taikun/alerting_profile"
	"github.com/itera-io/terraform-provider-taikun/taikun/app_instance"
	"github.com/itera-io/terraform-provider-taikun/taikun/backup_credential"
	"github.com/itera-io/terraform-provider-taikun/taikun/backup_policy"
	"github.com/itera-io/terraform-provider-taikun/taikun/billing"
	"github.com/itera-io/terraform-provider-taikun/taikun/catalog"
	"github.com/itera-io/terraform-provider-taikun/taikun/catalog_project_binding"
	"github.com/itera-io/terraform-provider-taikun/taikun/cc_aws"
	"github.com/itera-io/terraform-provider-taikun/taikun/cc_azure"
	"github.com/itera-io/terraform-provider-taikun/taikun/cc_gcp"
	"github.com/itera-io/terraform-provider-taikun/taikun/cc_openstack"
	"github.com/itera-io/terraform-provider-taikun/taikun/cc_proxmox"
	"github.com/itera-io/terraform-provider-taikun/taikun/cc_vsphere"
	"github.com/itera-io/terraform-provider-taikun/taikun/cc_zadara"
	"github.com/itera-io/terraform-provider-taikun/taikun/flavors"
	"github.com/itera-io/terraform-provider-taikun/taikun/kubeconfig"
	"github.com/itera-io/terraform-provider-taikun/taikun/kubernetes_profile"
	"github.com/itera-io/terraform-provider-taikun/taikun/organization"
	"github.com/itera-io/terraform-provider-taikun/taikun/policy_profile"
	"github.com/itera-io/terraform-provider-taikun/taikun/project"
	"github.com/itera-io/terraform-provider-taikun/taikun/repository"
	"github.com/itera-io/terraform-provider-taikun/taikun/showback"
	"github.com/itera-io/terraform-provider-taikun/taikun/slack"
	"github.com/itera-io/terraform-provider-taikun/taikun/standalone_profile"
	"github.com/itera-io/terraform-provider-taikun/taikun/user"
	"github.com/itera-io/terraform-provider-taikun/taikun/virtual_cluster"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
)

func init() {
	log.SetPrefix("[ERROR] ")

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
			"taikun_access_profile":              access_profile.DataSourceTaikunAccessProfile(),
			"taikun_access_profiles":             access_profile.DataSourceTaikunAccessProfiles(),
			"taikun_alerting_profile":            alerting_profile.DataSourceTaikunAlertingProfile(),
			"taikun_alerting_profiles":           alerting_profile.DataSourceTaikunAlertingProfiles(),
			"taikun_app_instance":                app_instance.DataSourceTaikunAppInstance(),
			"taikun_app_instances":               app_instance.DataSourceTaikunAppInstances(),
			"taikun_backup_credential":           backup_credential.DataSourceTaikunBackupCredential(),
			"taikun_backup_credentials":          backup_credential.DataSourceTaikunBackupCredentials(),
			"taikun_billing_credential":          billing.DataSourceTaikunBillingCredential(),
			"taikun_billing_credentials":         billing.DataSourceTaikunBillingCredentials(),
			"taikun_billing_rule":                billing.DataSourceTaikunBillingRule(),
			"taikun_billing_rules":               billing.DataSourceTaikunBillingRules(),
			"taikun_catalog":                     catalog.DataSourceTaikunCatalog(),
			"taikun_catalogs":                    catalog.DataSourceTaikunCatalogs(),
			"taikun_catalog_project_binding":     catalog_project_binding.DataSourceTaikunCatalogProjectBinding(),
			"taikun_catalog_project_bindings":    catalog_project_binding.DataSourceTaikunCatalogProjectBindings(),
			"taikun_cloud_credential_aws":        cc_aws.DataSourceTaikunCloudCredentialAWS(),
			"taikun_cloud_credential_azure":      cc_azure.DataSourceTaikunCloudCredentialAzure(),
			"taikun_cloud_credential_gcp":        cc_gcp.DataSourceTaikunCloudCredentialGCP(),
			"taikun_cloud_credential_openstack":  cc_openstack.DataSourceTaikunCloudCredentialOpenStack(),
			"taikun_cloud_credential_proxmox":    cc_proxmox.DataSourceTaikunCloudCredentialProxmox(),
			"taikun_cloud_credential_vsphere":    cc_vsphere.DataSourceTaikunCloudCredentialVsphere(),
			"taikun_cloud_credential_zadara":     cc_zadara.DataSourceTaikunCloudCredentialZadara(),
			"taikun_cloud_credentials_aws":       cc_aws.DataSourceTaikunCloudCredentialsAWS(),
			"taikun_cloud_credentials_azure":     cc_azure.DataSourceTaikunCloudCredentialsAzure(),
			"taikun_cloud_credentials_gcp":       cc_gcp.DataSourceTaikunCloudCredentialsGCP(),
			"taikun_cloud_credentials_openstack": cc_openstack.DataSourceTaikunCloudCredentialsOpenStack(),
			"taikun_cloud_credentials_proxmox":   cc_proxmox.DataSourceTaikunCloudCredentialsProxmox(),
			"taikun_cloud_credentials_vsphere":   cc_vsphere.DataSourceTaikunCloudCredentialsVsphere(),
			"taikun_cloud_credentials_zadara":    cc_zadara.DataSourceTaikunCloudCredentialsZadara(),
			"taikun_flavors":                     flavors.DataSourceTaikunFlavors(),
			"taikun_images_aws":                  cc_aws.DataSourceTaikunImagesAWS(),
			"taikun_images_azure":                cc_azure.DataSourceTaikunImagesAzure(),
			"taikun_images_gcp":                  cc_gcp.DataSourceTaikunImagesGCP(),
			"taikun_images_openstack":            cc_openstack.DataSourceTaikunImagesOpenStack(),
			"taikun_images_proxmox":              cc_proxmox.DataSourceTaikunImagesProxmox(),
			"taikun_images_vsphere":              cc_vsphere.DataSourceTaikunImagesVsphere(),
			"taikun_images_zadara":               cc_zadara.DataSourceTaikunImagesZadara(),
			"taikun_kubeconfig":                  kubeconfig.DataSourceTaikunKubeconfig(),
			"taikun_kubeconfigs":                 kubeconfig.DataSourceTaikunKubeconfigs(),
			"taikun_kubernetes_profile":          kubernetes_profile.DataSourceTaikunKubernetesProfile(),
			"taikun_kubernetes_profiles":         kubernetes_profile.DataSourceTaikunKubernetesProfiles(),
			"taikun_organization":                organization.DataSourceTaikunOrganization(),
			"taikun_organizations":               organization.DataSourceTaikunOrganizations(),
			"taikun_policy_profile":              policy_profile.DataSourceTaikunPolicyProfile(),
			"taikun_policy_profiles":             policy_profile.DataSourceTaikunPolicyProfiles(),
			"taikun_project":                     project.DataSourceTaikunProject(),
			"taikun_projects":                    project.DataSourceTaikunProjects(),
			"taikun_repository":                  repository.DataSourceTaikunRespository(),
			"taikun_repositories":                repository.DataSourceTaikunRepositories(),
			"taikun_showback_credential":         showback.DataSourceTaikunShowbackCredential(),
			"taikun_showback_credentials":        showback.DataSourceTaikunShowbackCredentials(),
			"taikun_showback_rule":               showback.DataSourceTaikunShowbackRule(),
			"taikun_showback_rules":              showback.DataSourceTaikunShowbackRules(),
			"taikun_slack_configuration":         slack.DataSourceTaikunSlackConfiguration(),
			"taikun_slack_configurations":        slack.DataSourceTaikunSlackConfigurations(),
			"taikun_standalone_profile":          standalone_profile.DataSourceTaikunStandaloneProfile(),
			"taikun_standalone_profiles":         standalone_profile.DataSourceTaikunStandaloneProfiles(),
			"taikun_user":                        user.DataSourceTaikunUser(),
			"taikun_users":                       user.DataSourceTaikunUsers(),
			"taikun_virtual_cluster":             virtual_cluster.DataSourceTaikunVirtualCluster(),
			"taikun_virtual_clusters":            virtual_cluster.DataSourceTaikunVirtualClusters(),
			//"taikun_images":                      taikun.dataSourceTaikunImages(), // DEPRECATED
		},
		ResourcesMap: map[string]*schema.Resource{
			"taikun_access_profile":                       access_profile.ResourceTaikunAccessProfile(),
			"taikun_alerting_profile":                     alerting_profile.ResourceTaikunAlertingProfile(),
			"taikun_app_instance":                         app_instance.ResourceTaikunAppInstance(),
			"taikun_backup_credential":                    backup_credential.ResourceTaikunBackupCredential(),
			"taikun_backup_policy":                        backup_policy.ResourceTaikunBackupPolicy(),
			"taikun_billing_credential":                   billing.ResourceTaikunBillingCredential(),
			"taikun_billing_rule":                         billing.ResourceTaikunBillingRule(),
			"taikun_catalog":                              catalog.ResourceTaikunCatalog(),
			"taikun_catalog_project_binding":              catalog_project_binding.ResourceTaikunCatalogProjectBinding(),
			"taikun_cloud_credential_aws":                 cc_aws.ResourceTaikunCloudCredentialAWS(),
			"taikun_cloud_credential_azure":               cc_azure.ResourceTaikunCloudCredentialAzure(),
			"taikun_cloud_credential_gcp":                 cc_gcp.ResourceTaikunCloudCredentialGCP(),
			"taikun_cloud_credential_openstack":           cc_openstack.ResourceTaikunCloudCredentialOpenStack(),
			"taikun_cloud_credential_proxmox":             cc_proxmox.ResourceTaikunCloudCredentialProxmox(),
			"taikun_cloud_credential_vsphere":             cc_vsphere.ResourceTaikunCloudCredentialVsphere(),
			"taikun_cloud_credential_zadara":              cc_zadara.ResourceTaikunCloudCredentialZadara(),
			"taikun_kubeconfig":                           kubeconfig.ResourceTaikunKubeconfig(),
			"taikun_kubernetes_profile":                   kubernetes_profile.ResourceTaikunKubernetesProfile(),
			"taikun_organization_billing_rule_attachment": organization.ResourceTaikunOrganizationBillingRuleAttachment(),
			"taikun_organization":                         organization.ResourceTaikunOrganization(),
			"taikun_policy_profile":                       policy_profile.ResourceTaikunPolicyProfile(),
			"taikun_project":                              project.ResourceTaikunProject(),
			"taikun_project_user_attachment":              project.ResourceTaikunProjectUserAttachment(),
			"taikun_repository":                           repository.ResourceTaikunRepository(),
			"taikun_showback_credential":                  showback.ResourceTaikunShowbackCredential(),
			"taikun_showback_rule":                        showback.ResourceTaikunShowbackRule(),
			"taikun_slack_configuration":                  slack.ResourceTaikunSlackConfiguration(),
			"taikun_standalone_profile":                   standalone_profile.ResourceTaikunStandaloneProfile(),
			"taikun_user":                                 user.ResourceTaikunUser(),
			"taikun_virtual_cluster":                      virtual_cluster.ResourceTaikunVirtualCluster(),
			//"taikun_cloud_credential":                     taikun.resourceTaikunCloudCredential(), // DEPRECATED

		},
		Schema: map[string]*schema.Schema{
			"api_host": {
				Type:         schema.TypeString,
				Description:  "Custom Taikun API host. Can be set with TAIKUN_API_HOST.",
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("TAIKUN_API_HOST", "api.taikun.cloud"),
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"email": {
				Type:          schema.TypeString,
				Description:   "Taikun email. Can be set with TAIKUN_EMAIL",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_EMAIL", nil),
				ConflictsWith: []string{"keycloak_email", "access_key"},
				RequiredWith:  []string{"password"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"password": {
				Type:          schema.TypeString,
				Description:   "Taikun password. Can be set with TAIKUN_PASSWORD.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_PASSWORD", nil),
				ConflictsWith: []string{"keycloak_password", "access_key"},
				RequiredWith:  []string{"email"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"access_key": {
				Type:          schema.TypeString,
				Description:   "Taikun access key. Can be set with TAIKUN_ACCESS_KEY.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_ACCESS_KEY", nil),
				ConflictsWith: []string{"email", "keycloak_email"},
				RequiredWith:  []string{"secret_key"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"secret_key": {
				Type:          schema.TypeString,
				Description:   "Taikun secret key. Can be set with TAIKUN_SECRET_KEY.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_SECRET_KEY", nil),
				ConflictsWith: []string{"email", "keycloak_email"},
				RequiredWith:  []string{"access_key"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"keycloak_email": {
				Type:          schema.TypeString,
				Description:   "Taikun Keycloak email. Can be set with TAIKUN_KEYCLOAK_EMAIL.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_KEYCLOAK_EMAIL", nil),
				ConflictsWith: []string{"email", "access_key"},
				RequiredWith:  []string{"keycloak_password"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"keycloak_password": {
				Type:          schema.TypeString,
				Description:   "Taikun Keycloak password. Can be set with TAIKUN_KEYCLOAK_PASSWORD.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_KEYCLOAK_PASSWORD", nil),
				ConflictsWith: []string{"email", "access_key"},
				RequiredWith:  []string{"keycloak_email"},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
		},
		ConfigureContextFunc: configureContextFunc,
	}
}

func configureContextFunc(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var (
		email     string
		password  string
		accessKey string
		secretKey string
		authMode  string
	)

	// 1. Try Keycloak
	if rawEmail, ok := d.GetOk("keycloak_email"); ok {
		email, _ = rawEmail.(string)
		rawPassword := d.Get("keycloak_password")
		password, _ = rawPassword.(string)
		authMode = "keycloak"
	} else if rawEmail, ok := d.GetOk("email"); ok {
		// 2. Try Username/Password
		email, _ = rawEmail.(string)
		rawPassword := d.Get("password")
		password, _ = rawPassword.(string)
		authMode = ""
	} else if rawAccessKey, ok := d.GetOk("access_key"); ok {
		// 3. Try Access Key mode
		accessKey, _ = rawAccessKey.(string)
		rawSecretKey := d.Get("secret_key")
		secretKey, _ = rawSecretKey.(string)
		authMode = "token"
	} else {
		return nil, diag.Errorf("You must define credentials using either keycloak_email, email/password, or access_key/secret_key")
	}

	apiHost, _ := d.Get("api_host").(string)

	client := tk.NewClientFromCredentials(email, password, accessKey, secretKey, authMode, apiHost)
	return client, nil
}
