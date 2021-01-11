package model

// Offer type
type Offer struct {
	Seller    string `json:"seller"`
	OfferID   string `json:"offer_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Quantity  int    `json:"quantity"`
	Available string `json:"available"`
}
