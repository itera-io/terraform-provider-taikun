package taikun

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getReadAfterOpTimeout(isUpdate bool) time.Duration {
	if isUpdate {
		return 1 * time.Minute
	}
	return 2 * time.Minute
}

func timedOut(err error) bool {
	timeoutErr, ok := err.(*resource.TimeoutError)
	return ok && timeoutErr.LastError == nil
}

func readAfterCreateWithRetries(readFunc schema.ReadContextFunc, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	isUpdate := false
	return readAfterOpWithRetries(readFunc, ctx, d, meta, isUpdate)
}

func readAfterUpdateWithRetries(readFunc schema.ReadContextFunc, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	isUpdate := true
	return readAfterOpWithRetries(readFunc, ctx, d, meta, isUpdate)
}

func readAfterOpWithRetries(readFunc schema.ReadContextFunc, ctx context.Context, d *schema.ResourceData, meta interface{}, isUpdate bool) diag.Diagnostics {
	retryErr := resource.RetryContext(ctx, getReadAfterOpTimeout(isUpdate), func() *resource.RetryError {
		readDiagnostics := readFunc(ctx, d, meta)
		if readDiagnostics != nil {

			if readDiagnostics[0].Summary == notFoundAfterCreateOrUpdateError {
				return resource.RetryableError(errors.New("failed to read after create/update"))
			}

			readErrors := diagnosticsToString(readDiagnostics)
			return resource.NonRetryableError(errors.New(readErrors))
		}
		return nil
	})
	if timedOut(retryErr) {
		if isUpdate {
			return diag.Errorf("timed out reading newly updated resource")
		}
		return diag.Errorf("timed out reading newly created resource")
	}
	return diag.FromErr(retryErr)
}
