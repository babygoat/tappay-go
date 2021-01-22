package tappay

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRecord(t *testing.T) {
	for _, tc := range []struct {
		name       string
		setupQuery func(*client) RecordParams
		wantError  bool
	}{
		{
			name: "Given empty filter returns no error",
		},
		{
			name: "Given filter rec_trade_id returns no error",
			setupQuery: func(c *client) RecordParams {
				return RecordParams{Filters: &RecordFilters{RecTradeID: getTestRecTradeID(t, c)}}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cli, _ := NewClient("partner_6ID1DoDlaPrfHw6HBZsULfTYtDmWs0q0ZZGKMBpp4YICWBxgK97eK3RM", WithServer(SandboxAPIURL))
			var query RecordParams
			if tc.setupQuery != nil {
				query = tc.setupQuery(cli)
			}
			records, err := cli.Records(context.Background(), query)
			if err != nil {
				t.Errorf("unexpected get record error: %v", err)
			}
			if tc.wantError && records.Status != 0 && records.Status != 2 {
				t.Errorf("unexpected status: %d", records.Status)
			}
		})
	}
}

func getTestRecTradeID(t *testing.T, c *client) string {
	t.Helper()

	resp, _ := c.PayByPrime(context.Background(), PaymentPrimeParams{
		Prime:       "test_3a2fb2b7e892b914a03c95dd4dd5dc7970c908df67a49527c0a648b2bc9",
		MerchantID:  "GlobalTesting_CTBC",
		Amount:      100,
		OrderNumber: fmt.Sprintf("tappay-go-%d", time.Now().UnixNano()),
		Details:     "test-tappay-go-package",
		Cardholder: PaymentParamsCardholder{
			PhoneNumber: "0912345678",
			Name:        "tappay-go",
			Email:       "tappaygo@example.com",
		},
		Remember: false,
	})

	return resp.RecTradeID
}
