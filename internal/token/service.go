package token

import (
	"context"
	"time"

	"github.com/bencoronard/demo-go-common-libs/jwt"
)

type tokenService interface {
	issueToken(ctx context.Context) (string, error)
}

type tokenServiceImpl struct {
	iss jwt.Issuer
}

func NewTokenService(iss jwt.Issuer) tokenService {
	return &tokenServiceImpl{iss: iss}
}

func (t *tokenServiceImpl) issueToken(ctx context.Context) (string, error) {
	return t.iss.IssueToken("1", nil, nil, 300*time.Second, time.Time{})
}
