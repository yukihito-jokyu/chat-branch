package model

// APIレスポンスの標準フォーマット
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
