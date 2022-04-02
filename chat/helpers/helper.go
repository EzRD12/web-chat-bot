package helpers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ezrod12/chat/settings"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func IsAuthorized(request *http.Request, response http.ResponseWriter) bool {
	if request.Header["Authorization"] != nil {
		token, err := jwt.Parse(request.Header["Authorization"][0], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				msg := "Error in JWT token validation"
				return nil, fmt.Errorf("%s", msg)
			}
			return []byte(settings.Configuration.AppConfig.SecretKey), nil
		})

		if err != nil {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte(err.Error()))
			return false
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return true
		} else {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("Not Authorized"))
			return false
		}
	} else {
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte("Not Authorized"))
		return false
	}
}
