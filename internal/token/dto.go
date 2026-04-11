package token

type issueTokenDTO struct {
	Id  string `validate:"required,numeric,gt=0"`
	Ttl string `validate:"required,numeric,gt=0"`
}
