package template

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shipping/dto"
)

type MLTemplate struct {
}

func (t *MLTemplate) GetProduct(productId string) dto.ItemDTO {
	resp, err := http.Get("https://api.mercadolibre.com/items/" + productId)
	var response dto.ItemDTO
	if err != nil {
		return response
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)
	return response
}
