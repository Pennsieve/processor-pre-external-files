package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const basicAuthType = "basicAuth"
const bearerType = "bearer"

type Auth struct {
	Type string `json:"type"`
	// Params are left unmarshalled until SetAuthentication is called
	Params json.RawMessage `json:"params"`
}

func (a *Auth) unmarshalParams(unmarshalled any) error {
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
	case basicAuthType:
		var params BasicAuthParams
		if err := a.unmarshalParams(&params); err != nil {
			return err
		}
		request.SetBasicAuth(params.Username, params.Password)
	case bearerType:
		var params BearerAuthParams
		if err := a.unmarshalParams(&params); err != nil {
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
