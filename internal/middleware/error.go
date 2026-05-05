package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()

		var appErr *constants.AppError
		var valErrs validator.ValidationErrors

		if errors.As(err, &appErr) {
			helpers.SendResponseHTTP(c, appErr.StatusCode, appErr.Message, nil)
			return
		}

		if errors.As(err, &valErrs) {
			errStr := helpers.ConstructErrString(valErrs)
			helpers.SendResponseHTTP(c, http.StatusBadRequest, errStr, nil)
			return
		}

		helpers.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
	}
}
