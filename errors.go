package momo

import (
	"fmt"
)

type GenericError struct {
	BaseErr   error
	RequestID string `json:"requestId"`
}

func (e GenericError) Error() string {
	return fmt.Sprintf("Momo: %s", e.BaseErr.Error())
}

type Error struct {
	Code          int    `json:"resultCode"`
	LocalMesssage string `json:"localMessage"`
	Message       string `json:"message"`
	RequestID     string `json:"requestId"`
}

func (e Error) Error() string {
	return fmt.Sprintf("Momo: (%d) %s", e.Code, e.Message)
}

var retryableCodes = []int{}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*Error); ok {
		for _, code := range retryableCodes {
			if code == e.Code {
				return true
			}
		}
	}

	return false
}
