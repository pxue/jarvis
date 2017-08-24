package slack

import (
	"context"
	"net/http"

	"github.com/goji/param"
	"github.com/goware/lg"
	"github.com/pxue/jarvis/lib/ws"
)

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

func SlackCtx(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			ws.Respond(w, http.StatusBadRequest, err)
			return
		}

		payload := &SlackForm{}
		if err := param.Parse(r.PostForm, payload); err != nil {
			ws.Respond(w, http.StatusBadRequest, err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "slack", payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(handler)
}

func GenericHandler(w http.ResponseWriter, r *http.Request) {
	slack := r.Context().Value("slack").(*SlackForm)

	lg.Warnf("%+v", slack)
}
