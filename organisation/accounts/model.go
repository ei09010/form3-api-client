package accounts

import (
	"time"
)

// Response containing all the information related to accounts
type AccountResponse struct {
	*AccountData
	apiErrorMessage
}

type AccountData struct {
	Data  *Data  `json:"data" gorm:"type:data"`
	Links *Links `json:"links" gorm:"type:links"`
}

type Account struct {
	Data  *Data  `json:"data" gorm:"type:data"`
	Links *Links `json:"links" gorm:"type:links"`
}

type Data struct {
	Attributes     *AccountAttributes `json:"attributes" gorm:"type:attributes"`
	CreatedOn      time.Time          `json:"created_on" gorm:"type:created_on"`
	ID             string             `json:"id" gorm:"type:id"`
	ModifiedOn     time.Time          `json:"modified_on" gorm:"type:modified_on"`
	OrganisationID string             `json:"organisation_id" gorm:"type:organisation_id"`
	Type           string             `json:"type" gorm:"type:type"`
	Version        int                `json:"version" gorm:"type:version"`
}

type AccountAttributes struct {
	AccountClassification string   `json:"account_classification" gorm:"type:account_classification"`
	AlternativeNames      []string `json:"alternative_names" gorm:"type:alternative_names"`
	BankID                string   `json:"bank_id" gorm:"type:bank_id"`
	BankIDCode            string   `json:"bank_id_code" gorm:"type:bank_id_code"`
	BaseCurrency          string   `json:"base_currency" gorm:"type:base_currency"`
	Bic                   string   `json:"bic" gorm:"type:bic"`
	Country               string   `json:"country" gorm:"type:country"`
	Name                  []string `json:"name" gorm:"type:name"`
}

type Links struct {
	Self string `json:"self" gorm:"type:self"`
}
