package tokens

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type TokenEnv struct {
	SigningKeyPath    string
	ValidationKeyPath string
	Issuer            string
}

type TokenController struct {
	TokenEnv
	signingKey    *rsa.PrivateKey
	validationKey *rsa.PublicKey
}

func Load(envs map[string]string) TokenEnv {
	return TokenEnv{
		SigningKeyPath:    envs["TOKEN_SIGNING_KEY_PATH"],
		ValidationKeyPath: envs["TOKEN_VALIDATION_KEY_PATH"],
	}
}

func Setup(env TokenEnv) (c *TokenController, err error) {
	c = &TokenController{TokenEnv: env}

	signingKey, err := ioutil.ReadFile(env.SigningKeyPath)
	if err != nil {
		err = errors.Wrapf(err, "could not read signing key file at "+env.SigningKeyPath)
		return
	}
	c.signingKey, err = jwt.ParseRSAPrivateKeyFromPEM(signingKey)
	if err != nil {
		return
	}

	validationKey, err := ioutil.ReadFile(env.ValidationKeyPath)
	if err != nil {
		err = errors.Wrapf(err, "could not read validation key file at "+env.ValidationKeyPath)
		return
	}
	c.validationKey, err = jwt.ParseRSAPublicKeyFromPEM(validationKey)
	if err != nil {
		return
	}
	return
}

type authCustomClaims struct {
	Name string `json:"name"`
	User bool   `json:"user"`
	jwt.StandardClaims
}

func (c *TokenController) GenerateToken(email string, isUser bool) (token string, err error) {
	claims := &authCustomClaims{
		email,
		isUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    c.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token, err = jwtToken.SignedString(c.SigningKeyPath)
	if err != nil {
		return
	}
	return
}

func (c *TokenController) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte(c.SigningKeyPath), nil
	})

}
