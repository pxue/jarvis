package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/goware/lg"
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

	router := web.New(&web.Handler{})
	hostURL := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if err := lg.SetLevelString("debug"); err != nil {
		lg.Fatal(err)
	}

	lg.Infof("Jarvis started")
	if err := graceful.ListenAndServe(hostURL, router); err != nil {
		lg.Fatal(err)
	}

	graceful.Wait()
}
