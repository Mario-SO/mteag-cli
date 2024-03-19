package api

import (
	"encoding/json"
	"net/http"

	"github.com/mario-so/mteag-cli/item"
)

const baseURI string = "https://api.scryfall.com/cards/named?fuzzy="

func GetCard(cardName string) (item.Card, error) {
	url := baseURI + cardName
	resp, err := http.Get(url)
	if err != nil {
		return item.Card{}, err
	}

	defer resp.Body.Close()

	var card item.Card

	err = json.NewDecoder(resp.Body).Decode(&card)
	if err != nil {
		return item.Card{}, err
	}

	return card, nil
}
