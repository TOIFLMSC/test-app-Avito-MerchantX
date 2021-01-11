package sqlstore

import (
	"database/sql"

	"github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/model"
	"github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/store"
	u "github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/utils"
)

// OfferRepository struct
type OfferRepository struct {
	store *Store
}

// NewOffer func
func (r *OfferRepository) NewOffer(o *model.Offer) error {
	return r.store.db.QueryRow("INSERT INTO offers (seller, offer_id, name, price, quantity, available) VALUES ($1, $2, $3, $4, $5, $6) RETURNING seller",
		o.Seller,
		o.OfferID,
		o.Name,
		o.Price,
		o.Quantity,
		o.Available,
	).Scan(&o.Seller)
}

// GetOfferWithAllSpecs func
func (r *OfferRepository) GetOfferWithAllSpecs(seller string, offerid string, name string) (*[]model.Offer, error) {
	rows, err := r.store.db.Query(
		"SELECT seller, offer_id, name, price, quantity, available FROM offers WHERE seller = $1 AND offer_id = $2 AND name LIKE $3",
		seller,
		offerid,
		name+"%",
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	array, err := u.ConvertRowToArray(rows)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

// GetOfferWithIDandName func
func (r *OfferRepository) GetOfferWithIDandName(offerid string, name string) (*[]model.Offer, error) {
	rows, err := r.store.db.Query(
		"SELECT seller, offer_id, name, price, quantity, available FROM offers WHERE offer_id = $1 AND name LIKE $2",
		offerid,
		name+"%",
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	array, err := u.ConvertRowToArray(rows)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

// GetOfferWithSelandName func
func (r *OfferRepository) GetOfferWithSelandName(seller string, name string) (*[]model.Offer, error) {
	rows, err := r.store.db.Query(
		"SELECT seller, offer_id, name, price, quantity, available FROM offers WHERE seller = $1 AND name LIKE $2",
		seller,
		name+"%",
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	array, err := u.ConvertRowToArray(rows)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

// GetOfferWithSelandID func
func (r *OfferRepository) GetOfferWithSelandID(seller string, offerid string) (*[]model.Offer, error) {
	rows, err := r.store.db.Query(
		"SELECT seller, offer_id, name, price, quantity, available FROM offers WHERE seller = $1 AND offer_id = $2",
		seller,
		offerid,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	array, err := u.ConvertRowToArray(rows)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

// GetOfferWithName func
func (r *OfferRepository) GetOfferWithName(name string) (*[]model.Offer, error) {
	rows, err := r.store.db.Query(
		"SELECT seller, offer_id, name, price, quantity, available FROM offers WHERE name LIKE $1",
		name+"%",
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	array, err := u.ConvertRowToArray(rows)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

// GetOfferWithID func
func (r *OfferRepository) GetOfferWithID(offerid string) (*[]model.Offer, error) {
	rows, err := r.store.db.Query(
		"SELECT seller, offer_id, name, price, quantity, available FROM offers WHERE offer_id = $1",
		offerid,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	array, err := u.ConvertRowToArray(rows)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

// GetOfferWithSel func
func (r *OfferRepository) GetOfferWithSel(seller string) (*[]model.Offer, error) {
	rows, err := r.store.db.Query(
		"SELECT seller, offer_id, name, price, quantity, available FROM offers WHERE seller = $1",
		seller,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	defer rows.Close()

	array, err := u.ConvertRowToArray(rows)
	if err != nil {
		return nil, err
	}

	return &array, nil
}

// CheckForOffer func
func (r *OfferRepository) CheckForOffer(seller string, offerID string) (bool, error) {
	o := &model.Offer{}
	if err := r.store.db.QueryRow(
		"SELECT seller, offer_id FROM offers WHERE seller = $1 AND offer_id = $2",
		seller,
		offerID,
	).Scan(
		&o.Seller,
		&o.OfferID,
	); err != nil {
		if err == sql.ErrNoRows {
			return true, store.ErrRecordNotFound
		}
		return false, err
	}
	return false, nil
}

// UpdateOffer func
func (r *OfferRepository) UpdateOffer(o *model.Offer) error {

	return r.store.db.QueryRow("UPDATE offers SET price = $1, quantity = $2 WHERE seller = $3 AND offer_id = $4 RETURNING offer_id",
		o.Price,
		o.Quantity,
		o.Seller,
		o.OfferID,
	).Scan(&o.OfferID)
}

// DeleteOffer func
func (r *OfferRepository) DeleteOffer(seller string, offerID string) error {

	_, err := r.store.db.Exec("DELETE FROM offers WHERE seller = $1 AND offer_id = $2",
		seller,
		offerID,
	)
	return err
}
