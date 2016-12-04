package main

import (
	"syscall"
	"time"

	"github.com/goware/lg"
	"github.com/pxue/jarvis/data"
	"github.com/pxue/jarvis/web"
	"github.com/zenazn/goji/graceful"
)

func main() {
	graceful.AddSignal(syscall.SIGINT, syscall.SIGTERM)
	graceful.Timeout(10 * time.Second) // Wait timeout for handlers to finish.
	graceful.PreHook(func() {
		lg.Info("waiting for requests to finish..")
	})
	graceful.PostHook(func() {
		lg.Info("finishing up...")
	})

	lg.Info("Jarvis starting on 5331")

	conf := &data.DBConf{}
	db, err := data.NewDBSession(conf)
	if err != nil {
		lg.Fatal(err)
	}

	router := web.New(&web.Handler{DB: db})
	if err := graceful.ListenAndServe(":5331", router); err != nil {
		lg.Fatal(err)
	}

	graceful.Wait()
}
