package token

import (
	"context"
	"strconv"
	"time"

	"github.com/bencoronard/demo-go-bff-web/internal/permission"
	"github.com/bencoronard/demo-go-common-libs/jwt"
)

type TokenService interface {
	IssueToken(ctx context.Context, id uint, ttl time.Duration) (string, error)
}

type tokenService struct {
	iss jwt.Issuer
	pr  permission.PermissionRepo
}

func NewTokenService(iss jwt.Issuer, pr permission.PermissionRepo) TokenService {
	return &tokenService{iss: iss, pr: pr}
}

func (t *tokenService) IssueToken(ctx context.Context, id uint, ttl time.Duration) (string, error) {
	perms, err := t.pr.ListAllPermissions(ctx)
	if err != nil {
		return "", err
	}

	claims := make(map[string]any, len(perms))
	for _, p := range perms {
		claims[p.Permission] = p.ID
	}

	return t.iss.IssueToken(strconv.FormatUint(uint64(id), 10), nil, claims, ttl, time.Time{})
}
