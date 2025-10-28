go mod init backend
go mod tidy

go install github.com/swaggo/swag/cmd/swag@latest
swag init


go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate create -ext sql -dir ./migrations -seq init_table