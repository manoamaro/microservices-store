package models

type ErrorResponse struct {
	Status string `json:"status"`
}

type VerifyResponse struct {
	Audiences []string `json:"audiences"`
	Flags     []string `json:"flags"`
}
