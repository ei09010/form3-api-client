package integration

import (
	"ei09010/form3-api-client/organisation/accounts"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
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

// func TestMain(m *testing.M) {
// 	setup()
// 	code := m.Run()
// 	teardown()
// 	os.Exit(code)
// }

// func setup() {
// 	// Do something here.

// 	fmt.Printf("\033[1;36m%s\033[0m", "> Setup completed\n")
// }

// func teardown() {
// 	// Do something here.

// 	fmt.Printf("\033[1;36m%s\033[0m", "> Teardown completed")
// 	fmt.Printf("\n")
// }

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		"localhost", "5432", "root", "interview_accountapi", "password")
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	s.dbConn = conn
}

// Fetch
func (s *e2eTestSuite) TestFetch_FetchesAccount_ReturnsAccount() {

	// Arrange
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL("http://localhost:8080"))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	storedTestAccount := generateAccountDataToStore(id)

	s.NoError(s.dbConn.Create(storedTestAccount).Error)

	expectedAccountData := generatedExpectedAccountToBeReturnedByAPI(id)

	// Act

	fetchedAccountData, err := accountsClient.Fetch(id)

	s.Require().NoError(err)

	// Assert

	assertAccountData(s.Suite, expectedAccountData, fetchedAccountData)

}

func (s *e2eTestSuite) TestFetch_FetchesNonExistentAccount_Returns404Error() {

	// Arrange
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL("http://localhost:8080"))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	// Act

	fetchedAccountData, err := accountsClient.Fetch(id)

	// Assert

	assert.Nil(s.T(), fetchedAccountData, "Fetched account data should be nil")

	assert.Equal(s.T(), fmt.Sprintf("Error message returned by the API | 404 | record %v does not exist", id), err.Error())
}

// Create
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

	assertAccountData(s.Suite, accountDataToStore, storedAccountData)
}

func (s *e2eTestSuite) TestCreate_CreatesDuplicateAccount_Returns409Conflict() {

	// Arrange

	// this url has to be a env variable
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL("http://localhost:8080"))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	accountDataToStore := generatedExpectedAccountToBeReturnedByAPI(id)

	storedTestAccount := generateAccountDataToStore(id)

	s.NoError(s.dbConn.Create(storedTestAccount).Error)

	// Act

	storedAccountData, err := accountsClient.Create(accountDataToStore)

	// Assert

	assert.Nil(s.T(), storedAccountData, "Stored account data returned should be nil")

	assert.Equal(s.T(), "Error message returned by the API | 409 | Account cannot be created as it violates a duplicate constraint", err.Error(), "Error message didn't match the expected")
}

// Delete
func (s *e2eTestSuite) TestDelete_DeleteAccount_ReturnsNilError() {

	// Arrange

	// this url has to be a env variable
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL("http://localhost:8080"))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	storedTestAccount := generateAccountDataToStore(id)

	s.NoError(s.dbConn.Create(storedTestAccount).Error)

	expectedVersion := 0

	// Act

	err = accountsClient.Delete(id, expectedVersion)

	s.Require().NoError(err)

}

func (s *e2eTestSuite) TestDelete_DeleteANonExistentccount_Returns404Error() {

	// Arrange

	// this url has to be a env variable
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL("http://localhost:8080"))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	expectedVersion := 0

	// Act

	err = accountsClient.Delete(id, expectedVersion)

	assert.Equal(s.T(), "Error message returned by the API | 404 | ", err.Error(), "Error message didn't match the expected")

}
