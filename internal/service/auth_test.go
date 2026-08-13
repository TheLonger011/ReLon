package service

//
//import (
//	"context"
//	"errors"
//	"github.com/TheLonger011/ReLon/internal/models"
//	"testing"
//)
//
//type fakeUserRepo struct {
//	err error
//}
//
//func (f *fakeUserRepo) CreateUser(ctx context.Context, user *models.User) error {
//	return f.err
//}
//
//func TestRegister(t *testing.T) {
//	tests := []struct {
//		name        string
//		repoErr     error
//		expectError bool
//	}{
//		{name: "success", repoErr: nil, expectError: false},
//		{name: "repository error", repoErr: errors.New("email already exists"), expectError: true},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			repo := &fakeUserRepo{err: tt.repoErr}
//			authService := NewAuthService(repo)
//			user, err := authService.Register(context.Background(), "testuser", "test@test.com", "pass123")
//			if tt.expectError == true && err == nil {
//				t.Errorf("Register() error = %v, wantErr %v", err, tt.expectError)
//			}
//			if tt.expectError == false && err != nil {
//				t.Errorf("Register() error = %v, wantErr %v", err, tt.expectError)
//			}
//			if !tt.expectError && user == nil {
//				t.Errorf("Register() error = %v, wantErr %v", err, tt.expectError)
//			}
//		})
//	}
//
//}
