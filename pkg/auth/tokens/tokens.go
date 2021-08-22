package tokens

import (
	"crypto/rsa"
	"io/ioutil"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/smartnuance/saas-kit/pkg/lib/tokens"
)

type TokenEnv struct {
	SigningKeyPath    string
	ValidationKeyPath string
	// Issuer is the issuer string of JWT tokens; defaults to service name
	Issuer string
}

type TokenController struct {
	TokenEnv
	signingKey    *rsa.PrivateKey
	ValidationKey *rsa.PublicKey
}

func Load(envs map[string]string, serviceName string) TokenEnv {
	issuer := envs["TOKEN_ISSUER"]
	if len(issuer) == 0 {
		issuer = serviceName
	}
	return TokenEnv{
		SigningKeyPath:    envs["TOKEN_SIGNING_KEY_PATH"],
		ValidationKeyPath: envs["TOKEN_VALIDATION_KEY_PATH"],
		Issuer:            issuer,
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
	c.ValidationKey, err = jwt.ParseRSAPublicKeyFromPEM(validationKey)
	if err != nil {
		return
	}
	return
}

func (c *TokenController) GenerateAccessToken(userID, instanceID int, isUser bool, role string) (token string, err error) {
	claims := tokens.AccessTokenClaims{
		User:     isUser,
		Role:     role,
		Instance: instanceID,
		StandardClaims: jwt.StandardClaims{
			Audience:  strconv.Itoa(instanceID),
			Subject:   strconv.Itoa(userID),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    c.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims)

	token, err = jwtToken.SignedString(c.signingKey)
	if err != nil {
		err = errors.Wrap(err, "signing access token failed")
		return
	}
	return
}

func (c *TokenController) GenerateRefreshToken(userID, instanceID int, isUser bool) (token string, expiresAt time.Time, err error) {
	expiresAt = time.Now().Add(time.Hour * 24 * 7)
	claims := jwt.StandardClaims{
		Audience:  strconv.Itoa(instanceID),
		Subject:   strconv.Itoa(userID),
		ExpiresAt: expiresAt.Unix(),
		Issuer:    c.Issuer,
		IssuedAt:  time.Now().Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token, err = jwtToken.SignedString(c.signingKey)
	if err != nil {
		err = errors.Wrap(err, "signing refresh token failed")
		return
	}
	return
}
