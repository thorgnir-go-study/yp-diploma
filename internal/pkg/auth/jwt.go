package auth

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtWrapper struct {
	secretKey       string
	Issuer          string
	ExpirationHours int64
}

func NewJwtWrapper(secretKey string, issuer string, expirationHours int64) *JwtWrapper {
	return &JwtWrapper{secretKey: secretKey, Issuer: issuer, ExpirationHours: expirationHours}
}

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

var (
	ErrClaimsParseError = errors.New("couldn't parse claims")
	ErrTokenExpired     = errors.New("token expired")
)

func (j *JwtWrapper) GenerateToken(userID uuid.UUID) (string, error) {
	claims := &CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *JwtWrapper) ValidateToken(signedToken string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrClaimsParseError
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, ErrTokenExpired
	}

	return claims, nil
}
