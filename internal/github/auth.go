package github

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yemiwebby/code-review-agent/config"
)

type AppAuthenticator struct {
	AppID      string
	PrivateKey *rsa.PrivateKey
}

func NewAppAuthenticator(appID string) (*AppAuthenticator, error) {
	// privateKeyData, err := ioutil.ReadFile(privateKeyPath)
	privateKeyBase64 := config.GithubPrivateKey
	if privateKeyBase64 == "" {
		return nil, fmt.Errorf("GITHUB_PRIVATE_KEY environment variable not set")
	}

	privateKeyData, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(a.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
}
