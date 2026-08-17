package service

import (
	"context"
	"errors"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type fakeUserRepo struct {
	err  error
	user *models.User
}

func (f *fakeUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return f.err
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		repoErr     error
		expectError bool
	}{
		{name: "success", repoErr: nil, expectError: false},
		{name: "repository error", repoErr: errors.New("email already exists"), expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeUserRepo{err: tt.repoErr}
			authService := NewAuthService(repo, "test-secret")
			user, err := authService.Register(context.Background(), "testuser", "test@test.com", "pass123")
			if tt.expectError == true && err == nil {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.expectError)
			}
			if tt.expectError == false && err != nil {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.expectError)
			}
			if !tt.expectError && user == nil {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.expectError)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	testUser := &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: string(hash),
	}

	tests := []struct {
		name        string
		repoErr     error
		repoUser    *models.User
		expectError bool
		password    string
	}{
		{name: "success", repoErr: nil, expectError: false, repoUser: testUser, password: "pass123"},
		{name: "user not found", repoErr: errors.New("not found"), repoUser: nil, password: "pass123", expectError: true},
		{name: "wrong password", repoErr: nil, repoUser: testUser, password: "wrong-password", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeUserRepo{err: tt.repoErr, user: tt.repoUser}
			authService := NewAuthService(repo, "test-secret")
			user, token, err := authService.Login(context.Background(), "testuser", tt.password)
			if tt.expectError == true && err == nil {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.expectError)
			}
			if tt.expectError == false && err != nil {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.expectError)
			}
			if !tt.expectError && user == nil {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.expectError)
			}
			if !tt.expectError && token == "" {
				t.Errorf("Login() token = %q, expected non-empty token", token)
			}
		})
	}
}

func (f *fakeUserRepo) GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error) {
	return f.user, f.err
}

func (f *fakeUserRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return nil, nil
}
