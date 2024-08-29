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

// unmarshalParams  unmarshalled should be a pointer to one of the *AuthParams structs defined in this package
func (a *Auth) unmarshalParams(unmarshalled any) error {
	if err := json.Unmarshal(a.Params, &unmarshalled); err != nil {
		return fmt.Errorf("error unmarshalling %s params: %w", a.Type, err)
	}
	return nil
}

// UnmarshallParams if no error, return value type will either be BasicAuthParams or BearerAuthParams
func (a *Auth) UnmarshallParams() (any, error) {
	if a == nil {
		return nil, nil
	}
	switch a.Type {
	case BasicAuthType:
		var params BasicAuthParams
		if err := a.unmarshalParams(&params); err != nil {
			return nil, err
		}
		return params, nil
	case BearerType:
		var params BearerAuthParams
		if err := a.unmarshalParams(&params); err != nil {
			return nil, err
		}
		return params, nil
	default:
		return nil, fmt.Errorf("unknown auth type: %s", a.Type)
	}
}

func (a *Auth) SetAuthentication(request *http.Request) error {
	authParams, err := a.UnmarshallParams()
	if err != nil {
		return err
	}
	switch p := authParams.(type) {
	case nil:
		// No auth, so nothing to set
		return nil
	case BasicAuthParams:
		request.SetBasicAuth(p.Username, p.Password)
	case BearerAuthParams:
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Token))
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
