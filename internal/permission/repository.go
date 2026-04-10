package permission

import (
	"context"

	"gorm.io/gorm"
)

type PermissionRepo interface {
	ListAllPermissions(ctx context.Context) ([]Permission, error)
}

type permissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) PermissionRepo {
	return &permissionRepo{db: db}
}

func (p *permissionRepo) ListAllPermissions(ctx context.Context) ([]Permission, error) {
	var perms []Permission
	if err := p.db.Find(&perms).Error; err != nil {
		return []Permission{}, nil
	}
	return perms, nil
}
