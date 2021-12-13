package core

type ProductStockResponse struct {
	StatusCode int32        `json:"status_code"`
	Message    string       `json:"message"`
	Data       ProductStock `json:"data"`
}
