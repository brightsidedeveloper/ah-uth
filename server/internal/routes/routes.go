package routes

import (
	"server/internal/handler"

	"github.com/go-chi/chi/v5"
)

func MountRoutes(r chi.Router, h *handler.Handler) {

	api := chi.NewRouter()
	api.Get("/users", h.GetUsers)
	api.Post("/user", h.PostUser)

	api.Get("/auth/{provider}/callback", h.GetAuthCallback)
	api.Get("/auth/{provider}", h.GetAuth)
	api.Get("/auth/{provider}/logout", h.GetLogout)
	api.Get("/auth/gothic", h.GetUserStore)

	r.Mount("/api", api)
}
