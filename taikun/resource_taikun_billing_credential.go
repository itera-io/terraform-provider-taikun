package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/ops_credentials"
	"strconv"
)

func resourceTaikunBillingCredentialRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := apiClient.client.OpsCredentials.OpsCredentialsList(ops_credentials.NewOpsCredentialsListParams().WithV(ApiVersion).WithID(&id), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.Payload.TotalCount == 1 {
		rawBillingCredential := response.GetPayload().Data[0]

		if err := data.Set("created_by", rawBillingCredential.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", strconv.Itoa(int(rawBillingCredential.ID))); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_locked", rawBillingCredential.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawBillingCredential.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawBillingCredential.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawBillingCredential.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", strconv.Itoa(int(rawBillingCredential.OrganizationID))); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawBillingCredential.OrganizationName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("prometheus_password", rawBillingCredential.PrometheusPassword); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("prometheus_url", rawBillingCredential.PrometheusURL); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("prometheus_username", rawBillingCredential.PrometheusUsername); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}
