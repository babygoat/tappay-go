package tappay

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRefund(t *testing.T) {
	for _, tc := range []struct {
		name         string
		setup        func(*testing.T, *client) (string, error)
		requireSetup bool
		wantStatus   int
	}{
		{
			name: "Given empty rec_trade_id returns error",
			setup: func(t *testing.T, c *client) (string, error) {
				return "Invalid_trade_id", nil
			},
			wantStatus: 11000,
		},
		{
			name:       "Given valid rec_trade_id during setup returns success",
			setup:      setupPayment,
			wantStatus: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, _ := NewClient("partner_6ID1DoDlaPrfHw6HBZsULfTYtDmWs0q0ZZGKMBpp4YICWBxgK97eK3RM", WithServer(SandboxAPIURL))
			var recTradeID string
			if tc.setup != nil {
				var err error
				recTradeID, err = tc.setup(t, client)
				if err != nil {
					t.Errorf("unable to setup the target refund payment")
					t.Fail()
				}

			}
			refund, err := client.Refund(context.Background(), RefundParams{RecTradeID: recTradeID})
			if err != nil {
				t.Errorf("unexpected refund error, err: %v", err)
			}

			if refund != nil {
				if tc.wantStatus != refund.Status {
					t.Errorf("expect status: %d, got %d", tc.wantStatus, refund.Status)
				}
			}
			_ = refund
		})
	}
}

func setupPayment(t *testing.T, c *client) (string, error) {
	t.Helper()

	resp, err := c.PayByPrime(context.Background(), PaymentPrimeParams{
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

	return resp.RecTradeID, err
}
