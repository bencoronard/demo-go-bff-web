package config

import (
	xjwt "github.com/bencoronard/demo-go-common-libs/jwt"
)

func NewJwtIssuer() (xjwt.Issuer, error) {
	iss, err := xjwt.NewAsymmIssuer("BFF", nil)
	if err != nil {
		return nil, err
	}
	return iss, nil
}
