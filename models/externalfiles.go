package models

type ExternalFileParam struct {
	URL     string  `json:"url"`
	Auth    *Auth   `json:"auth,omitempty"`
	Queries []Query `json:"queries,omitempty"`
}

type Query map[string]string
