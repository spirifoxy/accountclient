package accounts

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	acc := NewAccount("123", "GB", []string{"Samantha Holder"}, &AccountAttributes{BankID: "xxx"})

	_, err := uuid.FromString(acc.Data.ID)
	require.Nil(t, err)

	assert.Equal(t, "123", acc.Data.OrganisationID)
	assert.Equal(t, "accounts", acc.Data.Type)
}

func TestNewAccountRequiredParameters(t *testing.T) {
	acc := NewAccount("123", "XX", []string{"Samantha Holder"}, &AccountAttributes{Country: "GB"})

	assert.Equal(t, "XX", acc.Data.Attributes.Country)
}

func TestNewAccountWithoutAttributes(t *testing.T) {
	acc := NewAccount("123", "GB", []string{"Samantha Holder"}, nil)

	assert.IsType(t, &AccountAttributes{}, acc.Data.Attributes)
	assert.Equal(t, "GB", acc.Data.Attributes.Country)
}
