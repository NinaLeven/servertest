package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os/signal"
	"strings"
	"syscall"

	"servertest/internal/api"
	"servertest/internal/memstorage"
	"servertest/internal/servertest"
)

func main() {
	configPath := flag.String("c", "config/dev.yaml", "Path to config .yaml")

	flag.Parse()

	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	storage := memstorage.NewMemStorage()
	controller := servertest.NewController()

	apiServer := api.NewServer(storage, controller)

	handler := api.HandlerWithOptions(apiServer, api.ChiServerOptions{
		BaseURL: strings.TrimPrefix(config.BasePath, "/"),
	})

	ctx := context.Background()
	serverCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	httpServer := http.Server{
		Addr:    config.HttpAddr,
		Handler: handler,
		BaseContext: func(l net.Listener) context.Context {
			return serverCtx
		},
	}
	defer httpServer.Close()

	go func() {
		<-serverCtx.Done()
		err = httpServer.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
