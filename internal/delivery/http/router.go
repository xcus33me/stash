package httpdelivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(rooms *RoomHandler, files *FileHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/rooms", func(r chi.Router) {
		r.Post("/", rooms.Create)

		r.Route("/{roomID}", func(r chi.Router) {
			r.Get("/", rooms.Get)
			r.Delete("/", rooms.Delete)
			r.Post("/lock", rooms.Lock)
			r.Post("/unlock", rooms.Unlock)

			r.Get("/files", files.List)
			r.Post("/files", files.Upload)
		})
	})

	r.Route("/files/{fileID}", func(r chi.Router) {
		r.Get("/", files.Download)
		r.Delete("/", files.Delete)
	})

	return r
}
