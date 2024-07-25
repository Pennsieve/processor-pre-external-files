package models

type Integration struct {
	Uuid          string               `json:"uuid"`
	ApplicationID int64                `json:"applicationId"`
	PackageIDs    []string             `json:"packageIds"`
	Params        []ExternalFilesParam `json:"params"`
	Workflow      any                  `json:"workflow"`
}
