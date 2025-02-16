package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type JWTDecoder struct {
	host     string
	endpoint string
	pubkey   string
}

// Assuming right now that endpoint for public key is always "token/publickey"
func NewJWTDecoder(host string, endpoint string, pubkey string) *JWTDecoder {
	return &JWTDecoder{
		host:     host,
		endpoint: endpoint,
		pubkey:   pubkey,
	}
}

func (v *JWTDecoder) getPublicKey() (string, error) {
	url := fmt.Sprintf("%s/%s", v.host, v.endpoint)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve public key: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve public key: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)
	}

	publicKey := string(body)
	if publicKey == "" {
		return "", errors.New("public key not found")
	}

	return publicKey, nil
}

func (v *JWTDecoder) DecodeToken(tokenString string) (jwt.MapClaims, error) {
	if v.pubkey == "" {
		pubkey, err := v.getPublicKey()
		if err != nil {
			return nil, err
		}
		v.pubkey = pubkey
	}

	// Decode PEM into *rsa.PublicKey
	block, _ := pem.Decode([]byte(v.pubkey))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	// Parse token with claims verification
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return rsaPubKey, nil
	}, jwt.WithAudience("ALL"))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	return claims, nil
}
