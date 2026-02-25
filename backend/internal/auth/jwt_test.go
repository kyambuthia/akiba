package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTVerifyRejectsWrongIssuer(t *testing.T) {
	goodMgr := NewJWTManager("secret", "akiba-api")
	badIssuerMgr := NewJWTManager("secret", "other-issuer")

	token, err := badIssuerMgr.IssueAccessToken("u1", time.Hour)
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	if _, err := goodMgr.Verify(token); err == nil {
		t.Fatalf("expected issuer verification error")
	}
}

func TestJWTVerifyRejectsWrongAlgorithm(t *testing.T) {
	mgr := NewJWTManager("secret", "akiba-api")
	claims := Claims{
		Sub: "u1",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "akiba-api",
			Subject:   "u1",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	if _, err := mgr.Verify(tokenString); err == nil {
		t.Fatalf("expected algorithm verification error")
	}
}
