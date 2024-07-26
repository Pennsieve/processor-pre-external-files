package models

import "log/slog"

type ExternalFileParam struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Auth  *Auth  `json:"auth,omitempty"`
	Query Query  `json:"query,omitempty"`
}

type Query map[string]string

func (e ExternalFileParam) Logger(logger *slog.Logger) *slog.Logger {
	authType := "none"
	if e.Auth != nil {
		authType = e.Auth.Type
	}

	return logger.With(slog.Group("externalFileParam",
		slog.String("url", e.URL),
		slog.String("name", e.Name),
		slog.String("authType", authType),
		slog.Any("query", e.Query),
	))
}
