package item

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
