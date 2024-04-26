package main

import (
	"context"
	"fmt"
	"github.com/alitto/pond"
	"github.com/antithesishq/antithesis-sdk-go/assert"
	"github.com/antithesishq/antithesis-sdk-go/lifecycle"
	"github.com/antithesishq/antithesis-sdk-go/random"
	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"go.uber.org/atomic"
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
			Timeout: 10 * time.Second,
			//Transport: httpclient.NewDebugHTTPTransport(http.DefaultTransport),
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

func randomBigInt() *big.Int {
	v := random.GetRandom()
	ret := big.NewInt(0)
	ret.SetString(fmt.Sprintf("%d", v), 10)
	return ret
}

func runWorkload(ctx context.Context, client *sdk.Formance) {
	const count = 100

	fmt.Println("Creating ledger...")
	_, err := client.Ledger.V2CreateLedger(ctx, operations.V2CreateLedgerRequest{
		Ledger: "default",
	})
	if !assert.Always(err == nil, "ledger should have been created", Details{
		"error": fmt.Sprintf("%+v\n", err),
	}) {
		return
	}

	pool := pond.New(20, 10000)

	totalAmount := big.NewInt(0)

	hasError := atomic.NewBool(false)

	fmt.Printf("Insert %d transactions...\r\n", count)
	for i := 0; i < count; i++ {
		amount := randomBigInt()
		totalAmount = totalAmount.Add(totalAmount, amount)
		pool.Submit(func() {
			if err := runTrade(ctx, client, amount); err != nil {
				hasError.CompareAndSwap(false, true)
			}
		})
	}

	pool.StopAndWait()

	if !assert.Always(!hasError.Load(), "all transactions should have been written", Details{
		"error": fmt.Sprintf("%+v\n", err),
	}) {
		return
	}

	fmt.Println("Checking balance of 'world'...")
	account, err := client.Ledger.V2GetAccount(ctx, operations.V2GetAccountRequest{
		Address: "world",
		Expand:  pointer.For("volumes"),
		Ledger:  "default",
	})
	if !assert.Always(err == nil, "we should be able to query account 'world'", Details{
		"error": fmt.Sprintf("%+v\n", err),
	}) {
		return
	}

	output := account.V2AccountResponse.Data.Volumes["USD/2"].Output
	if !assert.Always(output != nil, "Expect output of world for USD/2 to be not empty", Details{}) {
		return
	}
	fmt.Printf("Expect output of world to be %s and got %d\r\n", totalAmount, output)
	assert.Always(
		output.Cmp(totalAmount) == 0,
		"output of 'world' should match",
		Details{
			"output": output,
		},
	)
}

func runTrade(ctx context.Context, client *sdk.Formance, amount *big.Int) error {
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
	assert.Always(err == nil, "creating transaction from @world to @bank should always return a nil error", Details{
		"error": fmt.Sprintf("%+v\n", err),
	})

	return err
}