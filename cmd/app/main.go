package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"golazy.dev/lazyroutes"
	appinit "sample_app/app/init"
)

func main() {
	ctx := appinit.Context(context.Background())
	mux := lazyroutes.New(ctx)
	appinit.Draw(ctx, mux)

	server := &http.Server{
		Addr:    listenAddr(),
		Handler: mux,
	}
	log.Printf("listening on %s", server.Addr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func listenAddr() string {
	addr := os.Getenv("ADDR")
	if addr == "" {
		return ":8080"
	}
	if _, err := strconv.ParseUint(addr, 10, 16); err == nil {
		return ":" + addr
	}
	return addr
}
