module github.com/bencoronard/demo-go-bff-web

go 1.26

replace github.com/bencoronard/demo-go-common-libs => ../demo-go-common-libs

require (
	github.com/bencoronard/demo-go-common-libs v0.0.0-20260408150834-2b7a80219142
	github.com/labstack/echo/v5 v5.1.0
	go.uber.org/fx v1.24.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
)
