package main

import (
	"context"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/config"
	"github.com/TheLonger011/ReLon/internal/database"
	"github.com/TheLonger011/ReLon/internal/handler"
	"github.com/TheLonger011/ReLon/internal/repository"
	"github.com/TheLonger011/ReLon/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}
	log.Printf("Config: %+v", cfg)

	ctx := context.Background()
	conn, err := database.ConnectDB(ctx, &cfg.DB)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	log.Println("successfully connected to database")

	defer conn.Close()

	r := chi.NewRouter()

	repo := repository.NewUserRepository(conn)
	serv := service.NewAuthService(repo)
	hand := handler.NewAuthHandler(serv)

	r.Post("/register", hand.Register)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

}
