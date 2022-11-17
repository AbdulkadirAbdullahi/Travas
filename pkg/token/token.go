package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
	"time"
)

type TravasClaims struct {
	jwt.RegisteredClaims
	Email string
	Name  string
}

var secretKey = os.Getenv("TRAVAS_KEY")

func GenerateToken(email, name string) (string, string, error) {
	travasClaims := TravasClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "travasAdmin",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(24 * time.Hour),
			},
		},
		Email: email,
		Name:  name,
	}
	refTravasClaims := &jwt.RegisteredClaims{
		Issuer:   "travasAdmin",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(48 * time.Hour),
		},
	}
	travasToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, travasClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	refTravasToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refTravasClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	return travasToken, refTravasToken, nil
}

func ParseTokenString(tokenString string) (*TravasClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TravasClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Fatalf("error while parsing token with it claims %v", err)
	}
	claims, ok := token.Claims.(*TravasClaims)
	if !ok {
		log.Fatalf("error %v controller not authorized access", http.StatusUnauthorized)
	}
	if err := claims.Valid(); err != nil {
		log.Fatalf("error %v %s", http.StatusUnauthorized, err)
	}
	return claims, nil
}