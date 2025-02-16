package tests

import (
	"avitoMerch/internal/app"
	"avitoMerch/internal/config"
	"database/sql"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

func setupApp() (*mux.Router, *httptest.ResponseRecorder, *sql.DB) {
	cfg := config.LoadConfig()

	a, err := app.NewApp(cfg)
	if err != nil {
		panic("Failed to create app: " + err.Error())
	}

	router := mux.NewRouter()
	a.RegisterRoutes(router)

	return router, httptest.NewRecorder(), a.Db
}
func clearDatabaseAuthTest(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		return err
	}
	return nil
}
func clearDatabaseBuyItemTest(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM transactions")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM inventory")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		return err
	}
	return nil
}

func clearDatabaseInfoTest(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM transactions")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM inventory")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		return err
	}
	return nil
}

func clearDatabaseSendCoinTest(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM transactions")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM inventory")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		return err
	}
	return nil
}
