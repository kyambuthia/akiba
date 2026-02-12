package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	issuer string
}

func NewJWTManager(secret, issuer string) *JWTManager {
	return &JWTManager{secret: []byte(secret), issuer: issuer}
}

func (j *JWTManager) IssueAccessToken(userID string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := Claims{Sub: userID, RegisteredClaims: jwt.RegisteredClaims{Issuer: j.issuer, Subject: userID, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(ttl))}}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(j.secret)
}

func (j *JWTManager) Verify(tokenString string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
