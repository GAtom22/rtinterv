package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"retargetly-exercise/structs"
	"time"
	"os"
	"github.com/dgrijalva/jwt-go"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		response := structs.APIResponse{}
		//Check user credentials
		ok, u := isValidUser(r)
		if !ok {
			response.SendInfoMessage(w, "Error: Wrong credentials, user not found or wrong format, please use JSON", http.StatusUnauthorized)
			return
		}
		//if OK generate and send auth token with expiration date
		token, expirationDate, err := createToken(u, 10)
		if err != nil {
			response.SendInfoMessage(w, "Error while generating token", http.StatusInternalServerError)
			return
		}

		loginData := structs.LoginResponseItem{
			Token:   token,
			Expires: expirationDate,
		}
		response.Content = loginData
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	default:
		// Response for other http Methods that are not POST
		response := structs.APIResponse{}
		response.SendInfoMessage(w, "Not implemented, try with POST method", http.StatusNotImplemented)
	}
}

func isValidUser(r *http.Request) (bool, structs.User) {
	// Decode the request body data to the User struct
	var u structs.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Println("Error while decoding user credentials")
		return false, structs.User{}
	}
	// Check user's credentials
	if u.UserId == "usuario" && u.Password == "contraseña" {
		return true, u
	}
	return false, structs.User{}
}

func createToken(u structs.User, expirationTimeInMinutes int) (string, string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	tokenExpirationTime := time.Now().Add(time.Minute * time.Duration(expirationTimeInMinutes))
	// Set JWT token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = u.UserId
	claims["exp"] = tokenExpirationTime
	// Generate encoded token.
	t, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET_KEY")))
	if err != nil {
		return "", "", err
	}

	return t, formatExpirationTime(tokenExpirationTime), nil
}

func formatExpirationTime(t time.Time) string {
	expTimeString := t.Format(time.RFC3339)
	// Remove the last 6 characters (-03:00)
	expTimeString = expTimeString[:len(expTimeString)-6]

	return expTimeString
}
