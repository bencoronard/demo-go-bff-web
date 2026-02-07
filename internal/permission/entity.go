package permission

type Permission struct {
	ID         uint
	Permission string
}

func (Permission) TableName() string {
	return "public.permissions"
}
