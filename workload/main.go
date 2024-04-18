package main

import (
	"context"
	"fmt"
	"github.com/antithesishq/antithesis-sdk-go/assert"
	"github.com/antithesishq/antithesis-sdk-go/lifecycle"
	"github.com/antithesishq/antithesis-sdk-go/random"
	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/httpclient"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"golang.org/x/sync/errgroup"
	"math"
	"math/big"
	"net/http"
	"os"
	"time"
)

type Details map[string]any

func main() {
	ctx := context.Background()
	client := sdk.New(
		sdk.WithServerURL("http://gateway:8080"),
		sdk.WithClient(&http.Client{
			Timeout:   5 * time.Second,
			Transport: httpclient.NewDebugHTTPTransport(http.DefaultTransport),
		}),
	)

	waitServicesReady(ctx, client)

	// signals that the system is up and running
	lifecycle.SetupComplete(Details{"Ledger": "Available"})

	runWorkload(ctx, client)
}

func waitServicesReady(ctx context.Context, client *sdk.Formance) {
	fmt.Println("Waiting for services to be ready")
	waitingServicesCtx, cancel := context.WithDeadline(ctx, time.Now().Add(30*time.Second))
	defer cancel()

	for {
		select {
		case <-time.After(time.Second):
			fmt.Println("Trying to get info of the ledger...")
			_, err := client.Ledger.GetInfo(ctx)
			if err != nil {
				fmt.Printf("error pinging ledger: %s\r\n", err)
				continue
			}
			return
		case <-waitingServicesCtx.Done():
			fmt.Printf("timeout waiting for services to be ready\r\n")
			os.Exit(1)
		}
	}
}

func runWorkload(ctx context.Context, client *sdk.Formance) {
	const count = 100

	fmt.Println("Creating ledger...")
	_, err := client.Ledger.V2CreateLedger(ctx, operations.V2CreateLedgerRequest{
		Ledger: "default",
	})
	assert.Always(err == nil, "ledger should have been created", Details{
		"error": fmt.Sprintf("%+v\n", err),
	})

	totalAmount := big.NewInt(0)

	fmt.Printf("Insert %d transactions...\r\n", count)
	grp, _ := errgroup.WithContext(ctx)
	for i := 0; i < count; i++ {
		amount := big.NewInt(int64(math.Abs(float64(random.GetRandom()))))
		totalAmount = totalAmount.Add(totalAmount, amount)
		grp.Go(runTrade(ctx, client, amount))
	}

	err = grp.Wait()
	assert.Always(err == nil, "all transactions should have been written", Details{
		"error": fmt.Sprintf("%+v\n", err),
	})

	fmt.Println("Checking balance of 'world'...")
	account, err := client.Ledger.V2GetAccount(ctx, operations.V2GetAccountRequest{
		Address: "world",
		Expand:  pointer.For("volumes"),
		Ledger:  "default",
	})
	assert.Always(err == nil, "we should be able to query account 'world'", Details{
		"error": fmt.Sprintf("%+v\n", err),
	})
	if err == nil {
		output := account.V2AccountResponse.Data.Volumes["USD/2"].Output
		assert.Always(
			output.Cmp(totalAmount) == 0,
			// TODO: confirm this logic is correct
			"output of 'world' should be 0",
			Details{
				"output": output,
			},
		)
	}
}

func runTrade(ctx context.Context, client *sdk.Formance, amount *big.Int) func() error {
	return func() error {
		orderID := fmt.Sprint(int64(math.Abs(float64(random.GetRandom()))))

		_, err := client.Ledger.V2CreateTransaction(ctx, operations.V2CreateTransactionRequest{
			V2PostTransaction: shared.V2PostTransaction{
				Postings: []shared.V2Posting{{
					Amount:      amount,
					Asset:       "USD/2",
					Destination: fmt.Sprintf("orders:%s", orderID),
					Source:      "world",
				}},
			},
			Ledger: "default",
		})
		assert.Always(err == nil, "creating transaction from @world to @bank should always return a nil error", Details{})

		return nil
	}
}
