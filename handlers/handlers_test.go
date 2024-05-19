package handlers_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vodeno/handlers"
)

func mockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	if err != nil {
		return nil, nil
	}
	return gormDB, mock
}

func TestHandler_CreateEmail(t *testing.T) {
	db, mock := mockDB()
	handlerTest := handlers.New(db)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"emails\"").WithArgs("simple text", "jan.kowalski@example.com", sqlmock.AnyArg(), 1, "Interview").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	req, err := http.NewRequest("POST", "/api/messages", strings.NewReader("{\"email\":\"jan.kowalski@example.com\",\"title\":\"Interview\",\"content\":\"simple text\",\"mailing_id\":1, \"insert_time\": \"2020-04-24T05:42:38.725412916Z\"}"))
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerTest.CreateEmail)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
}
