package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type HTTPEnv struct {
	Port string
}

type HTTPServer struct {
	HTTPEnv
	handler http.Handler
}

func SetupHTTP(env HTTPEnv, handler http.Handler) HTTPServer {
	return HTTPServer{HTTPEnv: env, handler: handler}
}

func (s *HTTPServer) Serve(ctx context.Context) (err error) {
	srv := &http.Server{
		Addr:    ":" + s.Port,
		Handler: s.handler,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err)
		}
	}()

	<-ctx.Done()
	log.Info().Msg("graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Stack().Err(err).Msg("error during shutdown")
	}
	log.Info().Msg("...graceful shutdown done")

	return
}
