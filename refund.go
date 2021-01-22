package tappay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// RefundParams defines the parameters for performing refund operation
// More details in: https://docs.tappaysdk.com/tutorial/zh/back.html#refund-api
type RefundParams struct {
	RecTradeID     string          `json:"rec_trade_id"`
	BankRefundID   string          `json:"bank_refund_id,omitempty"`
	Amount         string          `json:"amount,omitempty"`
	AdditionalData json.RawMessage `json:"additional_data,omitempty"`
}

// MarshalMap implements the Marshaler interface
func (r RefundParams) MarshalMap() (map[string]interface{}, error) {
	p, err := json.Marshal(&r)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal RefundParams: %v, err: %v", r, err)
	}
	var m map[string]interface{}
	if err = json.Unmarshal(p, &m); err != nil {
		return nil, fmt.Errorf("cannot unmarshal RefundParams into map, err: %v", err)
	}

	return m, nil
}

// RefundResponse defines the API response returns by TapPay server after refund request
// More details in: https://docs.tappaysdk.com/tutorial/zh/back.html#response15
type RefundResponse struct {
	Status         int    `json:"status"`
	Msg            string `json:"msg"`
	RefundID       string `json:"refund_id"`
	RefundAmount   int    `json:"refund_amount"`
	IsCaptured     bool   `json:"is_captured"`
	BankResultCode string `json:"bank_result_code"`
	BankResultMsg  string `json:"bank_result_msg"`
	Currency       string `json:"currency"`
}

// refundPath defines the refund service path
const refundPath = "/tpc/transaction/refund"

// Refund issues a refund request according to input RefundParams
// and returns the parsed RefundResponse from TapPay server
func (c *client) Refund(ctx context.Context, params RefundParams) (*RefundResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, serviceRefund, params)
	if err != nil {
		return nil, err
	}

	rawResp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	var resp RefundResponse
	if err = json.Unmarshal(rawResp, &resp); err != nil {
		return nil, fmt.Errorf("cannot unmarshal RefundResponse, err: %v", err)
	}

	return &resp, nil
}
