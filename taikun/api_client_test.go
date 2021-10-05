package taikun

import (
	"testing"

	"github.com/itera-io/taikungoclient/client"
)

func TestExpiredTokenHasExpired(t *testing.T) {

	apiClientWithExpiredToken := apiClient{
		client: client.Default,
		token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjB9.JWKPB-5Q8rTYzl-MfhRGpP9WpDpQxC7JkIAGFMDZnpg",
	}

	if !apiClientWithExpiredToken.hasTokenExpired() {
		t.Fatalf("API client's token has expired but hasTokenExpired returned false")
	}

}

func TestValidTokenHasNotExpired(t *testing.T) {

	apiClientWithExpiredToken := apiClient{
		client: client.Default,
		token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjEwMDAwMDAwMDAwMDB9.2r926yROSOUx-ivZ3RExUWfRtXjoMpSL_F4dxXzfXnY",
	}

	if apiClientWithExpiredToken.hasTokenExpired() {
		t.Fatalf("API client's token is still valid but hasTokenExpired returned true")
	}

}
