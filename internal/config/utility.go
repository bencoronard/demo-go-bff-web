package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	xjwt "github.com/bencoronard/demo-go-common-libs/jwt"
)

func NewJwtIssuer(p *Properties) (xjwt.Issuer, error) {
	block, _ := pem.Decode([]byte(p.Secret.Crypto.PrivateKey))
	if block == nil {
		return nil, errors.New("failed to parse private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse private key")
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	issuer, err := xjwt.NewAsymmIssuer("bff", rsaKey)
	if err != nil {
		return nil, err
	}

	return issuer, nil
}
