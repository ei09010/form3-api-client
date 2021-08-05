package accounts

import (
	"errors"
	"time"
)

const (
	AccountsPath        = "/v1/organisation/accounts"
	DefaultTimeOutValue = time.Duration(10000 * time.Millisecond)
)

// Error Standard Types
var (
	ApiHttpErrorType      = errors.New("API Error")
	UnmarshallingError    = errors.New("UnmarshallingError")
	BaseUrlParsingError   = errors.New("BaseUrlParsingError")
	PathParsingError      = errors.New("PathParsingError")
	FinalUrlParsingError  = errors.New("FinalUrlParsingError")
	BuildingRequestError  = errors.New("BuildingRequestError")
	ExecutingRequestError = errors.New("ExecutingRequestError")
	ResponseReadError     = errors.New("ResponseReadError")
)
