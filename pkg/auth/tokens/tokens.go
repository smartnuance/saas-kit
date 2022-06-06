package tokens

import (
	"crypto/rsa"
	"io/ioutil"
	"time"

	"github.com/friendsofgo/errors"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/smartnuance/saas-kit/pkg/lib"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
	"github.com/smartnuance/saas-kit/pkg/lib/tokens"
)

type TokenEnv struct {
	SigningKeyPath    string
	ValidationKeyPath string
	// Issuer is the issuer string of JWT tokens; defaults to service name
	Issuer string
	// Audience is the audience string of JWT tokens; defaults to DefaultAudience
	Audience string
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
	audience := envs["TOKEN_AUDIENCE"]
	if len(audience) == 0 {
		audience = lib.DefaultAudience
	}
	return TokenEnv{
		SigningKeyPath:    envs["TOKEN_SIGNING_KEY_PATH"],
		ValidationKeyPath: envs["TOKEN_VALIDATION_KEY_PATH"],
		Issuer:            issuer,
		Audience:          audience,
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

func (c *TokenController) GenerateAccessToken(userID, instanceID string, role roles.Role) (token string, err error) {
	claims := tokens.AccessTokenClaims{
		Purpose:  tokens.AccessPurpose,
		Role:     string(role),
		Instance: instanceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			Issuer:    c.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  []string{c.Audience},
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

func (c *TokenController) GenerateRefreshToken(userID, instanceID string) (token string, expiresAt time.Time, err error) {
	expiresAt = time.Now().Add(time.Hour * 24 * 7)
	claims := tokens.RefreshTokenClaims{
		Purpose:  tokens.RefreshPurpose,
		Instance: instanceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    c.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  []string{c.Audience},
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token, err = jwtToken.SignedString(c.signingKey)
	if err != nil {
		err = errors.Wrap(err, "signing refresh token failed")
		return
	}
	return
}
