package taikun

import (
	"testing"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
// var providerFactories = map[string]func() (*schema.Provider, error){
// 	"scaffolding": func() (*schema.Provider, error) {
// 		return New("dev")(), nil
// 	},
// }

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
