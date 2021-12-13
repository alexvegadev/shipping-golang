package core

type ErrorMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
