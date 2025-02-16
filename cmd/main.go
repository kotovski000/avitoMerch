package main

import (
	"avitoMerch/internal/app"
	"avitoMerch/internal/config"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Ошибка при инициализации приложения: %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}).Methods("GET")

	application.RegisterRoutes(router)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Сервер слушает на %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
