package tappay

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPayByPrime(t *testing.T) {
	for _, tc := range []struct {
		name       string
		params     PaymentPrimeParams
		wantStatus int
	}{
		{
			name: "Given valid params for PayByPrime returns status 0",
			params: PaymentPrimeParams{
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
			},
			wantStatus: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cli, _ := NewClient("partner_6ID1DoDlaPrfHw6HBZsULfTYtDmWs0q0ZZGKMBpp4YICWBxgK97eK3RM", WithServer(SandboxAPIURL))
			resp, err := cli.PayByPrime(context.Background(), tc.params)
			if err != nil {
				t.Errorf("unexpected pay-by-prime error, err: %v", err)
			}
			if tc.wantStatus != resp.Status {
				t.Errorf("expected pay-by-prime status: %d, got :%d", tc.wantStatus, resp.Status)
			}
		})

	}
}
