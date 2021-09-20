package test

import (
	"github.com/stretchr/testify/assert"
)

func iCallIsHealthyMethod() error {
	scenarioState.subject, scenarioState.lastError = f.c.IsHealthy()
	if scenarioState.lastError != nil {
		return scenarioState.lastError
	}

	return nil
}

func iShouldReceiveTrueAsResponse() error {
	return assertExpectedAndActual(
		assert.Equal, true, scenarioState.subject,
		"expected IsHealthy response to be true, but it is %v", scenarioState.subject,
	)
}
