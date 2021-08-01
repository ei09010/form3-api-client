package accounts

import (
	"fmt"
	"strings"
)

// errrors

// arguments
type AccountsCreationArg struct {
}

const (
	HttResponseStandardError = iota
	UnmarshallingError       = iota
	UrlParsingError          = iota
	RequestError             = iota
	ResponseError            = iota
)

type ApiHttpError struct {
	ErrorMessage string `json:"error_message"`
}

type ClientError struct {
	ErrorType    int
	HttpCode     int
	ErrorMessage string
}

func (hError *ClientError) Error() string {

	var sb strings.Builder

	baseErrorMessage := fmt.Sprintf("Type: %d | Message: %s", hError.ErrorType, hError.ErrorMessage)

	sb.WriteString(baseErrorMessage)

	if hError.HttpCode != 0 {
		sb.WriteString(fmt.Sprintf("| Code: %d", hError.HttpCode))
	}

	return fmt.Sprint(sb.String())

}
