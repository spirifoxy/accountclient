package test

import (
	uuid "github.com/satori/go.uuid"
	"github.com/spirifoxy/accountclient/pkg/accounts"
)

func buildAccountWithUUID(u string) *accounts.Account {
	// override autogenerated uuid just for the simplicity of testing
	acc := accounts.NewAccount(u, "GB", []string{"Samantha Holder"}, nil)
	acc.Data.ID = u
	return acc
}

func buildAccountWithoutName() *accounts.Account {
	return accounts.NewAccount(uuid.NewV4().String(), "GB", nil, nil)
}

func buildAccountWithBrokenCountry() *accounts.Account {
	return accounts.NewAccount(uuid.NewV4().String(), "hi", []string{"Samantha Holder"}, nil)
}

func buildAccountWithBrokenOrganizationID() *accounts.Account {
	return accounts.NewAccount("somebrokenuuid", "GB", []string{"Samantha Holder"}, nil)
}

func iTryToCreateTheAccountWithUUID(uuid string) error {
	acc := buildAccountWithUUID(uuid)
	scenarioState.subject, scenarioState.lastError = f.c.Create(acc)
	return nil
}

func iAccidentallyTryToCreateAnAccountWithoutName() error {
	acc := buildAccountWithoutName()
	scenarioState.subject, scenarioState.lastError = f.c.Create(acc)
	return nil
}

func iAccidentallyTryToCreateAnAccountWithBrokenCountry() error {
	acc := buildAccountWithBrokenCountry()
	scenarioState.subject, scenarioState.lastError = f.c.Create(acc)
	return nil
}

func iAccidentallyTryToCreateAnAccountWithBrokenOrganizationID() error {
	acc := buildAccountWithBrokenOrganizationID()
	scenarioState.subject, scenarioState.lastError = f.c.Create(acc)
	return nil
}

func iCreatedAnAccountWithUUID(uuid string) error {
	_ = iTryToCreateTheAccountWithUUID(uuid)
	return iGetTheAccount()
}