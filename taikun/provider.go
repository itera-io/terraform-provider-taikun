package taikun

import (
	"context"
	"fmt"
	"strings"

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
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Deprecated != "" {
			desc += " " + s.Deprecated
		}
		if len(s.ConflictsWith) != 0 {
			desc += " Conflicts with:"
			for index, element := range s.ConflictsWith {
				desc += fmt.Sprintf("`%s`", element)
				if index != len(s.ConflictsWith)-1 {
					desc += ", "
				}
			}
			desc += "."
		}
		if len(s.RequiredWith) != 0 {
			desc += " Required with:"
			for index, element := range s.RequiredWith {
				desc += fmt.Sprintf("`%s`", element)
				if index != len(s.ConflictsWith)-1 {
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
			"taikun_billing_credentials": dataSourceTaikunBillingCredentials(),
			"taikun_billing_credential":  dataSourceTaikunBillingCredential(),
			"taikun_access_profiles":     dataSourceTaikunAccessProfiles(),
			"taikun_access_profile":      dataSourceTaikunAccessProfile(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"taikun_access_profile": resourceTaikunAccessProfile(),
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:          schema.TypeString,
				Description:   "Taikun email.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_EMAIL", nil),
				ConflictsWith: []string{"keycloak_email"},
				RequiredWith:  []string{"password"},
			},
			"password": {
				Type:          schema.TypeString,
				Description:   "Taikun password.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_PASSWORD", nil),
				ConflictsWith: []string{"keycloak_password"},
				RequiredWith:  []string{"email"},
			},
			"keycloak_email": {
				Type:          schema.TypeString,
				Description:   "Taikun keycloak email.",
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_KEYCLOAK_EMAIL", nil),
				ConflictsWith: []string{"email"},
				RequiredWith:  []string{"keycloak_password"},
			},
			"keycloak_password": {
				Type:          schema.TypeString,
				Description:   "Taikun keycloak password.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("TAIKUN_KEYCLOAK_PASSWORD", nil),
				ConflictsWith: []string{"password"},
				RequiredWith:  []string{"keycloak_email"},
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

	return &apiClient{
		client:              client.Default,
		email:               email.(string),
		password:            password.(string),
		useKeycloakEndpoint: keycloakEnabled,
	}, nil
}
