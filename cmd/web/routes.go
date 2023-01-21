package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	cors "github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*", "https://*", "http://*"},
		AllowOriginFunc:    nil,
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     nil,
		AllowCredentials:   false,
		MaxAge:             300,
		OptionsPassthrough: false,
		Debug:              false,
	}))

	mux.Use(SessionLoad)

	mux.Get("/", app.Home)
	mux.Get("/occasions", app.Occasions)
	mux.Get("/occasions/{id}", app.Occasion)
	mux.Post("/occasions", app.UpdateOccasion)

	mux.Get("/resize-image", app.ResizeForm)
	mux.Post("/api/resize", app.ResizeDo)

	mux.Get("/browse-photos", app.BrowsePhotos)
	mux.Post("/api/browse-photos", app.BrowseDirectory)

	mux.Get("/api/image", app.GetImage)

	return mux
}
