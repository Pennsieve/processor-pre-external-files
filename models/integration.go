package models

type Integration struct {
	Uuid          string            `json:"uuid"`
	ApplicationID int64             `json:"applicationId"`
	PackageIDs    []string          `json:"packageIds"`
	Params        IntegrationParams `json:"params"`
	Workflow      any               `json:"workflow"`
}

type IntegrationParams struct {
	ExternalFiles []ExternalFileParam `json:"externalFiles"`
}
