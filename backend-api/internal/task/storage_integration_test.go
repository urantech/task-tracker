package task

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"backend-api/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestStorage_GetDailyReports(t *testing.T) {
	// Arrange
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:15-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	defer func() {
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("pgx", connStr)
	require.NoError(t, err)
	defer db.Close()

	goose.SetBaseFS(migrations.EmbedMigrations)
	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	err = goose.Up(db, ".")
	require.NoError(t, err)

	storage := NewStorage(db)

	user1ID := createTestUser(t, db, "user1@example.com")
	user2ID := createTestUser(t, db, "user2@example.com")
	user3ID := createTestUser(t, db, "user3@example.com")
	user4ID := createTestUser(t, db, "user4@example.com")

	today := time.Now().Format("2006-01-02")

	// User 1: 2 pending, 1 completed today
	createTestTask(t, db, user1ID, "Task 1", StatusTodo, "2024-01-01")
	createTestTask(t, db, user1ID, "Task 2", StatusInProgress, "2024-01-01")
	createTestTask(t, db, user1ID, "Task 3", StatusDone, today)

	// User 2: 1 pending
	createTestTask(t, db, user2ID, "Task 4", StatusTodo, "2024-01-01")

	// User 3: 1 completed today
	createTestTask(t, db, user3ID, "Task 5", StatusDone, today)

	// User 4: no tasks

	// Act
	reports, err := storage.GetDailyReports(ctx)
	require.NoError(t, err)

	// Assert
	assert.Len(t, reports, 3)

	reportMap := make(map[int64]DailyReport)
	for _, r := range reports {
		reportMap[r.UserId] = r
	}

	r1, exists := reportMap[user1ID]
	assert.True(t, exists)
	assert.Equal(t, 2, r1.PendingCount)
	assert.Equal(t, 1, r1.CompletedCount)

	r2, exists := reportMap[user2ID]
	assert.True(t, exists)
	assert.Equal(t, 1, r2.PendingCount)
	assert.Equal(t, 0, r2.CompletedCount)

	r3, exists := reportMap[user3ID]
	assert.True(t, exists)
	assert.Equal(t, 0, r3.PendingCount)
	assert.Equal(t, 1, r3.CompletedCount)

	_, exists = reportMap[user4ID]
	assert.False(t, exists)
}

func createTestUser(t *testing.T, db *sql.DB, email string) int64 {
	var id int64
	err := db.QueryRow(`
		INSERT INTO users (email, password, role_id)
		VALUES ($1, $2, (SELECT id FROM roles WHERE name = 'ROLE_USER'))
		RETURNING id
	`, email, "password").Scan(&id)
	require.NoError(t, err)
	return id
}

func createTestTask(t *testing.T, db *sql.DB, userID int64, title string, status TaskStatus, updatedAt string) {
	_, err := db.Exec(`
		INSERT INTO tasks (title, description, status, user_id, created_at, updated_at)
		VALUES ($1, '', $2, $3, $4, $4)
	`, title, status, userID, updatedAt)
	require.NoError(t, err)
}
