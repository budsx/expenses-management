package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/budsx/expenses-management/entity"
	"github.com/budsx/expenses-management/model"
	"github.com/budsx/expenses-management/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_AuthenticateUser(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name     string
		email    string
		password string
		mock     func(server *TestService)
		want     *model.LoginResponse
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "success - employee login",
			email:    "employee@example.com",
			password: "password123",
			mock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "employee@example.com").
					Return(&entity.User{
						ID:           1,
						Email:        "employee@example.com",
						Name:         "Test Employee",
						Role:         int(util.USER_ROLE_EMPLOYEE),
						PasswordHash: string(hashedPassword),
						CreatedAt:    time.Now(),
					}, nil).
					Times(1)
			},
			want: &model.LoginResponse{
				Token:     "mocked_token", // Will be mocked
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
			wantErr: false,
		},
		{
			name:     "success - manager login",
			email:    "manager@example.com",
			password: "password123",
			mock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "manager@example.com").
					Return(&entity.User{
						ID:           2,
						Email:        "manager@example.com",
						Name:         "Test Manager",
						Role:         int(util.USER_ROLE_MANAGER),
						PasswordHash: string(hashedPassword),
						CreatedAt:    time.Now(),
					}, nil).
					Times(1)
			},
			want: &model.LoginResponse{
				Token:     "mocked_token",
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
			wantErr: false,
		},
		{
			name:     "failure - user not found",
			email:    "nonexistent@example.com",
			password: "password123",
			mock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "nonexistent@example.com").
					Return(nil, errors.New("user not found")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name:     "failure - wrong password",
			email:    "employee@example.com",
			password: "wrongpassword",
			mock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "employee@example.com").
					Return(&entity.User{
						ID:           1,
						Email:        "employee@example.com",
						Name:         "Test Employee",
						Role:         int(util.USER_ROLE_EMPLOYEE),
						PasswordHash: string(hashedPassword),
						CreatedAt:    time.Now(),
					}, nil).
					Times(1)
			},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name:     "failure - empty email",
			email:    "",
			password: "password123",
			mock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "").
					Return(nil, errors.New("invalid email")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name:     "failure - empty password",
			email:    "employee@example.com",
			password: "",
			mock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "employee@example.com").
					Return(&entity.User{
						ID:           1,
						Email:        "employee@example.com",
						Name:         "Test Employee",
						Role:         int(util.USER_ROLE_EMPLOYEE),
						PasswordHash: string(hashedPassword),
						CreatedAt:    time.Now(),
					}, nil).
					Times(1)
			},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name:     "failure - database error",
			email:    "employee@example.com",
			password: "password123",
			mock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "employee@example.com").
					Return(nil, errors.New("database connection failed")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
			errMsg:  "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServerWithUserRepo(t)
			defer server.MockCtrl.Finish()

			ctx := context.Background()
			tt.mock(server)

			got, err := server.Service.AuthenticateUser(ctx, tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.NotEmpty(t, got.Token)
			assert.Greater(t, got.ExpiresAt, time.Now().Unix())
		})
	}
}

