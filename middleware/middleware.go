package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// JWTMiddleware function to check JWT token
func JWTMiddleware(jwtKey []byte) echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		SigningKey:    jwtKey,
		SigningMethod: "HS256",
		ContextKey:    "user",
		TokenLookup:   "header:Authorization",
		AuthScheme:    "Bearer",
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			var status int
			var message string

			// Default status and message
			status = http.StatusUnauthorized
			message = "Unauthorized"

			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					message = "Token has expired"
					logrus.WithFields(logrus.Fields{
						"error": "Token expired",
						"cause": err,
					}).Warn("JWT Middleware Error")
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					message = "Invalid token signature"
					logrus.WithFields(logrus.Fields{
						"error": "Invalid token signature",
						"cause": err,
					}).Warn("JWT Middleware Error")
				} else {
					message = "Invalid token"
					logrus.WithFields(logrus.Fields{
						"error": "Invalid token",
						"cause": err,
					}).Warn("JWT Middleware Error")
				}
			} else {
				// Log all other types of errors
				logrus.WithFields(logrus.Fields{
					"error": "Unauthorized access",
					"cause": err,
				}).Error("JWT Middleware Error")
			}

			return c.JSON(status, map[string]string{"error": message})
		},
	}
	return middleware.JWTWithConfig(config)
}

// RoleMiddleware function to check user role
func RoleMiddleware(requiredRole string, jwtKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token not found"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			role, ok := claims["role"].(string)
			if !ok || role != requiredRole {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have permission to access this route"})
			}

			return next(c)
		}
	}
}
