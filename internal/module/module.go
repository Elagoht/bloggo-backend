package module

import (
	"database/sql"

	"github.com/go-chi/chi"
)

type Module interface {
	RegisterModule(database *sql.DB, router *chi.Mux)
}
