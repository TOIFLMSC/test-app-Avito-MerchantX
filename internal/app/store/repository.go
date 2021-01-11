package store

import "github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/model"

// OfferRepository interface
type OfferRepository interface {
	NewOffer(*model.Offer) error
	UpdateOffer(*model.Offer) error
	DeleteOffer(string, string) error
	GetOfferWithAllSpecs(string, string, string) (*[]model.Offer, error)
	GetOfferWithIDandName(string, string) (*[]model.Offer, error)
	GetOfferWithSelandName(string, string) (*[]model.Offer, error)
	GetOfferWithSelandID(string, string) (*[]model.Offer, error)
	GetOfferWithName(string) (*[]model.Offer, error)
	GetOfferWithID(string) (*[]model.Offer, error)
	GetOfferWithSel(string) (*[]model.Offer, error)
	CheckForOffer(string, string) (bool, error)
}
