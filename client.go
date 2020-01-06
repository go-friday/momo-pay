package momo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	fhttp "github.com/go-friday/http"
	"golang.org/x/crypto/openpgp"
)

type HttpDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type client struct {
	client   HttpDoer
	host     string
	partner  string
	password string
}

func NewClient(
	host string,
	partnerCode string,
	walletPassword string,
	partnerKey, momoKey openpgp.EntityList,
) Client {
	return NewClientWithHTTP(host, partnerCode, walletPassword, &http.Client{
		Timeout: 10 * time.Second,
		Transport: &fhttp.TransportPGP{
			PrivateKey: partnerKey,
			Recipients: momoKey,
		},
	})
}

func NewClientWithHTTP(
	host string,
	partnerCode string,
	walletPassword string,
	httpDoer HttpDoer,
) Client {
	return &client{
		client:   httpDoer,
		host:     host,
		partner:  partnerCode,
		password: walletPassword,
	}
}

func (c *client) SendRequest(ctx context.Context, path string, req Request, res Response) error {
	requestID := req.GenerateID()

	bs, err := json.Marshal(req)
	if err != nil {
		return &GenericError{
			BaseErr:   err,
			RequestID: requestID,
		}
	}

	httpReq, err := http.NewRequest("POST", c.host+path, bytes.NewReader(bs))
	if err != nil {
		return &GenericError{
			BaseErr:   err,
			RequestID: requestID,
		}
	}

	httpReq = httpReq.WithContext(ctx)
	httpReq.Header.Set("Partner-Code", c.partner)
	httpRes, err := c.client.Do(httpReq)
	if err != nil {
		return &GenericError{
			BaseErr:   err,
			RequestID: requestID,
		}
	}
	defer httpRes.Body.Close()

	if err = json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return &GenericError{
			BaseErr:   err,
			RequestID: requestID,
		}
	}

	if err = res.CheckError(); err != nil {
		return err
	}

	return nil
}

func (c *client) GetBalance(ctx context.Context) (int64, error) {
	req := &BalanceRequest{
		Password: c.password,
	}
	res := &BalanceResponse{}
	err := c.SendRequest(ctx, "/api/pay/balance", req, res)
	if err == nil {
		return res.Amount, nil
	}
	return 0, err
}

func (c *client) CheckInfo(ctx context.Context, walletID string) (*AccountInfo, error) {
	req := &CheckInfoRequest{
		WalletID: walletID,
	}
	res := &CheckInfoResponse{}
	err := c.SendRequest(ctx, "/api/pay/check-info", req, res)
	if err == nil {
		return &res.AccountInfo, nil
	}
	return nil, err
}

func (c *client) PaymentPay(ctx context.Context, info *PayInfo) (*PayTransaction, error) {
	req := &PaymentPayRequest{
		Password: c.password,
		PayInfo:  info,
	}
	res := &PaymentPayResponse{}
	err := c.SendRequest(ctx, "/api/payment/pay", req, res)
	if err == nil {
		return &res.PayTransaction, nil
	}
	return nil, err
}

func (c *client) PaymentStatus(ctx context.Context, paymentID string) (Response, error) {
	req := &PaymentStatusRequest{
		PaymentID: paymentID,
	}
	res := &PaymentStatusResponse{}
	err := c.SendRequest(ctx, "/api/payment/status", req, res)
	if err == nil {
		basic := &BasicResponse{}
		if err = json.Unmarshal(res.Data, basic); err != nil {
			return nil, &GenericError{
				BaseErr:   err,
				RequestID: req.RequestID,
			}
		}

		// Detect method base on ReferenceID
		if strings.HasPrefix(basic.ReferenceID, "PAY_") {
			pay := &PaymentPayResponse{}
			if err = json.Unmarshal(res.Data, pay); err != nil {
				return nil, &GenericError{
					BaseErr:   err,
					RequestID: req.RequestID,
				}
			}
			return pay, nil
		}

		return nil, nil
	}

	return nil, &GenericError{
		BaseErr:   err,
		RequestID: req.RequestID,
	}
}

func (c *client) GetTransactions(ctx context.Context, date time.Time) error {
	res := &TransactionsResponse{}
	req := &TransactionsRequest{
		Password: c.password,
		Date:     date.Format("2006/01/02"),
	}
	err := c.SendRequest(ctx, "/api/payment/get-trans", req, res)
	if err == nil {
		return nil
	}
	return err
}
