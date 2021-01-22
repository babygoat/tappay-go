package tappay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// payByPrimePath defines the path of pay-by-prime service
const payByPrimePath = "/tpc/payment/pay-by-prime"

// PaymentParamsCardholder defines the field `cardholder` in request to pay-by-prime api
// See PaymentPrimeParams for more details
type PaymentParamsCardholder struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	ZipCode     string `json:"zip_code,omitempty"`
	Address     string `json:"address,omitempty"`
	NationalID  string `json:"national_id,omitempty"`
	MemberID    string `json:"member_id,omitempty"`
}

// PaymentParamsResultUrl defines the field `result_url` in request to pay-by-prime api
// See PaymentPrimeParams for more details
type PaymentParamsResultUrl struct {
	FrontendRedirectUrl string `json:"frontend_redirect_url"`
	BackendNotifyUrl    string `json:"backend_notify_url"`
}

// PaymentParamsCardholderVerify defines the field `cardholder_verify` in request to pay-by-prime api
// See PaymentPrimeParams for more details
type PaymentParamsCardholderVerify struct {
	PhoneNumber bool `json:"phone_number,omitempty"`
	NationalID  bool `json:"national_id,omitempty"`
}

// PaymentPrimeParams defines the parameters for performing pay-by-prime operation
// More details in: https://docs.tappaysdk.com/tutorial/zh/back.html#pay-by-prime-api
type PaymentPrimeParams struct {
	Prime              string                         `json:"prime"`
	MerchantID         string                         `json:"merchant_id"`
	MerchantGroupID    string                         `json:"merchant_group_id,omitempty"`
	Amount             int                            `json:"amount"`
	MerchandiseDetails *RecordMerchandiseDetails      `json:"merchandise_details,omitempty"`
	Currency           string                         `json:"currency,omitempty"`
	OrderNumber        string                         `json:"order_number,omitempty"`
	BankTransactionID  string                         `json:"bank_transaction_id,omitempty"`
	Details            string                         `json:"details"`
	Cardholder         PaymentParamsCardholder        `json:"cardholder"`
	CardholderVerify   *PaymentParamsCardholderVerify `json:"cardholder_verify,omitempty"`
	Instalment         int                            `json:"instalment,omitempty"`
	DelayCaptureInDays int                            `json:"delay_capture_in_days,omitempty"`
	ThreeDomainSecure  bool                           `json:"three_domain_secure,omitempty"`
	ResultUrl          *PaymentParamsResultUrl        `json:"result_url,omitempty"`
	Remember           bool                           `json:"remember,omitempty"`
	Redeem             bool                           `json:"redeem,omitempty"`
	AdditionalData     json.RawMessage                `json:"additional_data,omitempty"`
	EventCode          string                         `json:"event_code,omitempty"`
	ProductImageUrl    string                         `json:"product_image_url,omitempty"`
}

// MarshalMap implements the Marshaler interface
func (r PaymentPrimeParams) MarshalMap() (map[string]interface{}, error) {
	p, err := json.Marshal(&r)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal PaymentPrimeParams: %v, err: %v", r, err)
	}
	var m map[string]interface{}
	if err = json.Unmarshal(p, &m); err != nil {
		return nil, fmt.Errorf("cannot unmarshal PaymentPrimeParams into map, err: %v", err)
	}

	return m, nil
}

// PaymentCardSecret defines the field `card_secret` in PaymentPrimeResponse
// See PaymentPrimeResponse for more details
type PaymentCardSecret struct {
	CardToken string `json:"card_token"`
	CardKey   string `json:"card_key"`
}

// PaymentCardInfo defines the field `card_info` in PaymentPrimeResponse
// See PaymentPrimeResponse for more details
type PaymentCardInfo struct {
	RecordCardInfo
	ExpiryDate string `json:"expiry_date"`
}

// PaymentBankTransactionTime defines the field `bank_transaction_time` in PaymentPrimeResponse
// See PaymentPrimeResponse for more details
type PaymentBankTransactionTime struct {
	StartTimeMillis string `json:"start_time_millis"`
	EndTimeMillis   string `json:"end_time_millis"`
}

// PaymentRedeemExtraInfo defines the field `redeem_extra_info` in PaymentPrimeResponse
// See PaymentPrimeResponse for more details
type PaymentRedeemExtraInfo struct {
	RedeemUsed    string `json:"redeem_used"`
	CreditAmt     string `json:"credit_amt"`
	RedeemBalance string `json:"redeem_balance"`
	RedeemType    string `json:"redeem_type"`
}

// PaymentRedeemInfo defines the field `redeem_info` in PaymentPrimeResponse
// See PaymentPrimeResponse for more details
type PaymentRedeemInfo struct {
	RecordRedeemInfo
	ExtraInfo PaymentRedeemExtraInfo `json:"extra_info"`
}

// PaymentPrimeResponse defines the API response returns by TapPay server after pay-by-prime request
// More details in: https://docs.tappaysdk.com/tutorial/zh/back.html#response
type PaymentPrimeResponse struct {
	Status                int                         `json:"status"`
	Msg                   string                      `json:"msg"`
	RecTradeID            string                      `json:"rec_trade_id"`
	BankTransactionID     string                      `json:"bank_transaction_id"`
	AuthCode              string                      `json:"auth_code"`
	CardSecret            PaymentCardSecret           `json:"card_secret"`
	Amount                int                         `json:"amount"`
	Currency              string                      `json:"currency"`
	CardInfo              PaymentCardInfo             `json:"card_info"`
	OrderNumber           string                      `json:"order_number"`
	Acquirer              string                      `json:"acquirer"`
	TransactionTimeMillis int64                       `json:"transaction_time_millis"`
	BankTransactionTime   PaymentBankTransactionTime  `json:"bank_transaction_time"`
	BankResultCode        string                      `json:"bank_result_code"`
	BankResultMsg         string                      `json:"bank_result_msg"`
	PaymentUrl            string                      `json:"payment_url"`
	InstalmentInfo        RecordInstalmentInfo        `json:"instalment_info"`
	RedeemInfo            PaymentRedeemInfo           `json:"redeem_info"`
	CardIdentifier        string                      `json:"card_identifier"`
	MerchantReferenceInfo RecordMerchantReferenceInfo `json:"merchant_reference_info"`
	EventCode             string                      `json:"event_code"`
}

// PayByPrime issues a pay-by-prime request according to input PaymentPrimeParams
// and parses the response from TapPay server as PaymentPrimeResponse
func (c *client) PayByPrime(ctx context.Context, params PaymentPrimeParams) (*PaymentPrimeResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, servicePayByPrime, params)
	if err != nil {
		return nil, err
	}

	rawResp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	var resp PaymentPrimeResponse
	if err = json.Unmarshal(rawResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