func TestPasswordHashing(t *testing.T) {
	tests := []struct {
		name           string
		originalPwd    string
		testPwd        string
		expectedResult bool
	}{
		{
			name:           "correct password",
			originalPwd:    "testpassword123",
			testPwd:        "testpassword123",
			expectedResult: true,
		},
		{
			name:           "wrong password",
			originalPwd:    "testpassword123",
			testPwd:        "wrongpassword",
			expectedResult: false,
		},
		{
			name:           "empty password",
			originalPwd:    "testpassword123",
			testPwd:        "",
			expectedResult: false,
		},
		{
			name:           "case sensitive",
			originalPwd:    "TestPassword123",
			testPwd:        "testpassword123",
			expectedResult: false,
		},
		{
			name:           "special characters",
			originalPwd:    "P@ssw0rd!123",
			testPwd:        "P@ssw0rd!123",
			expectedResult: true,
		},
		{
			name:           "unicode characters",
			originalPwd:    "пароль123",
			testPwd:        "пароль123",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tt.originalPwd), bcrypt.DefaultCost)
			assert.NoError(t, err)
			assert.NotEmpty(t, hashedPassword)

			result := util.CheckPasswordHash(tt.testPwd, string(hashedPassword))
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestAuthService_UserRolesAndJWT(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name         string
		userID       int64
		email        string
		role         int
		roleName     string
		expectToken  bool
		expectExpiry bool
	}{
		{
			name:         "admin user authentication",
			userID:       1,
			email:        "admin@example.com",
			role:         int(util.USER_ROLE_ADMIN),
			roleName:     "Admin",
			expectToken:  true,
			expectExpiry: true,
		},
		{
			name:         "manager user authentication",
			userID:       2,
			email:        "manager@example.com",
			role:         int(util.USER_ROLE_MANAGER),
			roleName:     "Manager",
			expectToken:  true,
			expectExpiry: true,
		},
		{
			name:         "employee user authentication",
			userID:       3,
			email:        "employee@example.com",
			role:         int(util.USER_ROLE_EMPLOYEE),
			roleName:     "Employee",
			expectToken:  true,
			expectExpiry: true,
		},
		{
			name:         "user with invalid role",
			userID:       4,
			email:        "invalid@example.com",
			role:         999, // Invalid role
			roleName:     "Invalid",
			expectToken:  true, // JWT should still be generated
			expectExpiry: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServerWithUserRepo(t)
			defer server.MockCtrl.Finish()

			server.MockUserRepo.EXPECT().
				GetUserWithPassword(gomock.Any(), tt.email).
				Return(&entity.User{
					ID:           tt.userID,
					Email:        tt.email,
					Name:         "Test " + tt.roleName,
					Role:         tt.role,
					PasswordHash: string(hashedPassword),
					CreatedAt:    time.Now(),
				}, nil).
				Times(1)

			ctx := context.Background()
			got, err := server.Service.AuthenticateUser(ctx, tt.email, "password123")

			assert.NoError(t, err)
			assert.NotNil(t, got)

			if tt.expectToken {
				assert.NotEmpty(t, got.Token)
				assert.True(t, len(got.Token) > 20)
			}

			if tt.expectExpiry {
				assert.Greater(t, got.ExpiresAt, time.Now().Unix())
				assert.Less(t, got.ExpiresAt, time.Now().Add(48*time.Hour).Unix())
			}
		})
	}
}

func TestAuthService_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		setupMock   func(server *TestService)
		expectError bool
		expectedMsg string
	}{
		{
			name:     "null context",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "test@example.com").
					Return(nil, errors.New("context deadline exceeded")).
					Times(1)
			},
			expectError: true,
			expectedMsg: "invalid credentials",
		},
		{
			name:     "extremely long email",
			email:    "very_long_email_address_that_exceeds_normal_limits_" + string(make([]byte, 200)) + "@example.com",
			password: "password123",
			setupMock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("invalid email format")).
					Times(1)
			},
			expectError: true,
			expectedMsg: "invalid credentials",
		},
		{
			name:     "sql injection attempt in email",
			email:    "admin@example.com'; DROP TABLE users; --",
			password: "password123",
			setupMock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "admin@example.com'; DROP TABLE users; --").
					Return(nil, errors.New("invalid email")).
					Times(1)
			},
			expectError: true,
			expectedMsg: "invalid credentials",
		},
		{
			name:     "user with corrupted password hash",
			email:    "corrupted@example.com",
			password: "password123",
			setupMock: func(server *TestService) {
				server.MockUserRepo.EXPECT().
					GetUserWithPassword(gomock.Any(), "corrupted@example.com").
					Return(&entity.User{
						ID:           1,
						Email:        "corrupted@example.com",
						Name:         "Corrupted User",
						Role:         int(util.USER_ROLE_EMPLOYEE),
						PasswordHash: "corrupted_hash_not_bcrypt",
						CreatedAt:    time.Now(),
					}, nil).
					Times(1)
			},
			expectError: true,
			expectedMsg: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServerWithUserRepo(t)
			defer server.MockCtrl.Finish()

			tt.setupMock(server)

			ctx := context.Background()
			got, err := server.Service.AuthenticateUser(ctx, tt.email, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.expectedMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
