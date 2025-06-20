package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type jwtAuth struct {
	secret string
	iss    string
	sub    string
}

func NewJwtAuth(secret, iss, sub string) *jwtAuth {
	return &jwtAuth{
		secret: secret,
		iss:    iss,
		sub:    sub,
	}
}

func (j *jwtAuth) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *jwtAuth) VerifyToken(token string) (*jwt.Token, error) {

	// validate and verify the token
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(j.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithSubject(j.sub),
		jwt.WithIssuer(j.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

}
