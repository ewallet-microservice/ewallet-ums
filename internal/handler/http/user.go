package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
	"github.com/mhasnanr/ewallet-ums/internal/models/dto"
)

type UserService interface {
	Register(context.Context, *models.User) (*models.User, error)
	Login(context.Context, models.LoginRequest) (dto.LoginResponse, error)
	UpdateTokenByRefreshToken(context.Context, string, *helpers.ClaimToken) (string, error)
	Logout(context.Context, string) error
}

type AuthMiddleware interface {
	MiddlewareRefreshToken(c *gin.Context)
	MiddlewareAccessToken(c *gin.Context)
}

type UserHandler struct {
	service        UserService
	authMiddleware AuthMiddleware
}

func NewUserHandler(svc UserService, authMiddleware AuthMiddleware) *UserHandler {
	return &UserHandler{service: svc, authMiddleware: authMiddleware}
}

func (h *UserHandler) RegisterRoute(c *gin.Engine) {
	userV1 := c.Group("/users/v1")
	userV1.POST("/register", h.registerUser)
	userV1.POST("/login", h.login)
	userV1.DELETE("/logout", h.authMiddleware.MiddlewareAccessToken, h.logout)
	userV1.GET("/token/refresh", h.authMiddleware.MiddlewareRefreshToken, h.refreshToken)
	userV1.GET("/token/validate", h.authMiddleware.MiddlewareAccessToken, h.validateToken)

}

func (h *UserHandler) registerUser(c *gin.Context) {
	var req models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(constants.ErrorBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		c.Error(err)
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusCreated, constants.UserCreated, user)
}

func (h *UserHandler) login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.Error(err)
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.LoginSucceed, res)
}

func (h *UserHandler) logout(c *gin.Context) {
	token, ok := c.Get("accessToken")
	if !ok {
		c.Error(constants.ErrorFailedToGetToken)
		return
	}

	accessToken, ok := token.(string)
	if !ok {
		c.Error(constants.ErrorFailedToParseToken)
		return
	}

	err := h.service.Logout(c.Request.Context(), accessToken)
	if err != nil {
		c.Error(err)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusNoContent, constants.LoginSucceed, nil)
}

func (h *UserHandler) refreshToken(c *gin.Context) {
	var response dto.RefreshTokenResponse

	token, ok := c.Get("refreshToken")
	if !ok {
		c.Error(constants.ErrorFailedToGetToken)
		return
	}

	refreshToken, ok := token.(string)
	if !ok {
		c.Error(constants.ErrorFailedToParseToken)
		return
	}

	val, ok := c.Get("claim")
	if !ok {
		c.Error(constants.ErrorFailedToExtractClaims)
		return
	}

	claim, ok := val.(*helpers.ClaimToken)
	if !ok {
		c.Error(constants.ErrorFailedToParseClaims)
		return
	}

	newToken, err := h.service.UpdateTokenByRefreshToken(c.Request.Context(), refreshToken, claim)
	if err != nil {
		c.Error(constants.ErrorFailedToUpdateToken)
		return
	}

	response.Token = newToken

	helpers.SendResponseHTTP(c, http.StatusOK, constants.NewToken, response)
}

func (h *UserHandler) validateToken(c *gin.Context) {
	var response dto.ValidateTokenResponse

	token, ok := c.Get("accessToken")
	if !ok {
		c.Error(constants.ErrorFailedToGetToken)
		return
	}

	_, ok = token.(string)
	if !ok {
		c.Error(constants.ErrorFailedToParseToken)
		return
	}

	val, ok := c.Get("claim")
	if !ok {
		c.Error(constants.ErrorFailedToExtractClaims)
		return
	}

	claim, ok := val.(*helpers.ClaimToken)
	if !ok {
		c.Error(constants.ErrorFailedToParseClaims)
		return
	}

	response.UserID = claim.UserID
	response.Email = claim.Email
	response.Username = claim.Username
	response.Fullname = claim.Username

	helpers.SendResponseHTTP(c, http.StatusOK, constants.ValidToken, response)
}
