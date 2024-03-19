package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const baseURI string = "https://api.scryfall.com/cards/named?fuzzy="

type Card struct {
	Name      string `json:"name"`
	ManaCost  string `json:"mana_cost"`
	ImageUris struct {
		Large string `json:"large"`
	} `json:"image_uris"`
	Prices struct {
		Eur string `json:"eur"`
		Usd string `json:"usd"`
	} `json:"prices"`
}

func GetCards(cardName string) (*Card, error) {
	url := baseURI + cardName
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var card Card
	err = json.NewDecoder(resp.Body).Decode(&card)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func main() {
	// Get details for a specific card
	card, _ := GetCards("jace sculptor")

	// Print the details
	log.Println(card.Name, card.ManaCost, card.ImageUris.Large, card.Prices.Usd, card.Prices.Eur)
}
