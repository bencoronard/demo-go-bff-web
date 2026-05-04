package token

type issueTokenDTO struct {
	Id  string `validate:"required,numeric,gt=0" msg:"ID must be a positive number"`
	Ttl string `validate:"required,numeric,gt=0" msg:"TTL must be a positive number"`
}
