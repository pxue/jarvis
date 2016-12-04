package web

import (
	"net/http"

	"github.com/goji/param"
	"github.com/goware/lg"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pxue/jarvis/data"
	"github.com/pxue/jarvis/lib/ws"
)

type Handler struct {
	DB *data.Database
}

type SlackForm struct {
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
}

func GenericHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ws.Respond(w, http.StatusBadRequest, err)
		return
	}

	payload := &SlackForm{}
	if err := param.Parse(r.PostForm, payload); err != nil {
		ws.Respond(w, http.StatusBadRequest, err)
		return
	}

	lg.Warnf("%+v", payload)
}

func New(h *Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.WithValue("db", h.DB))

	r.Post("/", GenericHandler)

	return r
}
