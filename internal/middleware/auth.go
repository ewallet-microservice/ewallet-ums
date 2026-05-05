package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
)

type SessionRepository interface {
	GetUserSessionByRefreshToken(context.Context, string) error
}

type JWTManager interface {
	GenerateToken(user models.User, tokenType string) (string, error)
	ValidateToken(ctx context.Context, token string) (*helpers.ClaimToken, error)
}

type AuthMiddleware struct {
	sessionRepo SessionRepository
	jwtManager  JWTManager
}

func NewAuthMiddleware(sessionRepo SessionRepository, jwtManager JWTManager) *AuthMiddleware {
	return &AuthMiddleware{sessionRepo, jwtManager}
}

func (a *AuthMiddleware) MiddlewareRefreshToken(c *gin.Context) {
	var log = bootstrap.Log

	auth := c.Request.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {
		log.Infow("invalid or missing bearer token")
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	refreshToken := strings.TrimPrefix(auth, "Bearer ")
	if refreshToken == "" {
		log.Infow("invalid token")
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	err := a.sessionRepo.GetUserSessionByRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		log.Infow("failed to get user session on DB: ", err)
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	claim, err := a.jwtManager.ValidateToken(c.Request.Context(), refreshToken)
	if err != nil {
		log.Infow(err.Error())
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Infow("jwt token is expired: ", claim.ExpiresAt)
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	c.Set("refreshToken", refreshToken)
	c.Set("claim", claim)
	c.Next()
}

func (a *AuthMiddleware) MiddlewareAccessToken(c *gin.Context) {
	var log = bootstrap.Log

	auth := c.Request.Header.Get("Authorization")

	if !strings.HasPrefix(auth, "Bearer ") {
		log.Infow("invalid or missing bearer token")
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		log.Infow("invalid token")
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	claim, err := a.jwtManager.ValidateToken(c.Request.Context(), token)
	if err != nil {
		log.Infow(err.Error())
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Infow("jwt token is expired: ", claim.ExpiresAt)
		c.Error(constants.ErrorUnauthorized)
		c.Abort()
		return
	}

	c.Set("accessToken", token)
	c.Set("claim", claim)
	c.Next()

}
