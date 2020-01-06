package momo

import (
	"time"
)

type AccountInfo struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	IsNonBank bool   `json:"isNonBank"`
}

type PayInfo struct {
	Amount      int64     `json:"amount"`
	Created     time.Time `json:"created,omitempty"`
	Description string    `json:"description,omitempty"`
	PersonalID  string    `json:"personalId,omitempty"`
	WalletID    string    `json:"walletId,omitempty"`

	ValidatePersonalID  bool   `json:"validatePersionalId"`
	RejectIfOverBalance bool   `json:"rejectIfOverBalance"`
	WalletName          string `json:"walletName,omitempty"`

	RequiredOtp  bool   `json:"requiredOtp"`
	ChecksumKey  string `json:"checksumKey,omitempty"`
	NotifyURL    string `json:"notifyUrl,omitempty"`
	VerifyOtpURL string `json:"verifyOtpUrl,omitempty"`
}

type PayTransaction struct {
	AcceptedAmount int64  `json:"acceptedAmount"`
	PaymentRef     string `json:"paymentRef"`
}

type Transaction struct {
	*Error
	WalletID    string    `json:"walletId"`
	Amount      int64     `json:"amount"`
	PaymentDate time.Time `json:"paymentDate"`
	PaymentRef  string    `json:"paymentRef"`
}
