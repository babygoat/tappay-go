package tappay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RecordStatus int

const (
	RecordStatusError RecordStatus = iota + -1
	RecordStatusAuth
	RecordStatusOK
	RecordStatusPartialRefunded
	RecordStatusRefunded
	RecordStatusPending
	RecordStatusCancel
)

// recordPath defines the path of record query service
const recordPath = "/tpc/transaction/query"

// RecordFilterTime defines the filter query with a range of transaction times
type RecordFilterTime struct {
	StartTime int64 `json:"start_time,omitempty"`
	EndTime   int64 `json:"end_time,omitempty"`
}

// RecordFilterAmount defines the filter to query within a range of transaction amounts
type RecordFilterAmount struct {
	UpperLimit int `json:"upper_limit,omitempty"`
	LowerLimit int `json:"lower_limit,omitempty"`
}

// RecordFilterCardholder defines the filter to query w.r.t the cardholder properties
type RecordFilterCardholder struct {
	PhoneNumber string `json:"phone_number,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
}

// RecordFilters defines a collection of filters that can be used within the records query operation
type RecordFilters struct {
	Time              *RecordFilterTime       `json:"time,omitempty"`
	Amount            *RecordFilterAmount     `json:"amount,omitempty"`
	Cardholder        *RecordFilterCardholder `json:"cardholder,omitempty"`
	MerchantID        []string                `json:"merchant_id,omitempty"`
	RecordStatus      int                     `json:"record_status,omitempty"`
	RecTradeID        string                  `json:"rec_trade_id,omitempty"`
	OrderNumber       string                  `json:"order_number,omitempty"`
	BankTransactionID string                  `json:"bank_transaction_id,omitempty"`
	Currency          string                  `json:"currency,omitempty"`
}

// RecordSort defines the field(time or amount) and type(ascending/descending) of the returned records
type RecordSort struct {
	Attribute    string `json:"attribute,omitempty"`
	IsDescending bool   `json:"is_descending,omitempty"`
}

// RecordParams defines the params for performing records query operation
// More details in: https://docs.tappaysdk.com/tutorial/zh/back.html#record-api
type RecordParams struct {
	RecordsPerPage int            `json:"records_per_page,omitempty"`
	Page           int            `json:"page,omitempty"`
	Filters        *RecordFilters `json:"filters,omitempty"`
	OrderBy        *RecordSort    `json:"order_by,omitempty"`
}

func (r RecordParams) MarshalMap() (map[string]interface{}, error) {
	p, err := json.Marshal(&r)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal RecordParams: %v, err: %v", r, err)
	}
	var m map[string]interface{}
	if err = json.Unmarshal(p, &m); err != nil {
		return nil, fmt.Errorf("cannot unmarshal RecordParams into map, err: %v", err)
	}

	return m, nil
}

// RecordCardholder defines the `cardholder` field in Record.
// See Record for more details.
type RecordCardholder struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

// RecordMerchandiseDetails defines the `merchandise_details` field in Record.
// See Record for more details.
type RecordMerchandiseDetails struct {
	NoRebateAmount int `json:"no_rebate_amount"`
}

// RecordMerchantReferenceInfo defines the `merchant_reference_info` field in Record.
// See Record for more details.
type RecordMerchantReferenceInfo struct {
	AffiliateCodes []string `json:"affiliate_codes"`
}

// RecordEInvoiceCarrier defines the `e_invoice_info` field in Record.
// See Record for more details.
type RecordEInvoiceCarrier struct {
	Type       int    `json:"type"`
	Number     string `json:"number"`
	Donation   bool   `json:"donation"`
	DonationID string `json:"donation_id"`
}

// RecordInstalmentInfo defines the `instalment_info` field in Record.
// See Record for more details.
type RecordInstalmentInfo struct {
	NumberOfInstalments int `json:"number_of_instalments"`
	FirstPayment        int `json:"first_payment"`
	EachPayment         int `json:"each_payment"`
}

// RecordPayInfo defines the `pay_info` field in Record.
// See Record for more details.
type RecordPayInfo struct {
	Method                 string `json:"method"`
	MaskedCreditCardNumber string `json:"masked_credit_card_number"`
	Point                  int    `json:"point"`
	Discount               int    `json:"discount"`
	CreditCard             int    `json:"credit_card"`
	Balance                int    `json:"balance"`
	BankAccount            int    `json:"bank_account"`
}

// RecordRedeemInfo defines the `redeem_info` field in Record.
// See Record for more details.
type RecordRedeemInfo struct {
	UsedPoint    string `json:"used_point"`
	Balance      string `json:"balance"`
	OffsetAmount string `json:"offset_amount"`
	DueAmount    string `json:"due_amount"`
}

// RecordCardInfo defines the `card_info` field in Record.
// See Record for more details.
type RecordCardInfo struct {
	BinCode     string `json:"bin_code"`
	LastFour    string `json:"last_four"`
	Issuer      string `json:"issuer"`
	IssuerZhTw  string `json:"issuer_zh_tw"`
	BankID      string `json:"bank_id"`
	Funding     int    `json:"funding"`
	Type        int    `json:"type"`
	Level       string `json:"level"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}

