package utils

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetReadAfterOpTimeout(isUpdate bool) time.Duration {
	if isUpdate {
		return 1 * time.Minute
	}
	return 2 * time.Minute
}

func TimedOut(err error) bool {
	//timeoutErr, ok := err.(*resource.TimeoutError)
	timeoutErr, ok := err.(*retry.TimeoutError)
	return ok && timeoutErr.LastError == nil
}

func ReadAfterCreateWithRetries(readFunc schema.ReadContextFunc, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	isUpdate := false
	return readAfterOpWithRetries(readFunc, ctx, d, meta, isUpdate)
}

func ReadAfterUpdateWithRetries(readFunc schema.ReadContextFunc, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	isUpdate := true
	return readAfterOpWithRetries(readFunc, ctx, d, meta, isUpdate)
}

func readAfterOpWithRetries(readFunc schema.ReadContextFunc, ctx context.Context, d *schema.ResourceData, meta interface{}, isUpdate bool) diag.Diagnostics {
	retryErr := retry.RetryContext(ctx, GetReadAfterOpTimeout(isUpdate), func() *retry.RetryError {
		readDiagnostics := readFunc(ctx, d, meta)
		if readDiagnostics != nil {

			if readDiagnostics[0].Summary == NotFoundAfterCreateOrUpdateError {
				return retry.RetryableError(errors.New("failed to read after create/update"))
			}

			readErrors := diagnosticsToString(readDiagnostics)
			return retry.NonRetryableError(errors.New(readErrors))
		}
		return nil
	})
	if TimedOut(retryErr) {
		if isUpdate {
			return diag.Errorf("timed out reading newly updated resource")
		}
		return diag.Errorf("timed out reading newly created resource")
	}
	return diag.FromErr(retryErr)
}
