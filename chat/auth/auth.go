package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ezrod12/chat/models"
	"github.com/ezrod12/chat/services"
	"github.com/ezrod12/chat/settings"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string
	UserId   string
}

func AuthUser(response http.ResponseWriter, authRequest models.AuthenthicationRequest, context context.Context, collection *mongo.Collection) {
	response.Header().Set("Content-Type", "application/json")
	var dbUser models.User
	authRequest.Username = strings.ToLower(authRequest.Username)
	dbUser, err := services.GetUserByUsername(authRequest.Username, collection, context)

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte("user not found"))
		return
	}

	userPass := []byte(authRequest.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		log.Println(passErr)
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(`{"response":"Invalid Password!"}`))
		return
	}

	jwtToken, err := GenerateJWT(&Claims{Username: authRequest.Username, UserId: dbUser.Id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	services.UpdateLastConnection(dbUser, collection, context)
	response.Write([]byte(`{"token":"` + jwtToken + `", "username": "` + dbUser.Username + `" , "userId": "` + dbUser.Id + `" }`))
}

func GenerateJWT(c *Claims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = c.Username
	claims["userId"] = c.UserId
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(settings.Configuration.AppConfig.SecretKey))
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}

	return tokenString, nil
}

func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := settings.Configuration.AppConfig.SecretKey
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}
