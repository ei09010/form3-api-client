package integration

import (
	"context"
	"ei09010/form3-api-client/accounts"
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

var envVar = &EnvVar{}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {

	envVar.InitEnvVariables()

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		envVar.DatabaseHostUrl, envVar.DatabasePort, envVar.DatabaseUser, envVar.DatabaseName, envVar.DatabasePwd)

	conn, err := gorm.Open("postgres", dbUri)

	if err != nil {
		fmt.Print(err)
	}

	s.dbConn = conn

	fmt.Printf(" > TestSuite Setup is complete with the following connection to DB: \n %s", fmt.Sprint(dbUri))
}

func (s *e2eTestSuite) TearDownSuite() {

	s.dbConn.Delete(&Account{})

	fmt.Printf(" > TestSuite TearDown is complete")

}

func (s *e2eTestSuite) SetupTest() {

	s.Require().NoError(s.dbConn.DB().Ping())

	s.dbConn.Delete(&Account{})
}

// Fetch
func (s *e2eTestSuite) TestFetch_FetchesAccount_ReturnsAccount() {

	// Arrange

	accountsClient, err := accounts.NewClient(accounts.WithBaseURL(envVar.ApplicationUrl))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	storedTestAccount := generateAccountDataToStore(id)

	s.NoError(s.dbConn.Create(storedTestAccount).Error)

	expectedAccountData := generatedExpectedAccountToBeReturnedByAPI(id)

	ctx := context.Background()

	// Act

	fetchedAccountData, err := accountsClient.Fetch(ctx, id)

	s.Require().NoError(err)

	// Assert

	assertAccountData(s.Suite, expectedAccountData, fetchedAccountData)

}

func (s *e2eTestSuite) TestFetch_FetchesNonExistentAccount_Returns404Error() {

	// Arrange

	accountsClient, err := accounts.NewClient(accounts.WithBaseURL(envVar.ApplicationUrl))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	ctx := context.Background()

	// Act

	fetchedAccountData, err := accountsClient.Fetch(ctx, id)

	// Assert

	assert.Nil(s.T(), fetchedAccountData, "Fetched account data should be nil")

	assert.Equal(s.T(), fmt.Sprintf("Error message returned by the API | 404 | record %v does not exist", id), err.Error())
}

// Create
func (s *e2eTestSuite) TestCreate_CreatesAccount_ReturnsAccountCreated() {

	// Arrange

	accountsClient, err := accounts.NewClient(accounts.WithBaseURL(envVar.ApplicationUrl))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	accountDataToStore := generatedExpectedAccountToBeReturnedByAPI(id)

	ctx := context.Background()

	// Act

	storedAccountData, err := accountsClient.Create(ctx, accountDataToStore)

	s.Require().NoError(err)

	// Assert

	assertAccountData(s.Suite, accountDataToStore, storedAccountData)
}

func (s *e2eTestSuite) TestCreate_CreatesDuplicateAccount_Returns409Conflict() {

	// Arrange
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL(envVar.ApplicationUrl))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	accountDataToStore := generatedExpectedAccountToBeReturnedByAPI(id)

	storedTestAccount := generateAccountDataToStore(id)

	s.NoError(s.dbConn.Create(storedTestAccount).Error)

	ctx := context.Background()

	// Act

	storedAccountData, err := accountsClient.Create(ctx, accountDataToStore)

	// Assert

	assert.Nil(s.T(), storedAccountData, "Stored account data returned should be nil")

	assert.Equal(s.T(), "Error message returned by the API | 409 | Account cannot be created as it violates a duplicate constraint", err.Error(), "Error message didn't match the expected")
}

// Delete
func (s *e2eTestSuite) TestDelete_DeleteAccount_ReturnsNilError() {

	// Arrange
	accountsClient, err := accounts.NewClient(accounts.WithBaseURL(envVar.ApplicationUrl))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	storedTestAccount := generateAccountDataToStore(id)

	s.NoError(s.dbConn.Create(storedTestAccount).Error)

	expectedVersion := 0

	ctx := context.Background()

	// Act

	err = accountsClient.Delete(ctx, id, expectedVersion)

	s.Require().NoError(err)

}

func (s *e2eTestSuite) TestDelete_DeleteANonExistentccount_Returns404Error() {

	// Arrange

	accountsClient, err := accounts.NewClient(accounts.WithBaseURL(envVar.ApplicationUrl))

	s.Require().NoError(err)

	id, err := uuid.NewUUID()

	s.Require().NoError(err)

	expectedVersion := 0

	ctx := context.Background()

	// Act

	err = accountsClient.Delete(ctx, id, expectedVersion)

	assert.Equal(s.T(), "Error message returned by the API | 404 | ", err.Error(), "Error message didn't match the expected")

}
