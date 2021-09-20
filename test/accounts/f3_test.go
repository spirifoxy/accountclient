package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cucumber/godog"
	_ "github.com/lib/pq"
	"github.com/spirifoxy/accountclient/pkg/accounts"
	"github.com/stretchr/testify/assert"
)

var f *clientFeature

type clientFeature struct {
	c   *accounts.Client
	dsn string
	db  *sql.DB
}

// assertExpectedAndActual is a helper function to allow the step function to call
// assertion functions where you want to compare an expected and an actual value.
func assertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	a(&t, expected, actual, msgAndArgs...)
	return t.err
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

// asserter is used to be able to retrieve the error reported by the called assertion
type asserter struct {
	err error
}

// Errorf is used by the called assertion to report an error
func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}

// getEnv tries to get environment varialbe,
// returns fallback in case of failure
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// initSuite inilitializes api client for testing
func initSuite() {
	f = &clientFeature{}

	baseURL := getEnv("API_ADDR", "http://localhost:8080")
	basePath := getEnv("BASE_PATH", "v1")
	timeout, err := strconv.Atoi(getEnv("REQ_TIMEOUT", "5"))
	if err != nil {
		timeout = 5
	}

	postgresHost := getEnv("PSQL_HOST", "localhost")
	postgresPort := getEnv("PSQL_PORT", "5432")
	postgresUser := getEnv("PSQL_USER", "interview_accountapi_user")
	postgresPass := getEnv("PSQL_PASSWORD", "123")
	postgresDB := getEnv("PSQL_DB", "interview_accountapi")

	d := time.Duration(timeout) * time.Second
	f.dsn = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser,
		postgresPass,
		postgresHost,
		postgresPort,
		postgresDB,
	)
	f.c, err = accounts.New(
		baseURL,
		accounts.WithBasePath(basePath),
		accounts.WithTimeout(d),
	)

	if err != nil {
		log.Fatal(fmt.Errorf("not able to create client: %v", err))
	}
}

// TestMain allows us to run e2e tests using go test command,
// which might be useful for debugging or running all the tests
// with only one command
func TestMain(m *testing.M) {
	opts := godog.Options{
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}

	status := godog.TestSuite{
		Name:                 "f3tests",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}

func InitializeTestSuite(sc *godog.TestSuiteContext) {
	sc.BeforeSuite(initSuite)
}

type State struct {
	lastError error
	subject   interface{}
}

var scenarioState *State

func InitializeScenario(sc *godog.ScenarioContext) {
	// Set up db connection and response recorder
	sc.BeforeScenario(func(sc *godog.Scenario) {
		db, err := sql.Open("postgres", f.dsn)
		if err != nil {
			log.Fatal(fmt.Errorf("not able to connect to db: %v", err))
		}

		f.db = db
		scenarioState = &State{}
	})

	// Clean up the db and variables used between scenarios
	sc.AfterScenario(func(sc *godog.Scenario, err error) {
		if f.db != nil {
			f.db.Exec(`DELETE FROM public."Account"`)
			f.db.Close()
		}

		scenarioState = nil
	})

	sc.Step(`^the response is success$`, theResponseIsSuccess)
	sc.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	sc.Step(`^I get the request status error$`, iGetTheRequestStatusError)
	sc.Step(`^I get the account$`, iGetTheAccount)

	// health
	sc.Step(`^I call IsHealthy method$`, iCallIsHealthyMethod)
	sc.Step(`^I should receive true as response$`, iShouldReceiveTrueAsResponse)

	// fetch
	sc.Step(`^the account exists with uuid (.*)$`, theAccountExistsWithUUID)
	sc.Step(`^I try to fetch the account with uuid (.*)$`, iTryToFetchTheAccountWithUUID)

	// create
	sc.Step(`^I created an account with uuid (.*)$`, iCreatedAnAccountWithUUID)
	sc.Step(`^I try to create the account with uuid (.*)$`, iTryToCreateTheAccountWithUUID)
	sc.Step(`^I accidentally try to create an account with no name$`, iAccidentallyTryToCreateAnAccountWithoutName)
	sc.Step(`^I accidentally try to create an account with broken country$`, iAccidentallyTryToCreateAnAccountWithBrokenCountry)
	sc.Step(`^I accidentally try to create an account with broken organization id$`, iAccidentallyTryToCreateAnAccountWithBrokenOrganizationID)

	// remove
	sc.Step(`^I removed the account with uuid (.*)$`, iRemovedTheAccountWithUUID)
	sc.Step(`^I try to remove an account with uuid (.*)$`, iTryToRemoveAnAccountWithUUID)
	sc.Step(`^I try to remove the modified account with uuid (.*)$`, iTryToRemoveTheModifiedAccountWithUUID)
}

func iGetTheAccount() error {
	return assertExpectedAndActual(
		assert.IsType, &accounts.Account{}, scenarioState.subject,
		"expected account not to be nil, but it is",
	)
}

func iGetTheRequestStatusError() error {
	return assertExpectedAndActual(
		assert.IsType, &accounts.RequestStatusError{}, scenarioState.lastError,
		"expected error to be RequestStatusError type",
	)
}

func theResponseIsSuccess() error {
	return assertExpectedAndActual(
		assert.IsType, nil, scenarioState.lastError,
		"expected response to be successfull, but error happened: %v", scenarioState.lastError,
	)
}

func theResponseCodeShouldBe(code int) error {
	err, ok := scenarioState.lastError.(*accounts.RequestStatusError)
	if !ok {
		return fmt.Errorf("last error is not request status error: %v", scenarioState.lastError)
	}

	return assertExpectedAndActual(
		assert.Equal, code, err.Code,
		"expected status code to be %d, but it is %d", code, err.Code,
	)
}
