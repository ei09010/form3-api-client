package accounts

import (
	"fmt"
	"strings"
	"time"
)

const (
	AccountsPath        = "/v1/organisation/accounts"
	DefaultTimeOutValue = time.Duration(10000 * time.Millisecond)
)

// Error Standard Types
const (
	HttResponseStandardError = iota
	UnmarshallingError       = iota
	BaseUrlParsingError      = iota
	PathParsingError         = iota
	FinalUrlParsingError     = iota
	BuildingRequestError     = iota
	ExecutingRequestError    = iota
	ResponseError            = iota
)

type BadStatusError struct {
	HttpCode int
	URL      string
}

type ClientError struct {
	ErrorType      int
	ErrorMessage   string
	BadStatusError *BadStatusError
}

func (c *ClientError) Error() string {

	var sb strings.Builder

	baseErrorMessage := fmt.Sprintf("Type: %d | %s", c.ErrorType, c.ErrorMessage)

	sb.WriteString(baseErrorMessage)

	if c.BadStatusError != nil {
		sb.WriteString(c.BadStatusError.Error())
	}

	return fmt.Sprint(sb.String())
}

func (b *BadStatusError) Error() string {
	return fmt.Sprintf("| Did not get 200 from %s, got %d ", b.URL, b.HttpCode)
}

// arguments
type AccountsCreationArg struct {
}
