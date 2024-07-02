package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/antithesishq/antithesis-sdk-go/assert"
	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
)

type Details map[string]any

func main() {
	ctx := context.Background()

	client := sdk.New(
		sdk.WithServerURL("http://gateway:8080"),
		sdk.WithClient(&http.Client{
			Timeout: 10 * time.Second,
			//Transport: httpclient.NewDebugHTTPTransport(http.DefaultTransport),
		}),
	)

	waitServicesReady(ctx, client)
	runWorkload(ctx, client)
}

func createLedger(ctx context.Context, client *sdk.Formance) error {
	fmt.Println("Creating ledger...")
	_, err := client.Ledger.V2CreateLedger(ctx, operations.V2CreateLedgerRequest{
		Ledger: "default",
	})

	if assert.Always(err == nil, "ledger should have been created", Details{
		"error": fmt.Sprintf("%+v\n", err),
	}); err != nil {
		return err
	}

	return nil
}
