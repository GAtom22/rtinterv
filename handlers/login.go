package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"retargetly-exercise/helpers"
	"retargetly-exercise/models"
	"time"
	"github.com/dgrijalva/jwt-go"
)

//LoginHandler handler for /login path - returns a token and its expiration date if login credentials are OK
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	switch r.Method {
		
	case http.MethodPost:
		response := models.APIResponse{}
		//Check user credentials
		ok, user, errMsg := isValidUser(r)
		if !ok {
			response.SendInfoMessage(w, errMsg, http.StatusBadRequest)
			return
		}
		//if OK generate and send auth token with expiration date
		// Generate token with user data and expiration time (in minutes)
		token, expirationDate, err := createToken(user, 15)
		if err != nil {
			response.SendInfoMessage(w, "Error while generating token", http.StatusInternalServerError)
			return
		}

		// send token data to client
		loginData := models.LoginResponseItem{
			Token:   token,
			Expires: expirationDate,
		}
		response.Content = loginData
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	default:
		// Response for other http Methods that are not POST
		response := models.APIResponse{}
		response.SendInfoMessage(w, "Not implemented, try with POST method", http.StatusNotImplemented)
	}
}

// isValidUser returns true if the user is valid, the user data as User struct and an error message describing the problem
func isValidUser(r *http.Request) (bool, models.User, string) {
	// Decode the request body data to the User struct
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Println("Error while decoding user credentials")
		return false, models.User{}, "Wrong format, please use JSON"
	}
	// Check user's credentials
	if u.UserName == "usuario" && u.Password == "contraseña" {
		return true, u, ""
	}
	return false, models.User{}, "Wrong credentials, user not found"
}

func createToken(u models.User, expirationTimeInMinutes float32) (string, string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	tokenExpirationTime := time.Now().Add(time.Minute * time.Duration(expirationTimeInMinutes)).Unix()
	// Set JWT token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = u.UserName
	claims["exp"] = tokenExpirationTime
	// Generate encoded token.
	t, err := token.SignedString([]byte(helpers.GetEnvVariable("TOKEN_SECRET_KEY")))
	if err != nil {
		return "", "", err
	}

	return t, helpers.FormatDate(tokenExpirationTime), nil
}
