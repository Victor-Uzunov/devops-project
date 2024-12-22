package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type ConfigOAuth2 struct {
	ClientID              string        `envconfig:"CLIENT_ID"`
	ClientSecret          string        `envconfig:"CLIENT_SECRET"`
	RedirectURL           string        `envconfig:"REDIRECT_URL"`
	Scopes                []string      `envconfig:"SCOPES"`
	OAuth2State           string        `envconfig:"OAUTH2_STATE"`
	JwtKey                string        `envconfig:"JWT_KEY"`
	JWTExpirationTime     time.Duration `envconfig:"JWT_EXPIRATION_TIME"`
	RefreshExpirationTime time.Duration `envconfig:"REFRESH_EXPIRATION_TIME"`
}

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type TokenParser struct {
	jwtKey string
}

func NewTokenParser(config ConfigOAuth2) *TokenParser {
	return &TokenParser{
		jwtKey: config.JwtKey,
	}
}

func (tp *TokenParser) ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tp.jwtKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, fmt.Errorf("invalid signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
