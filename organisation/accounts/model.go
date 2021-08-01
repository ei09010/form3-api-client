package accounts

import "time"

// Account represents an account in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.
// type AccountData struct {
// 	Attributes     *AccountAttributes `json:"attributes,omitempty"`
// 	ID             string             `json:"id,omitempty"`
// 	OrganisationID string             `json:"organisation_id,omitempty"`
// 	Type           string             `json:"type,omitempty"`
// 	Version        *int64             `json:"version,omitempty"`
// }

// type AccountAttributes struct {
// 	AccountClassification   *string  `json:"account_classification,omitempty"`
// 	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
// 	AccountNumber           string   `json:"account_number,omitempty"`
// 	AlternativeNames        []string `json:"alternative_names,omitempty"`
// 	BankID                  string   `json:"bank_id,omitempty"`
// 	BankIDCode              string   `json:"bank_id_code,omitempty"`
// 	BaseCurrency            string   `json:"base_currency,omitempty"`
// 	Bic                     string   `json:"bic,omitempty"`
// 	Country                 *string  `json:"country,omitempty"`
// 	Iban                    string   `json:"iban,omitempty"`
// 	JointAccount            *bool    `json:"joint_account,omitempty"`
// 	Name                    []string `json:"name,omitempty"`
// 	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
// 	Status                  *string  `json:"status,omitempty"`
// 	Switched                *bool    `json:"switched,omitempty"`
// }

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
