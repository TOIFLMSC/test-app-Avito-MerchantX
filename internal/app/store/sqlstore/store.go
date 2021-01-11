package sqlstore

import (
	"database/sql"

	"github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/store"

	_ "github.com/lib/pq"
)

// Store struct
type Store struct {
	db              *sql.DB
	offerRepository *OfferRepository
}

// New func
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Offer func
func (s *Store) Offer() store.OfferRepository {
	if s.offerRepository != nil {
		return s.offerRepository
	}

	s.offerRepository = &OfferRepository{
		store: s,
	}

	return s.offerRepository
}
