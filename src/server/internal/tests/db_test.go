package tests

import (
	"deploy-buddy/server/internal/config"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestLoadEnv(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("Error loading .env file: %s", err)
	}
}

func TestBuildDSN(t *testing.T) {
	expectedDSN := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "postgres")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("DB_TIMEZONE", "America/Sao_Paulo")

	dsn := config.BuildDSN()
	if dsn != expectedDSN {
		t.Errorf("DSN was incorrect, got: %s, want: %s", dsn, expectedDSN)
	}
}

func TestConnectDB(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "postgres")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("DB_TIMEZONE", "America/Sao_Paulo")

	db, err := config.ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to retrieve SQL DB from GORM: %s", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("Failed to close database connection: %s", err)
	}
}
