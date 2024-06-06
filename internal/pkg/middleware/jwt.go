// pkg/middleware/jwt.go
package middleware

import (
    "time"
    "net/http"
    "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func RefreshJWT(r *http.Request) (string, error) {
    tokenStr := r.Header.Get("Authorization")

    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        return "", err
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims.ExpiresAt = expirationTime.Unix()
    newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return newToken.SignedString(jwtKey)
}
