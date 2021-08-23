package integration_test

import (
	"database/sql/driver"
	"ei09010/form3-api-client/organisation/accounts"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// prepare

type e2eTestSuite struct {
	suite.Suite
	dbConnectionStr string
	dbConn          *gorm.DB
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		"localhost", "5432", "root", "interview_accountapi", "password") //Build connection string
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	s.dbConn = conn
}

func (s *e2eTestSuite) TestFetch_FetchesAccount_ReturnsAccount() {

	// Arrange
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL("http://localhost:8080"))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	storedTestAccount := generateAccountDataToStore(id)

	s.NoError(s.dbConn.Create(storedTestAccount).Error)

	exectedAccountData := generatedExpectedAccountToBeReturnedByAPI(id)

	// Act

	fetchedAccountData, err := accountsClient.Fetch(id)

	s.Require().NoError(err)

	// Assert

	assert.Equal(s.T(), exectedAccountData.Data.ID, fetchedAccountData.Data.ID, "ID from the fetched account, should match the expected to be returned by the API")

	assert.NotNil(s.T(), fetchedAccountData.Data.ModifiedOn, "ModifiedOn from the fetched account, should not be nil")

	assert.NotNil(s.T(), fetchedAccountData.Data.CreatedOn, "CreatedOn from the fetched account, should not be nil")

	assert.Equal(s.T(), exectedAccountData.Data.OrganisationID, fetchedAccountData.Data.OrganisationID, "OrganisationId from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Type, fetchedAccountData.Data.Type, "Type from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Version, fetchedAccountData.Data.Version, "Version from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Links.Self, fetchedAccountData.Links.Self, "Self links from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Attributes.AccountClassification, fetchedAccountData.Data.Attributes.AccountClassification, "AccountClassification from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Attributes.AlternativeNames, fetchedAccountData.Data.Attributes.AlternativeNames, "AlternativeNames from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Attributes.BankID, fetchedAccountData.Data.Attributes.BankID, "BankId from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Attributes.BankIDCode, fetchedAccountData.Data.Attributes.BankIDCode, "BankIdCode from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Attributes.BaseCurrency, fetchedAccountData.Data.Attributes.BaseCurrency, "BaseCurrency from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Attributes.Bic, fetchedAccountData.Data.Attributes.Bic, "Bic from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Attributes.Country, fetchedAccountData.Data.Attributes.Country, "Country from the fetched account, should match the expected to be returned by the API")

	assert.Equal(s.T(), exectedAccountData.Data.Attributes.Name, fetchedAccountData.Data.Attributes.Name, "Name from the fetched account, should match the expected to be returned by the API")
}

func (s *e2eTestSuite) TestCreate_CreatesAccount_ReturnsAccountCreated() {

	// Arrange

	// this url has to be a env variable
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL("http://localhost:8080"))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	accountDataToStore := generatedExpectedAccountToBeReturnedByAPI(id)

	// Act

	storedAccountData, err := accountsClient.Create(accountDataToStore)

	s.Require().NoError(err)

	// Assert

	// assert stored account data
	//v this is a placeholder
	fmt.Printf("%v", storedAccountData)
}

func (*Account) TableName() string {
	return "Account"
}

// JSONB Interface for JSONB Field of Account Table
// type JSONB []interface{}

// Value Marshal
func (a AccountAttributes) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *AccountAttributes) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func generateAccountDataToStore(id uuid.UUID) *Account {

	expectedAccountClassification := "Personal"
	expectedAlternativeNames := []string{"Alternative Names."}
	expectedBankId := "400300"
	expectedBankIdCode := "GBDSC"
	expectedBaseCurrency := "GBP"
	expectedBic := "NWBKGB22"
	expectedCountry := "GB"
	expectedName := []string{"Name of the account holder, up to four lines possible."}

	expectedOrganisationId := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	expectedVersion := 0

	return &Account{
		ID:             id,
		ModifiedOn:     time.Now(),
		CreatedOn:      time.Now(),
		IsDeleted:      false,
		IsLocked:       false,
		OrganisationID: expectedOrganisationId,
		Version:        expectedVersion,
		Record: &AccountAttributes{
			AccountClassification: expectedAccountClassification,
			AlternativeNames:      expectedAlternativeNames,
			BankID:                expectedBankId,
			BankIDCode:            expectedBankIdCode,
			BaseCurrency:          expectedBaseCurrency,
			Bic:                   expectedBic,
			Country:               expectedCountry,
			Name:                  expectedName,
		},
	}
}

type Account struct {
	ID             uuid.UUID          `gorm:"unique"`
	ModifiedOn     time.Time          `json:"modified_on" gorm:"type:modified_on"`
	CreatedOn      time.Time          `json:"created_on" gorm:"type:created_on"`
	OrganisationID string             `json:"organisation_id" gorm:"type:organisation_id"`
	Version        int                `json:"version" gorm:"type:version"`
	IsDeleted      bool               `gorm:"type:is_deleted"`
	IsLocked       bool               `gorm:"type:is_locked"`
	Record         *AccountAttributes `gorm:"type:jsonb" json:"record"`
}

func generatedExpectedAccountToBeReturnedByAPI(id uuid.UUID) *accounts.AccountData {

	expectedAccountClassification := "Personal"
	expectedAlternativeNames := []string{"Alternative Names."}
	expectedBankId := "400300"
	expectedBankIdCode := "GBDSC"
	expectedBaseCurrency := "GBP"
	expectedBic := "NWBKGB22"
	expectedCountry := "GB"
	expectedName := []string{"Name of the account holder, up to four lines possible."}

	timeLayout := "2006-01-02 15:04:05 -0700 MST"
	expectedCreatedOn := "2021-07-31 22:09:02 +0000 UTC"
	expectedCreatedOnTime, _ := time.Parse(timeLayout, expectedCreatedOn)

	expectedId := id

	expectedModifiedOn := "2021-07-31 22:09:02 +0000 UTC"
	expectedModifiedOnTime, _ := time.Parse(timeLayout, expectedModifiedOn)

	expectedOrganisationId := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	expectedType := "accounts"
	expectedVersion := 0
	expectedSelf := fmt.Sprintf("/v1/organisation/accounts/%s", id.String())

	return &accounts.AccountData{
		Data: &accounts.Data{
			Attributes: &accounts.AccountAttributes{
				AccountClassification: expectedAccountClassification,
				AlternativeNames:      expectedAlternativeNames,
				BankID:                expectedBankId,
				BankIDCode:            expectedBankIdCode,
				BaseCurrency:          expectedBaseCurrency,
				Bic:                   expectedBic,
				Country:               expectedCountry,
				Name:                  expectedName,
			},
			CreatedOn:      expectedCreatedOnTime,
			ID:             expectedId.String(),
			ModifiedOn:     expectedModifiedOnTime,
			OrganisationID: expectedOrganisationId,
			Type:           expectedType,
			Version:        expectedVersion,
		},
		Links: &accounts.Links{
			Self: expectedSelf,
		},
	}
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
	AccountClassification string   `json:"account_classification,omitempty"`
	AlternativeNames      []string `json:"alternative_bank_account_names,omitempty"`
	BankID                string   `json:"bank_id,omitempty"`
	BankIDCode            string   `json:"bank_id_code,omitempty"`
	BaseCurrency          string   `json:"base_currency,omitempty"`
	Bic                   string   `json:"bic,omitempty"`
	Country               string   `json:"country,omitempty"`
	Name                  []string `json:"name,omitempty"`
}

type Links struct {
	Self string `json:"self" gorm:"type:self"`
}
