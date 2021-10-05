package taikun

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/taikungoclient/client/auth"
	"github.com/itera-io/taikungoclient/client/keycloak"
	"github.com/itera-io/taikungoclient/models"
	"strings"
	"time"
)

type apiClient struct {
	client *client.Taikungoclient

	email               string
	password            string
	useKeycloakEndpoint bool

	token        string
	refreshToken string
}

type jwtData struct {
	Nameid     string `json:"nameid"`
	Email      string `json:"email"`
	UniqueName string `json:"unique_name"`
	Role       string `json:"role"`
	Nbf        int    `json:"nbf"`
	Exp        int    `json:"exp"`
	Iat        int    `json:"iat"`
}

func (apiClient *apiClient) AuthenticateRequest(c runtime.ClientRequest, _ strfmt.Registry) error {

	if len(apiClient.token) == 0 {

		if !apiClient.useKeycloakEndpoint {
			loginResult, err := apiClient.client.Auth.AuthLogin(
				auth.NewAuthLoginParams().WithV(ApiVersion).WithBody(
					&models.LoginCommand{Email: apiClient.email, Password: apiClient.password},
				), nil,
			)
			if err != nil {
				return err
			}
			apiClient.token = loginResult.Payload.Token
			apiClient.refreshToken = loginResult.Payload.RefreshToken
		} else {
			loginResult, err := apiClient.client.Keycloak.KeycloakLogin(
				keycloak.NewKeycloakLoginParams().WithV(ApiVersion).WithBody(
					&models.LoginWithKeycloakCommand{Email: apiClient.email, Password: apiClient.password},
				), nil,
			)
			if err != nil {
				return err
			}
			apiClient.token = loginResult.Payload.Token
			apiClient.refreshToken = loginResult.Payload.RefreshToken
		}

		fmt.Println(apiClient.token)
	}

	if apiClient.hasTokenExpired() {

		refreshResult, err := apiClient.client.Auth.AuthRefreshToken(
			auth.NewAuthRefreshTokenParams().WithV(ApiVersion).WithBody(
				&models.RefreshTokenCommand{
					RefreshToken: apiClient.refreshToken,
					Token:        apiClient.token,
				}), nil,
		)
		if err != nil {
			return err
		}

		apiClient.token = refreshResult.Payload.Token
		apiClient.refreshToken = refreshResult.Payload.RefreshToken
	}

	err := c.SetHeaderParam("Authorization", fmt.Sprintf("Bearer %s", apiClient.token))
	if err != nil {
		return err
	}

	return nil
}

func (apiClient *apiClient) hasTokenExpired() bool {
	jwtSplit := strings.Split(apiClient.token, ".")
	if len(jwtSplit) != 3 {
		return true
	}

	data, err := base64.RawURLEncoding.DecodeString(jwtSplit[1])
	if err != nil {
		return true
	}

	jwtData := jwtData{}
	err = json.Unmarshal(data, &jwtData)
	if err != nil {
		return true
	}

	tm := time.Unix(int64(jwtData.Exp), 0)

	return tm.Before(time.Now())
}
