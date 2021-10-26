package taikun

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func diagnosticsToString(diagnostics diag.Diagnostics) string {
	stringBuilder := strings.Builder{}
	for _, diagnostic := range diagnostics {
		stringBuilder.WriteString(diagnosticToString(diagnostic))
		stringBuilder.WriteString("\n")
	}
	return stringBuilder.String()
}

func diagnosticToString(diagnostic diag.Diagnostic) string {
	stringBuilder := strings.Builder{}
	if diagnostic.Severity == diag.Error {
		stringBuilder.WriteString("[ERROR] ")
	} else {
		stringBuilder.WriteString("[WARN] ")
	}
	stringBuilder.WriteString(diagnostic.Summary)
	stringBuilder.WriteString(": ")
	stringBuilder.WriteString(diagnostic.Detail)
	return stringBuilder.String()
}
