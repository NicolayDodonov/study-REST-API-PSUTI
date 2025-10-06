package handler

import (
	"net/http"
	"study-REST-API-PSUTI/internal/logger"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func New(db *sqlx.DB, logger *logger.Logger) *Handler {
	return &Handler{db: db, logger: logger}
}

func TODO(w http.ResponseWriter, r *http.Request) {}
