package accounts

import uuid "github.com/satori/go.uuid"

// Account represents an account in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.
type Account struct {
	Data *AccountData `json:"data"`
}

// AccountData represents all the account data.
type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        int64              `json:"version,omitempty"`
}

// AccountAttributes represents attributes in the account data.
type AccountAttributes struct {
	AccountClassification   string   `json:"account_classification,omitempty"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 string   `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            bool     `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  string   `json:"status,omitempty"`
	Switched                bool     `json:"switched,omitempty"`
}

// NewAccount is a basic builder for the account entity.
//
// Client side validation is omitted on purpose.
// Only the fields that are always required by docs are mandatory
// and will override those in attributes if attributes are provided.
func NewAccount(organizationID, country string, name []string, attr *AccountAttributes) *Account {
	if attr == nil {
		attr = &AccountAttributes{}
	}
	attr.Country = country
	attr.Name = name

	return &Account{
		Data: &AccountData{
			ID:             uuid.NewV4().String(),
			OrganisationID: organizationID,
			Type:           "accounts",
			Attributes:     attr,
		},
	}
}

// Health represents api status
type Health struct {
	Status string `json:"status"`
}

// ErrorResponse represents an API error returned in case of failure.
type ErrorResponse struct {
	Message string `json:"error_message"`
}
