module github.com/suprsokr/vanilladbc-mysql

go 1.25.5

replace github.com/suprsokr/vanilladbc => ../vanilladbc

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/suprsokr/vanilladbc v0.0.0-00010101000000-000000000000
)

require filippo.io/edwards25519 v1.1.0 // indirect
