package test

import (
	uuid "github.com/satori/go.uuid"
)

func iRemovedTheAccountWithUUID(u string) error {
	err := iTryToRemoveAnAccountWithUUID(u)
	if err != nil {
		return err
	}
	return theResponseIsSuccess()
}

func iTryToRemoveAnAccountWithUUID(u string) error {
	id, err := uuid.FromString(u)
	if err != nil {
		return err
	}
	scenarioState.lastError = f.c.Delete(id, 0)
	return nil
}

// The account after creation has version equal to 0,
// but we try to remove this entity providing another version number
func iTryToRemoveTheModifiedAccountWithUUID(u string) error {
	id, err := uuid.FromString(u)
	if err != nil {
		return err
	}
	scenarioState.lastError = f.c.Delete(id, 1)
	return nil
}
