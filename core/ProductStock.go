package core

type ProductStock struct {
	ProductId string `json:"product_id"`
	Deposit   string `json:"deposit" validate:"min=11,regexp=^([a-zA-Z]{2\\,}-(\\d{2\\,})-(\\d{2\\,})-(DE|IZ))$"`
	Location  string `json:"location" validate:"min=4"`
	Quantity  int64  `json:"quantity"`
}
