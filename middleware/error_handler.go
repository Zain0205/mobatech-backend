package middleware

import (
	"net/http"

	"backend/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if appErr, ok := err.(*utils.AppError); ok {
				c.JSON(appErr.Status, utils.BuildError(appErr.Code, appErr.Message, nil))
				return
			}

			// Default internal error
			c.JSON(http.StatusInternalServerError, utils.BuildError(utils.ErrInternal, err.Error(), nil))
		}
	}
}