// Record defines trade record returned from TapPay server after record query
// More details in: https://docs.tappaysdk.com/tutorial/zh/reference.html#trade_records
type Record struct {
	RecTradeID                 string                      `json:"rec_trade_id"`
	AuthCode                   string                      `json:"auth_code"`
	MerchantID                 string                      `json:"merchant_id"`
	MerchantName               string                      `json:"merchant_name"`
	AppName                    string                      `json:"app_name"`
	Time                       int64                       `json:"time"`
	Amount                     int                         `json:"amount"`
	RefundedAmount             int                         `json:"refunded_amount"`
	RecordStatus               RecordStatus                `json:"record_status"`
	BankTransactionID          string                      `json:"bank_transaction_id"`
	CapMillis                  int64                       `json:"cap_millis"`
	OriginalAmount             int                         `json:"original_amount"`
	BankTransactionStartMillis int64                       `json:"bank_transaction_start_millis"`
	BankTransactionEndMillis   int64                       `json:"bank_transaction_end_millis"`
	IsCaptured                 bool                        `json:"is_captured"`
	BankResultCode             string                      `json:"bank_result_code"`
	BankResultMsg              string                      `json:"bank_result_msg"`
	PartialCardNumber          string                      `json:"partial_card_number"`
	PaymentMethod              string                      `json:"payment_method"`
	Details                    string                      `json:"details"`
	Cardholder                 RecordCardholder            `json:"cardholder"`
	MerchandiseDetails         RecordMerchandiseDetails    `json:"merchandise_details"`
	Currency                   string                      `json:"currency"`
	MerchantReferenceInfo      RecordMerchantReferenceInfo `json:"merchant_reference_info"`
	EInvoiceCarrier            RecordEInvoiceCarrier       `json:"e_invoice_carrier"`
	ThreeDomainSecure          bool                        `json:"three_domain_secure"`
	PayByInstalment            bool                        `json:"pay_by_instalment"`
	InstalmentInfo             RecordInstalmentInfo        `json:"instalment_info"`
	OrderNumber                string                      `json:"order_number"`
	PayInfo                    RecordPayInfo               `json:"pay_info"`
	PayByRedeem                bool                        `json:"pay_by_redeem"`
	RedeemInfo                 RecordRedeemInfo            `json:"redeem_info"`
	CardIdentifier             string                      `json:"card_identifier"`
	CardInfo                   RecordCardInfo              `json:"card_info"`
}

// RecordResponse defines the API response returns from TapPay server after records query is issued
// More details in: https://docs.tappaysdk.com/tutorial/zh/back.html#response17
type RecordResponse struct {
	Status               int      `json:"status"`
	Msg                  string   `json:"msg"`
	RecordsPerPage       int      `json:"records_per_page"`
	Page                 int      `json:"page"`
	TotalPageCount       int      `json:"total_page_count"`
	NumberOfTransactions int64    `json:"number_of_transactions"`
	TradeRecords         []Record `json:"trade_records"`
}

// Records issues a record query request according to the input RecordParams
// and parse the response from TapPay server as RecordResponse
func (c *client) Records(ctx context.Context, params RecordParams) (*RecordResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, serviceRecord, params)
	if err != nil {
		return nil, err
	}

	rawResp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	var resp RecordResponse
	if err = json.Unmarshal(rawResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
