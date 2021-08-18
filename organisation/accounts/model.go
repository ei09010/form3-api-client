package accounts

import "time"

// Response containing all the information related to accounts
type AccountResponse struct {
	*AccountData
	apiErrorMessage
}

type AccountData struct {
	Data  *Data  `json:"data"`
	Links *Links `json:"links"`
}

type Data struct {
	Attributes     *AccountAttributes `json:"attributes"`
	CreatedOn      time.Time          `json:"created_on"`
	ID             string             `json:"id"`
	ModifiedOn     time.Time          `json:"modified_on"`
	OrganisationID string             `json:"organisation_id"`
	Type           string             `json:"type"`
	Version        int                `json:"version"`
}

type AccountAttributes struct {
	AccountClassification string   `json:"account_classification"`
	AlternativeNames      []string `json:"alternative_names"`
	BankID                string   `json:"bank_id"`
	BankIDCode            string   `json:"bank_id_code"`
	BaseCurrency          string   `json:"base_currency"`
	Bic                   string   `json:"bic"`
	Country               string   `json:"country"`
	Name                  []string `json:"name"`
}

type Links struct {
	Self string `json:"self"`
}
