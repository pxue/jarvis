package web

import (
	"github.com/pressly/chi"
	"github.com/pxue/jarvis/lib/linkedin"
	"github.com/pxue/jarvis/lib/mls"
	"github.com/pxue/jarvis/web/slack"
)

type Handler struct{}

func New(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Route("/slack", func(r chi.Router) {
		r.Use(slack.SlackCtx)
		r.Post("/", slack.GenericHandler)
	})

	r.Route("/linkedin", func(r chi.Router) {
		r.Get("/oauth", linkedin.OAuth)
		r.Get("/oauth/callback", linkedin.OAuthCallback)

		//r.Get("/connected", linkedin.GetConnected)
	})

	r.Route("/mls", func(r chi.Router) {
		r.Get("/crawl", mls.ParseListings)
	})

	return r
}
