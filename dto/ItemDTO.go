package dto

type ItemDTO struct {
	Id                string      `json:"id"`
	CategoryId        string      `json:"category_id"`
	AvailableQuantity int         `json:"available_quantity"`
	Shipping          ShippingDTO `json:"shipping"`
}
