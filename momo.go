package momo

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RequestID is mockable version to generate request ID,
var RequestID = func() string {
	return uuid.New().String()
}

const (
	codeSuccess       int = 0
	codePayoutSuccess int = 1000
)

type BasicRequest struct {
	RequestID string `json:"requestId"`
}

func (r *BasicRequest) GenerateID() string {
	if r.RequestID == "" {
		r.RequestID = RequestID()
	}
	return r.RequestID
}

type BasicResponse struct {
	*Error
	ReferenceID string `json:"referenceId"`
}

func (r BasicResponse) CheckError() error {
	if r.Code == codeSuccess || r.Code == codePayoutSuccess {
		return nil
	}

	return r.Error
}

type BalanceRequest struct {
	BasicRequest
	Password string `json:"password"`
}

type BalanceResponse struct {
	BasicResponse
	Amount int64 `json:"amount"`
}

type CheckInfoRequest struct {
	BasicRequest
	WalletID string `json:"walletId"`
}

type CheckInfoResponse struct {
	BasicResponse
	AccountInfo AccountInfo `json:"accountInfo"`
}

type AccountInfo struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	IsNonBank bool   `json:"isNonBank"`
}

type PaymentPayRequest struct {
	BasicRequest
	Password string `json:"password"`
	*PayInfo
}

type PaymentPayResponse struct {
	BasicResponse
	PayTransaction
}

type PaymentStatusRequest struct {
	BasicRequest
	PaymentID string `json:"paymentId"`
}

type PaymentStatusResponse struct {
	BasicResponse
	Data *PaymentPayResponse `json:"data"`
}

type TransactionsRequest struct {
	BasicRequest
	Password string `json:"password"`
	Date     string `json:"date"`
}

type TransactionsResponse struct {
	BasicResponse
	Data []Transaction `json:"data"`
}

type Request interface {
	GenerateID() string
}

type Response interface {
	CheckError() error
}

type Client interface {
	CheckInfo(ctx context.Context, walletID string) (*AccountInfo, error)
	GetBalance(ctx context.Context) (int64, error)
	GetTransactions(ctx context.Context, date time.Time) error
	PaymentPay(ctx context.Context, info *PayInfo) (*PayTransaction, error)
	PaymentStatus(ctx context.Context, paymentID string) (*PaymentPayResponse, error)
}
