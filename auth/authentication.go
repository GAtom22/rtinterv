package auth

import(
	"net/http"
	"strings"
	"fmt"
	"os"
	"github.com/dgrijalva/jwt-go"
)

func ExtractToken(r *http.Request) string {
  bearToken := r.Header.Get("Authorization")
  //expecting "Authorization": {token}
  strArr := strings.Split(bearToken, " ")
  if len(strArr) == 2 {
     return strArr[1]
  }
  return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
  tokenString := ExtractToken(r)
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
     //Make sure that the token method conform to "SigningMethodHMAC"
     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
     }
     return []byte(os.Getenv("TOKEN_SECRET_KEY")), nil
  })
  if err != nil {
     return nil, err
  }
  return token, nil
}

func TokenValid(r *http.Request) error {
  token, err := VerifyToken(r)
  if err != nil {
     return err
  }
  if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
     return err
  }
  return nil
}
