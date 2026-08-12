package main

import (
	"context"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/config"
	"github.com/TheLonger011/ReLon/internal/database"
	"github.com/TheLonger011/ReLon/internal/handler"
	"github.com/TheLonger011/ReLon/internal/middleware"
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

	userRepo := repository.NewUserRepository(conn)
	userServ := service.NewAuthService(userRepo, cfg.JWT.Secret)
	userHand := handler.NewAuthHandler(userServ)

	postRepo := repository.NewPostRepository(conn)
	postSer := service.NewPostService(postRepo)
	postHand := handler.NewPostHandler(postSer)

	voteRepo := repository.NewVoteRepository(conn)
	voteServ := service.NewVoteService(voteRepo)
	voteHand := handler.NewVoteHandler(voteServ)

	r.Post("/register", userHand.Register)
	r.Post("/login", userHand.Login)
	r.With(middleware.Auth(cfg.JWT.Secret)).Get("/me", userHand.Me)

	r.With(middleware.Auth(cfg.JWT.Secret)).Post("/posts", postHand.CreatePost)
	r.Get("/posts/{id}", postHand.GetByPostID)
	r.Get("/posts", postHand.GetPosts)

	r.With(middleware.Auth(cfg.JWT.Secret)).Post("/posts/{id}/vote", voteHand.Vote)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

}
