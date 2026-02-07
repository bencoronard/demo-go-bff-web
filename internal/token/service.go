package token

import (
	"context"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/bencoronard/demo-go-bff-web/internal/permission"
	"github.com/bencoronard/demo-go-common-libs/jwt"
)

type tokenService interface {
	issueToken(ctx context.Context) (string, error)
}

type tokenServiceImpl struct {
	iss jwt.Issuer
	pr  permission.PermissionRepo
}

func NewTokenService(iss jwt.Issuer, pr permission.PermissionRepo) tokenService {
	return &tokenServiceImpl{iss: iss, pr: pr}
}

func (t *tokenServiceImpl) issueToken(ctx context.Context) (string, error) {
	perms, err := t.pr.ListAllPermissions(ctx)
	if err != nil {
		return "", err
	}

	claims := make(map[string]any, len(perms))
	for _, p := range perms {
		claims[p.Permission] = p.ID
	}

	min, max := 10, 50
	id := rand.IntN(max-min) + min

	return t.iss.IssueToken(strconv.Itoa(id), nil, claims, 600*time.Second, time.Time{})
}
