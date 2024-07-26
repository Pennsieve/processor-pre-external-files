package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const BasicAuthType = "basicAuth"
const BearerType = "bearer"

type Auth struct {
	Type string `json:"type"`
	// Params are left unmarshalled until SetAuthentication is called
	Params json.RawMessage `json:"params"`
}

// UnmarshalParams  visible for testing. unmarshalled should be a pointer to one of the *AuthParams structs defined in this package
func (a *Auth) UnmarshalParams(unmarshalled any) error {
	if err := json.Unmarshal(a.Params, &unmarshalled); err != nil {
		return fmt.Errorf("error unmarshalling %s params: %w", a.Type, err)
	}
	return nil
}

func (a *Auth) SetAuthentication(request *http.Request) error {
	if a == nil {
		return nil
	}
	switch a.Type {
	case BasicAuthType:
		var params BasicAuthParams
		if err := a.UnmarshalParams(&params); err != nil {
			return err
		}
		request.SetBasicAuth(params.Username, params.Password)
	case BearerType:
		var params BearerAuthParams
		if err := a.UnmarshalParams(&params); err != nil {
			return err
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", params.Token))
	default:
		return fmt.Errorf("unknown auth type: %s", a.Type)
	}
	return nil
}

type BasicAuthParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BearerAuthParams struct {
	Token string `json:"token"`
}
