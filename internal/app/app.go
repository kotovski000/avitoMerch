package app

import (
	"avitoMerch/internal/config"
	"avitoMerch/internal/handler"
	"avitoMerch/internal/middleware"
	"avitoMerch/internal/repository"
	"avitoMerch/internal/service"
	"avitoMerch/internal/utils"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
)

type App struct {
	config             *config.Config
	Db                 *sql.DB
	authHandler        *handler.AuthHandler
	infoHandler        *handler.InfoHandler
	itemHandler        *handler.ItemHandler
	transactionHandler *handler.TransactionHandler
	authMiddleware     *middleware.AuthMiddleware
	jwtService         *utils.JWTService
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := OpenDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	jwtService := utils.NewJWTService(cfg.JWTSecretKey)

	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	authService := service.NewAuthService(userRepo, jwtService, cfg)
	itemService := service.NewItemService(itemRepo, userRepo, transactionRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)

	authMiddleware := middleware.NewAuthMiddleware(jwtService, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	infoHandler := handler.NewInfoHandler(transactionService, itemService, userRepo)
	itemHandler := handler.NewItemHandler(itemService, transactionService, userRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService, userRepo)

	return &App{
		config:             cfg,
		Db:                 db,
		authHandler:        authHandler,
		infoHandler:        infoHandler,
		itemHandler:        itemHandler,
		transactionHandler: transactionHandler,
		authMiddleware:     authMiddleware,
		jwtService:         jwtService,
	}, nil
}

func (a *App) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/auth", a.authHandler.Authenticate).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(a.authMiddleware.Authenticate)

	apiRouter.HandleFunc("/info", a.infoHandler.GetInfo).Methods("GET")
	apiRouter.HandleFunc("/buy/{item}", a.itemHandler.BuyItem).Methods("GET")
	apiRouter.HandleFunc("/sendCoin", a.transactionHandler.SendCoin).Methods("POST")
}

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.DatabaseURL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Подключено к базе данных")
	return db, nil
}
