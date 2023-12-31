package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	ID       uint64 `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

func CheckJWT(c echo.Context) (uint64, string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	ckTkn, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return 0, "", echo.ErrUnauthorized
		}
		return 0, "", echo.ErrBadRequest
	}
	tknStr := ckTkn.Value
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return 0, "", echo.ErrUnauthorized
		}
		return 0, "", echo.ErrBadRequest
	}
	if !tkn.Valid {
		return 0, "", echo.ErrUnauthorized
	}
	return claims.ID, claims.Username, nil
}

func GenJWT(c echo.Context, id uint64, un string) error {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		ID:       id,
		Username: un,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return echo.ErrInternalServerError
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		SameSite: 1,
		Secure:   true,
		HttpOnly: true,
	})
	return nil
}

func RefJWT(c echo.Context) error {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	ckTkn, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return echo.ErrUnauthorized
		}
		return echo.ErrBadRequest
	}
	tknStr := ckTkn.Value
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return echo.ErrUnauthorized
		}
		return echo.ErrBadRequest
	}
	if !tkn.Valid {
		return echo.ErrUnauthorized
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return echo.ErrTooEarly
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return echo.ErrInternalServerError
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		SameSite: 1,
		Secure:   true,
		HttpOnly: true,
	})
	return nil
}

func Logout(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}
