package unitTests

import (
	"avitoMerch/internal/config"
	"database/sql"
	"log"
)

type TestUser struct {
	Username string
	Password string
	Token    string
}

func setupTestDB() (*sql.DB, error) {
	cfg := config.LoadConfig()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Printf("Не удалось подключиться к базе данных: %v", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Printf("Не удалось проверить подключение к базе данных: %v", err)
		return nil, err
	}
	log.Println("Успешно подключено к базе данных")
	return db, nil
}
