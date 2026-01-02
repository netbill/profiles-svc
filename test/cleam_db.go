package test

import (
	"testing"

	"github.com/netbill/profiles-svc/cmd/migrations"
)

const TestDatabaseURL = "postgresql://postgres:postgres@localhost:7777/postgres?sslmode=disable"

func CleanDb(t *testing.T) {
	err := migrations.MigrateDown(TestDatabaseURL)
	if err != nil {
		t.Fatalf("migrate down: %v", err)
	}
	err = migrations.MigrateUp(TestDatabaseURL)
	if err != nil {
		t.Fatalf("migrate up: %v", err)
	}
}
