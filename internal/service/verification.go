package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"math/big"
	"time"
)

type VerificationRepo interface {
	MarkVerified(ctx context.Context, userID uuid.UUID) error
	GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error)
}

type VerificationService struct {
	redis        *redis.Client
	emailService *EmailService
	userRepo     VerificationRepo
}

func generateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

func NewVerificationService(redis *redis.Client, emailService *EmailService, userRepo VerificationRepo) *VerificationService {
	return &VerificationService{redis: redis, emailService: emailService, userRepo: userRepo}
}

func (s *VerificationService) SendVerificationCode(ctx context.Context, email string) error {
	code := generateCode()
	err := s.redis.Set(ctx, "verify:"+email, code, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("send verification code: %w", err)
	}
	if err := s.emailService.SendVerificationCode(email, code); err != nil {
		return fmt.Errorf("send verification code: %w", err)
	}
	return nil
}

func (s *VerificationService) VerifyCode(ctx context.Context, email, code string) error {
	storedCode, err := s.redis.Get(ctx, "verify:"+email).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("verify code expired or not found")
	}
	if err != nil {
		return fmt.Errorf("get verify code: %w", err)
	}
	if storedCode != code {
		return fmt.Errorf("verification code does not match")
	}
	user, err := s.userRepo.GetByEmailOrUsername(ctx, email)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if err := s.userRepo.MarkVerified(ctx, user.ID); err != nil {
		return fmt.Errorf("mark verified: %w", err)
	}
	if err := s.redis.Del(ctx, "verify:"+email).Err(); err != nil {
		return fmt.Errorf("delete verification code: %w", err)
	}
	return nil
}
