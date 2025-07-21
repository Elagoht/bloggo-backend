package module

import (
	"github.com/go-chi/chi"
)

type Module interface {
	RegisterModule(router *chi.Mux)
}
