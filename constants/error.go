package constants

import "net/http"

var (
	ErrBadRequest   = "bad request"
	ErrUnauthorized = "unauthorized"
)

var (
	ErrRequiredEmail         = "email is required"
	ErrRequiredUsername      = "username is required"
	ErrRequiredPassword      = "password is required"
	ErrRequiredFullName      = "full name is required"
	ErrDuplicateEmail        = "email is already registered"
	ErrUserNotFound          = "user not found"
	ErrSessionNotFound       = "user session not found"
	ErrFailedToGetToken      = "failed to get token"
	ErrFailedToParseToken    = "failed to parse token"
	ErrFailedToExtractClaims = "failed to extract claims"
	ErrFailedToParseClaims   = "failed to parse claims"
	ErrFailedToUpdateToken   = "failed to update token"
	ErrFailedToGenerateToken = "failed to generate token"
	ErrFailedToCreateSession = "failed to create session"
	ErrFailedToCreateWallet  = "failed to create wallet"
)

var ValidationErrorMap = map[string]map[string]string{
	"required": {
		"User.Username":         ErrRequiredUsername,
		"User.Email":            ErrRequiredEmail,
		"User.FullName":         ErrRequiredFullName,
		"User.Password":         ErrRequiredPassword,
		"LoginRequest.Email":    ErrRequiredEmail,
		"LoginRequest.Password": ErrRequiredPassword,
	},
}

type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statusCode int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Message: message}
}

var (
	ErrorBadRequest            = NewAppError(http.StatusBadRequest, ErrBadRequest)
	ErrorUnauthorized          = NewAppError(http.StatusUnauthorized, ErrUnauthorized)
	ErrorDuplicateEmail        = NewAppError(http.StatusConflict, ErrDuplicateEmail)
	ErrorRequiredEmail         = NewAppError(http.StatusBadRequest, ErrRequiredEmail)
	ErrorUserNotFound          = NewAppError(http.StatusNotFound, ErrUserNotFound)
	ErrorSessionNotFound       = NewAppError(http.StatusNotFound, ErrSessionNotFound)
	ErrorFailedToGetToken      = NewAppError(http.StatusInternalServerError, ErrFailedToGetToken)
	ErrorFailedToParseToken    = NewAppError(http.StatusInternalServerError, ErrFailedToParseToken)
	ErrorFailedToExtractClaims = NewAppError(http.StatusInternalServerError, ErrFailedToExtractClaims)
	ErrorFailedToParseClaims   = NewAppError(http.StatusInternalServerError, ErrFailedToParseClaims)
	ErrorFailedToUpdateToken   = NewAppError(http.StatusInternalServerError, ErrFailedToUpdateToken)
	ErrorFailedToGenerateToken = NewAppError(http.StatusInternalServerError, ErrFailedToGenerateToken)
	ErrorFailedToCreateSession = NewAppError(http.StatusInternalServerError, ErrFailedToCreateSession)
	ErrorFailedToCreateWallet  = NewAppError(http.StatusInternalServerError, ErrFailedToCreateWallet)
)

var (
	ErrorValidationUserEmail = NewAppError(http.StatusBadRequest, ErrRequiredEmail)
)
