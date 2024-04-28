package middleware

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ValidateBasicAuth(expectedUsername, expectedPassword string) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(expectedUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) == 1 {
			return true, nil
		}
		return false, nil
	})
}
