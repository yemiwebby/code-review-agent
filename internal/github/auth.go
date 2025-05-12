package github

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppAuthenticator struct {
	AppID      string
	PrivateKey *rsa.PrivateKey
}

func NewAppAuthenticator(appID, privateKeyPath string) (*AppAuthenticator, error) {
	// privateKeyData, err := ioutil.ReadFile(privateKeyPath)
	privateKeyData, err := os.ReadFile(privateKeyPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &AppAuthenticator{
		AppID:      appID,
		PrivateKey: privateKey,
	}, nil
}

func (a *AppAuthenticator) GenerateJWT() (string, error) {
	now := time.Now().UTC()
	exp := now.Add(9 * time.Minute)

	claims := jwt.MapClaims{
		"iss": a.AppID,
		"iat": now.Unix(),
		"exp": exp.Unix(),
	}

	fmt.Printf("Current UTC Time: %s\n", now)
	fmt.Printf("Generated JWT with iat: %d (%s), exp: %d (%s)\n",
		claims["iat"], time.Unix(claims["iat"].(int64), 0).UTC(),
		claims["exp"], time.Unix(claims["exp"].(int64), 0).UTC())

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(a.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	fmt.Printf("Generated JWT with iat: %d, exp: %d\n", claims["iat"], claims["exp"])

	return signedToken, nil
}
