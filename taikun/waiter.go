package taikun

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	readAfterCreateOperationTimeout = 2 * time.Minute
	readAfterUpdateOperationTimeout = 1 * time.Minute
)

func timedOut(err error) bool {
	timeoutErr, ok := err.(*resource.TimeoutError)
	return ok && timeoutErr.LastError == nil
}

type readFunc func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics

func readAfterCreateWithRetries(readFunc readFunc, ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resource.RetryContext(ctx, readAfterCreateOperationTimeout, func() *resource.RetryError {
		readErr := readFunc(ctx, data, meta)
		if readErr != nil {
			return resource.NonRetryableError(nil)
		}
		if data.Id() == "" {
			return resource.RetryableError(nil)
		}
		return nil
	})
	if timedOut(err) {
		return diag.Errorf("timed out reading newly created resource")
	}
	return nil
}

func readAfterUpdateWithRetries(readFunc readFunc, ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := resource.RetryContext(ctx, readAfterUpdateOperationTimeout, func() *resource.RetryError {
		readErr := readFunc(ctx, data, meta)
		if readErr != nil {
			return resource.NonRetryableError(nil)
		}
		if data.Id() == "" {
			return resource.RetryableError(nil)
		}
		return nil
	})
	if timedOut(err) {
		return diag.Errorf("timed out reading newly updated resource")
	}
	return nil
}
