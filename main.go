package main

import (
	"context"
	"employee-worklog-service/api/router"
	"employee-worklog-service/config"
	"employee-worklog-service/utils/logger"
	"employee-worklog-service/utils/validator"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
    c := config.New()
    l := logger.New(c.Server.Debug)
    v := validator.New()

    r := router.New(l, v)

    s := &http.Server{
        Addr: fmt.Sprintf(":%d", c.Server.Port),
        Handler: r,
        ReadTimeout: c.Server.ReadTimeout,
        WriteTimeout: c.Server.WriteTimeout,
        IdleTimeout: c.Server.IdleTimeout,
    }

    closed := make(chan struct{})
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
        <-sigint

        l.Info().Msgf("Shutting down server %v", s.Addr)

        ctx, cancel := context.WithTimeout(context.Background(), c.Server.ReadTimeout)
        defer cancel()

        if err := s.Shutdown(ctx); err != nil {
            l.Error().Msg("Server shutdown failed")
        }

        close(closed)
    }()

    l.Info().Msgf("Starting server %v", s.Addr)
    if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        l.Fatal().Err(err).Msg("Server startup failed")
    }

    <-closed
    l.Info().Msgf("Server shutdown successfully")
}
