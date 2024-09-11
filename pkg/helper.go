package pkg

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"strings"
	"time"
)

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

// ResponseJson is a utility function that sends a JSON response to the client
func ResponseJson(c echo.Context, status int, payload interface{}, message string) error {
	response := Response{
		Meta: Meta{
			Message: message,
			Code:    status,
			Status:  getStatusText(status),
		},
		Data: payload,
	}
	return c.JSON(status, response)
}

func getStatusText(status int) string {
	if status >= 200 && status < 300 {
		return "success"
	} else if status >= 400 && status < 500 {
		return "fail"
	} else {
		return "error"
	}
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash checks if a password matches its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func FormatValidationError(model interface{}, errors validator.ValidationErrors) string {
	var errorMessages []string
	modelType := reflect.TypeOf(model).Elem()

	for _, err := range errors {
		// Get the field name from the model's struct
		field, found := modelType.FieldByName(err.Field())
		fieldName := err.Field()
		if found {
			fieldName = field.Tag.Get("json")
			if fieldName == "" {
				fieldName = err.Field()
			}
		}

		// Generate custom error message based on the validation tag
		var message string
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", fieldName)
		case "numeric":
			message = fmt.Sprintf("%s must be a number", fieldName)
		case "unique":
			message = fmt.Sprintf("%s already exists", fieldName)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", fieldName)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", fieldName, err.Param())
		case "max":
			message = fmt.Sprintf("%s must not be longer than %s characters", fieldName, err.Param())
		case "gte":
			message = fmt.Sprintf("%s must be greater than or equal to %s", fieldName, err.Param())
		case "gt":
			message = fmt.Sprintf("%s must be greater than %s", fieldName, err.Param())
		case "lte":
			message = fmt.Sprintf("%s must be less than or equal to %s", fieldName, err.Param())
		case "lt":
			message = fmt.Sprintf("%s must be less than %s", fieldName, err.Param())
		case "eqfield":
			message = fmt.Sprintf("%s must be equal to %s", fieldName, err.Param())
		default:
			message = fmt.Sprintf("%s is invalid", fieldName)
		}

		errorMessages = append(errorMessages, message)
	}

	return strings.Join(errorMessages, ", ")
}

// GenerateToken generates a JWT token
func GenerateToken(userId string, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
