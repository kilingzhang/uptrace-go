module github.com/uptrace/uptrace-go/example/go-pg

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-pg/pg/v10 v10.6.2
	github.com/go-pg/pgext v0.2.0
	github.com/uptrace/uptrace-go v0.3.0
	go.opentelemetry.io/otel v0.13.0
)
