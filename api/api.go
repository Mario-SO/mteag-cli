package api

import (
	"encoding/json"
	"net/http"

	"github.com/mario-so/mteag-cli/item"
)

const baseURI string = "https://api.scryfall.com/cards/named?fuzzy="

func GetCards(cardName string) ([]item.Card, error) {
	url := baseURI + cardName
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var body = resp.Body
	defer resp.Body.Close()

	var data struct {
		Data []item.Card `json:"data"`
	}

	err = json.NewDecoder(body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}
