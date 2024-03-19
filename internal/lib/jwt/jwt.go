package jwt_encoder

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"url-shortener-api/internal/lib/user"
)

type JWTEncoder struct {
	key string
}

func New(key string) (*JWTEncoder, error) {
	if key == "" {
		return nil, fmt.Errorf("key is empty")
	}
	return &JWTEncoder{key: key}, nil
}

func (encoder *JWTEncoder) CreateToken(user *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString([]byte(encoder.key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (encoder *JWTEncoder) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(encoder.key), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
