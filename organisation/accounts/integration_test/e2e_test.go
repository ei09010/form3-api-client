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

	// gorm "gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	// "gorm.io/gorm"
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

	//gorm original
	// //"postgres://postgres:postgres@localhost/Account?sslmode=disable"
	// s.dbConnectionStr = "host=localhost user=root password=password port=5432"
	// //s.dbConnectionStr = "host=localhost user=interview_accountapi_user password=123 dbname=interview_accountapi port=5432" //config.DBConnectionURL()
	// db, err := gorm.Open(postgres.Open(s.dbConnectionStr), &gorm.Config{})
	// s.Require().NoError(err)
	// s.dbConn = db

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

	temp := generateAccountDataToStore(id)

	fmt.Println(id.String())

	s.NoError(s.dbConn.Create(temp).Error)

	// Act

	storedAccountData, err := accountsClient.Fetch(id)

	s.Require().NoError(err)

	// Assert

	// assert stored account data
	fmt.Printf("%v", storedAccountData)
}

func (s *e2eTestSuite) TestCreate_CreatesAccount_ReturnsAccountCreated() {

	// Arrange

	accountDataToStore := generateValidGenericAccountData()

	// Act
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL("http://localhost:8080"))

	s.Require().NoError(err)

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

// JSONB Interface for JSONB Field of yourTableName Table
type JSONB []interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
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

	// timeLayout := "2006-01-02 15:04:05 -0700 MST"
	//expectedCreatedOn := "2021-07-31 22:09:02 +0000 UTC"
	//expectedCreatedOnTime, _ := time.Parse(timeLayout, expectedCreatedOn)

	//expectedId := "86b3264e-0121-11ec-9a03-0242ac130003"

	// expectedModifiedOn := "2021-07-31 22:09:02 +0000 UTC"
	// expectedModifiedOnTime, _ := time.Parse(timeLayout, expectedModifiedOn)

	expectedOrganisationId := "fc1bd6f5-c3f5-44b2-b677-acd23cdde73c"
	//expectedType := "accounts"
	expectedVersion := 0
	//expectedSelf := "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

	byteArr, _ := json.Marshal(&AccountAttributes{
		AccountClassification: expectedAccountClassification,
		AlternativeNames:      expectedAlternativeNames,
		BankID:                expectedBankId,
		BankIDCode:            expectedBankIdCode,
		BaseCurrency:          expectedBaseCurrency,
		Bic:                   expectedBic,
		Country:               expectedCountry,
		Name:                  expectedName,
	})

	recordJs := json.RawMessage(string(byteArr))

	return &Account{
		Record: JSONB{recordJs},
		// record: JSONB{
		// 	"account_classification":         expectedAccountClassification,
		// 	"alternative_bank_account_names": expectedAlternativeNames,
		// 	"bank_id":                        expectedBankId,
		// 	"bank_id_code":                   expectedBankIdCode,
		// 	"base_currency":                  expectedBaseCurrency,
		// 	"bic":                            expectedBic,
		// 	"country":                        expectedCountry,
		// 	"name":                           expectedName,
		// }, //postgres.Jsonb{RawMessage: recordJs},
		ID:             id,
		ModifiedOn:     time.Now(),
		IsDeleted:      false,
		IsLocked:       false,
		OrganisationID: expectedOrganisationId,
		Version:        expectedVersion,
	}
}

type Account struct {
	ID             uuid.UUID `gorm:"unique"`
	ModifiedOn     time.Time `json:"modified_on" gorm:"type:modified_on"`
	OrganisationID string    `json:"organisation_id" gorm:"type:organisation_id"`
	Version        int       `json:"version" gorm:"type:version"`
	IsDeleted      bool      `gorm:"type:is_deleted"`
	IsLocked       bool      `gorm:"type:is_locked"`
	Record         JSONB     `gorm:"type:jsonb" json:"record"` //`sql:"type:JSONB"` //JSONB     `gorm:"type:jsonb"`
	//AccountAttributes *AccountAttributes `gorm:"type:record"`
}

// {"bic": "NWBKGB22",
// "name": ["Name of the account holder, up to four lines possible."],
// "bank_id": "400300",
// "country": "GB",
// "bank_id_code": "GBDSC",
// "base_currency": "GBP",
// "account_classification": "Personal",
// "alternative_bank_account_names": null}

func generateValidGenericAccountData() *accounts.AccountData {

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

	expectedId := "be38e265-9607-5c5c-a0e5-3003ea9cc4dd"

	expectedModifiedOn := "2021-07-31 22:09:02 +0000 UTC"
	expectedModifiedOnTime, _ := time.Parse(timeLayout, expectedModifiedOn)

	expectedOrganisationId := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	expectedType := "accounts"
	expectedVersion := 0
	expectedSelf := "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

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
			ID:             expectedId,
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
	gorm.Model
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
	AlternativeNames      []string `json:"alternative_bank_account_names" gorm:"type:alternative_bank_account_names"`
	BankID                string   `json:"bank_id" gorm:"type:bank_id"`
	BankIDCode            string   `json:"bank_id_code" gorm:"type:bank_id_code"`
	BaseCurrency          string   `json:"base_currency" gorm:"type:base_currency"`
	Bic                   string   `json:"bic" gorm:"type:bic"`
	Country               string   `json:"country" gorm:"type:country"`
	Name                  []string `json:"name" gorm:"type:name"`
}

type Links struct {
	gorm.Model
	Self string `json:"self" gorm:"type:self"`
}
