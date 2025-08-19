package services

import (
	"net/http"
	"time"

	"github.com/HanmaDevin/chatdev/dto"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

const (
	AccessTokenCookieName = "access_token"
	JwtSecretKey          = "some_secret_key"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateTokenAndSetCookies(user dto.UserDTO, c echo.Context) error {
	accessToken, exp, err := generateAccessToken(user.Username)
	if err != nil {
		return err
	}

	setTokenCookie(AccessTokenCookieName, accessToken, exp, c)
	return nil
}

func generateAccessToken(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(JwtSecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expirationTime, nil
}

func setTokenCookie(cookieName, token string, exp time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = token
	cookie.Expires = exp
	cookie.HttpOnly = true
	cookie.Secure = true // Set to true if using HTTPS
	cookie.Path = "/"
	c.SetCookie(cookie)
}
