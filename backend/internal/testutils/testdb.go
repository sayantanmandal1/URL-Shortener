package testutils

import (
	"context"
	"testing"
	"time"

	"url-shortener-backend/internal/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB represents a test database instance
type TestDB struct {
	Container testcontainers.Container
	DB        *database.DB
	URL       string
}

// SetupTestDB creates a test database using testcontainers
func SetupTestDB(t *testing.T) *TestDB {
	ctx := context.Background()

	// Create PostgreSQL container
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("Failed to start postgres container: %v", err)
	}

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Create connection pool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("Failed to create connection pool: %v", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	db := &database.DB{Pool: pool}

	// Run migrations
	if err := runMigrations(ctx, db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return &TestDB{
		Container: postgresContainer,
		DB:        db,
		URL:       connStr,
	}
}

// Cleanup closes the database connection and stops the container
func (tdb *TestDB) Cleanup(t *testing.T) {
	ctx := context.Background()

	if tdb.DB != nil && tdb.DB.Pool != nil {
		tdb.DB.Pool.Close()
	}

	if tdb.Container != nil {
		if err := tdb.Container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}
}

// runMigrations runs the database migrations
func runMigrations(ctx context.Context, db *database.DB) error {
	// Create tables
	schema := `
	-- Links table
	CREATE TABLE IF NOT EXISTS links (
		id VARCHAR(36) PRIMARY KEY,
		original_url TEXT NOT NULL,
		slug VARCHAR(50) UNIQUE NOT NULL,
		title VARCHAR(200),
		description TEXT,
		click_count INTEGER DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Clicks table for analytics
	CREATE TABLE IF NOT EXISTS clicks (
		id VARCHAR(36) PRIMARY KEY,
		link_id VARCHAR(36) REFERENCES links(id) ON DELETE CASCADE,
		clicked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		ip_address INET,
		user_agent TEXT,
		referrer TEXT,
		country VARCHAR(2),
		device_type VARCHAR(20)
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_links_slug ON links(slug);
	CREATE INDEX IF NOT EXISTS idx_links_created_at ON links(created_at);
	CREATE INDEX IF NOT EXISTS idx_clicks_link_id ON clicks(link_id);
	CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at);

	-- Function to update updated_at timestamp
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ language 'plpgsql';

	-- Trigger to automatically update updated_at
	CREATE TRIGGER update_links_updated_at 
		BEFORE UPDATE ON links 
		FOR EACH ROW 
		EXECUTE FUNCTION update_updated_at_column();
	`

	_, err := db.Pool.Exec(ctx, schema)
	return err
}

// ClearTables clears all test data from tables
func (tdb *TestDB) ClearTables(t *testing.T) {
	ctx := context.Background()

	queries := []string{
		"DELETE FROM clicks",
		"DELETE FROM links",
	}

	for _, query := range queries {
		if _, err := tdb.DB.Pool.Exec(ctx, query); err != nil {
			t.Fatalf("Failed to clear table: %v", err)
		}
	}
}

// CreateTestLink creates a test link in the database
func (tdb *TestDB) CreateTestLink(t *testing.T, id, originalURL, slug string) {
	ctx := context.Background()

	query := `
		INSERT INTO links (id, original_url, slug, title, description, click_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 0, NOW(), NOW())`

	_, err := tdb.DB.Pool.Exec(ctx, query, id, originalURL, slug, "Test Title", "Test Description")
	if err != nil {
		t.Fatalf("Failed to create test link: %v", err)
	}
}
