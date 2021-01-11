package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/TOIFLMSC/test-app-Avito-MerchantX/internal/app/store/sqlstore"
)

// Start func
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println("DB opened")

	store := sqlstore.New(db)

	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}

// newDB func
func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
