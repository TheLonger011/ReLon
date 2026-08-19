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

	redisClient, err := database.ConnectRedis(ctx, cfg.Redis.Addr)
	if err != nil {
		log.Fatal("Error connecting to redis: ", err)
	}

	defer redisClient.Close()

	IPLimiter := middleware.NewIPLimiter()

	from := cfg.Verification.From
	password := cfg.Verification.Password

	r := chi.NewRouter()

	userRepo := repository.NewUserRepository(conn)
	emailServ := service.NewEmailService(from, password)
	verificationServ := service.NewVerificationService(redisClient, emailServ, userRepo)
	userServ := service.NewAuthService(userRepo, cfg.JWT.Secret)
	userHand := handler.NewAuthHandler(userServ, verificationServ)

	postRepo := repository.NewPostRepository(conn)

	voteRepo := repository.NewVoteRepository(conn)
	voteServ := service.NewVoteService(voteRepo)
	voteHand := handler.NewVoteHandler(voteServ)

	commentRepo := repository.NewCommentRepository(conn)
	commentServ := service.NewCommentService(commentRepo)
	commentHand := handler.NewCommentHandler(commentServ)

	communityRepo := repository.NewCommunityRepository(conn)
	communityServ := service.NewCommunityService(communityRepo)
	communityHand := handler.NewCommunityHandler(communityServ)

	communityMemberRepo := repository.NewCommunityMemberRepository(conn)
	communityMemberServ := service.NewCommunityMemberService(communityMemberRepo, communityRepo)
	communityMemberHand := handler.NewCommunityMemberHandler(communityMemberServ)

	postSer := service.NewPostService(postRepo, communityMemberRepo)
	postHand := handler.NewPostHandler(postSer)

	communityRequestRepo := repository.NewCommunityRequestRepository(conn)
	communityJoinRequestServ := service.NewCommunityJoinRequestService(communityRequestRepo, communityMemberRepo, communityRepo)
	communityJoinRequestHand := handler.NewCommunityJoinRequestHandler(communityJoinRequestServ)

	r.With(middleware.RateLimit(IPLimiter)).Post("/register", userHand.Register)
	r.With(middleware.RateLimit(IPLimiter)).Post("/login", userHand.Login)
	r.With(middleware.Auth(cfg.JWT.Secret)).Get("/me", userHand.Me)
	r.Get("/posts/search", postHand.SearchPosts)
	r.With(middleware.Auth(cfg.JWT.Secret)).Delete("/posts/{id}", postHand.DeletePost)
	r.With(middleware.Auth(cfg.JWT.Secret)).Put("/posts/{id}", postHand.UpdatePost)
	r.With(middleware.RateLimit(IPLimiter)).Post("/verify", userHand.Verify)

	r.With(middleware.RateLimit(IPLimiter), middleware.Auth(cfg.JWT.Secret)).Post("/posts", postHand.CreatePost)
	r.Get("/posts/{id}", postHand.GetByPostID)
	r.Get("/posts", postHand.GetPosts)

	r.With(middleware.Auth(cfg.JWT.Secret)).Post("/posts/{id}/vote", voteHand.Vote)
	r.With(middleware.Auth(cfg.JWT.Secret)).Delete("/posts/{id}/vote", voteHand.RemoveVote)

	r.With(middleware.RateLimit(IPLimiter), middleware.Auth(cfg.JWT.Secret)).Post("/posts/{id}/comments", commentHand.CreateComment)
	r.Get("/posts/{id}/comments", commentHand.GetByPostID)
	r.With(middleware.Auth(cfg.JWT.Secret)).Delete("/comments/{id}", commentHand.DeleteComment)
	r.With(middleware.Auth(cfg.JWT.Secret)).Put("/comments/{id}", commentHand.UpdateComment)

	r.With(middleware.Auth(cfg.JWT.Secret)).Post("/communities", communityHand.CreateCommunity)
	r.Get("/communities/{id}", communityHand.GetCommunityByID)
	r.Get("/communities", communityHand.GetCommunities)
	r.With(middleware.Auth(cfg.JWT.Secret)).Delete("/communities/{id}", communityHand.DeleteCommunity)
	r.With(middleware.Auth(cfg.JWT.Secret)).Put("/communities/{id}", communityHand.UpdateCommunity)
	r.With(middleware.Auth(cfg.JWT.Secret)).Delete("/communities/{id}/leave", communityMemberHand.LeaveCommunity)
	r.With(middleware.Auth(cfg.JWT.Secret)).Post("/communities/{id}/join", communityJoinRequestHand.JoinCommunity)
	r.With(middleware.Auth(cfg.JWT.Secret)).Get("/communities/{id}/requests", communityJoinRequestHand.GetPendingRequests)
	r.With(middleware.Auth(cfg.JWT.Secret)).Post("/communities/{id}/requests/{requestId}/approve", communityJoinRequestHand.ApproveJoinRequest)
	r.With(middleware.Auth(cfg.JWT.Secret)).Post("/communities/{id}/requests/{requestId}/reject", communityJoinRequestHand.RejectJoinRequest)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
