package services

import (
	"context"
	"time"

	"github.com/mhasnanr/ewallet-ums/cmd/wallet"
	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
	"github.com/mhasnanr/ewallet-ums/internal/models/dto"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashed string, plain string) error
}

type JWTManager interface {
	GenerateToken(user models.User, tokenType string) (string, error)
	ValidateToken(ctx context.Context, token string) (*helpers.ClaimToken, error)
}

type UserRepository interface {
	Register(context.Context, *models.User) (*models.User, error)
	GetUserByEmail(context.Context, string) (models.User, error)
	DeleteUser(context.Context, int) error
}

type WalletAPI interface {
	CreateWallet(ctx context.Context, userID int) (*wallet.CreateWalletResponse, error)
}

type SessionRepository interface {
	CreateUserSession(context.Context, models.UserSession) error
	GetUserSessionByRefreshToken(context.Context, string) error
	UpdateTokenByRefreshToken(context.Context, string, string) error
	DeleteUserSessionByToken(context.Context, string) error
}

type UserService struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
	jwtManager  JWTManager
	pwHasher    PasswordHasher
	walletAPI   WalletAPI
}

func NewUserService(userRepo UserRepository, sessionRepo SessionRepository, jwtManager JWTManager, pwHasher PasswordHasher, walletAPI WalletAPI) *UserService {
	return &UserService{userRepo, sessionRepo, jwtManager, pwHasher, walletAPI}
}

func (s *UserService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	returnedUser, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if returnedUser.ID != 0 {
		return nil, constants.ErrorDuplicateEmail
	}

	hashedPassword, err := s.pwHasher.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	user, err = s.userRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	if _, err := s.walletAPI.CreateWallet(ctx, user.ID); err != nil {
		if err := s.userRepo.DeleteUser(ctx, user.ID); err != nil {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (dto.LoginResponse, error) {
	var (
		response dto.LoginResponse
		now      = time.Now()
	)

	returnedUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if returnedUser.ID == 0 {
		return response, constants.ErrorUserNotFound
	}

	if err != nil {
		return response, err
	}

	err = s.pwHasher.VerifyPassword(returnedUser.Password, req.Password)
	if err != nil {
		return response, err
	}

	token, err := s.jwtManager.GenerateToken(returnedUser, "token")
	if err != nil {
		return response, constants.ErrorFailedToGenerateToken
	}

	refreshToken, err := s.jwtManager.GenerateToken(returnedUser, "refreshToken")
	if err != nil {
		return response, constants.ErrorFailedToGenerateToken
	}

	response.Token = token

	userSession := models.UserSession{
		UserID:              uint(returnedUser.ID),
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        now.Add(helpers.MapTypeToken["token"]),
		RefreshTokenExpired: now.Add(helpers.MapTypeToken["refresh_token"]),
	}

	err = s.sessionRepo.CreateUserSession(ctx, userSession)
	if err != nil {
		return response, constants.ErrorFailedToCreateSession
	}

	response.Token = token
	response.RefreshToken = refreshToken

	return response, nil
}

func (s *UserService) UpdateTokenByRefreshToken(ctx context.Context, refreshToken string, claims *helpers.ClaimToken) (string, error) {
	var user models.User

	user.ID = claims.UserID
	user.Username = claims.Username
	user.Email = claims.Email
	user.FullName = claims.Fullname

	newToken, err := s.jwtManager.GenerateToken(user, "token")
	if err != nil {
		return newToken, constants.ErrorFailedToGenerateToken
	}

	err = s.sessionRepo.UpdateTokenByRefreshToken(ctx, newToken, refreshToken)
	if err != nil {
		return newToken, constants.ErrorFailedToUpdateToken
	}

	return newToken, nil
}

func (s *UserService) Logout(ctx context.Context, accessToken string) error {
	return s.sessionRepo.DeleteUserSessionByToken(ctx, accessToken)
}
