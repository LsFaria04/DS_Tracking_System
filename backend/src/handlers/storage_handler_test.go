package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	return gormDB, mock
}

func TestGetAllStorages_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gormDB, mock := newMockDB(t)

	// Expected rows
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Storage A").
		AddRow(2, "Storage B")

	// Expect the SQL query GORM will generate
	mock.ExpectQuery(`SELECT \* FROM "storages" ORDER BY id ASC`).
		WillReturnRows(rows)

	handler := StorageHandler{DB: gormDB}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAllStorages(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	mock.ExpectationsWereMet()
}

func TestGetAllStorages_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gormDB, mock := newMockDB(t)

	// Force DB error
	mock.ExpectQuery(`SELECT \* FROM "storages" ORDER BY id ASC`).
		WillReturnError(sql.ErrConnDone)

	handler := StorageHandler{DB: gormDB}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetAllStorages(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	mock.ExpectationsWereMet()
}