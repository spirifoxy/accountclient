package test

import (
	uuid "github.com/satori/go.uuid"
)

func theAccountExistsWithUUID(u string) error {
	stmt := `
		INSERT INTO public."Account" (id,organisation_id,"version",is_deleted,is_locked,created_on,modified_on,record)
		VALUES ('` + u + `'::uuid, '` + u + `'::uuid,0,false,false,'2021-01-01 11:11:11.123456','2021-01-01 11:11:11.123456','{"bic": "NWBKGB42", "name": ["Samantha Holder"], "bank_id": "400302", "country": "GB", "bank_id_code": "GBDSC", "base_currency": "GBP", "joint_account": false, "account_classification": "Personal", "account_matching_opt_out": false, "secondary_identification": "A1B2C3D4", "alternative_bank_account_names": ["Sam Holder"]}')
	`
	_, err := f.db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func iTryToFetchTheAccountWithUUID(u string) error {
	id, err := uuid.FromString(u)
	if err != nil {
		return err
	}
	scenarioState.subject, scenarioState.lastError = f.c.Fetch(id)
	return nil
}
