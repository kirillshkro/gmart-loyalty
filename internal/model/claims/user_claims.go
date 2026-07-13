package claims

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kirillshkro/gmart-loyalty/internal/config"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID  int
	userCfg *config.AuthConfig
}

func NewUserClaims() *UserClaims {
	uc := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		userCfg: config.GetAuthConfig(),
	}
	return uc
}

func (uc UserClaims) Token() (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, uc).SignedString([]byte(uc.userCfg.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (uc UserClaims) UserIDByToken(token string) (int, error) {
	if token == "" {
		return 0, fmt.Errorf("token is empty")
	}
	if _, err := jwt.ParseWithClaims(token, &uc, func(t *jwt.Token) (any, error) {
		return []byte(uc.userCfg.SecretKey), nil
	}); err != nil {
		return 0, err
	}

	return uc.UserID, nil
}

func (u UserClaims) UserCfg() *config.AuthConfig {
	return u.userCfg
}
