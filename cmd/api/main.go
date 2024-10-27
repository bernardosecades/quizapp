package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"

	"github.com/quizapp/pkg/env"

	"github.com/quizapp/internal/api/repository"

	"github.com/gorilla/mux"
	"github.com/quizapp/internal/api/server/handler/health"
	"github.com/quizapp/internal/api/server/handler/quiz"
	"github.com/quizapp/internal/api/server/middleware"
	"github.com/quizapp/internal/api/service"
	"github.com/rs/zerolog"
)

const timeOutHandlers = 5 * time.Second

func main() {
	// LOGGER
	loggerOutput := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(loggerOutput).With().Timestamp().Logger()

	// HANDLERS
	quizService := service.NewQuizService(repository.NewInMemoryQuestionsRepository())
	quizHandler := quiz.NewHandler(quizService)
	healthHandler := health.NewHandler()

	// ROUTER
	router := mux.NewRouter()

	v1 := router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/questions", quizHandler.Questions).Methods(http.MethodGet)
	v1.HandleFunc("/answers", quizHandler.Answers).Methods(http.MethodPost)

	router.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	// MIDDLEWARE
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Headers)

	// SERVER
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", env.GetEnv("PORT_API", "8082")),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// SHUTDOWN
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info().Msg(fmt.Sprintf("Received signal: %v. Initiating graceful shutdown...", sigChan))

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Fatal().Msg(fmt.Sprintf("HTTP shutdown error: %v", err))
		}
	}()

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:4000", // swagger
		},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	routerWithCors := c.Handler(router)

	// TIMEOUT endpoints handlers
	routerWithMiddlewares := http.TimeoutHandler(routerWithCors, timeOutHandlers, "Timeout!")
	http.Handle("/", routerWithMiddlewares)

	logger.Info().
		Str("PORT", server.Addr).
		Msg("HTTP server listening on port")

	// RUN SERVER
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal().
			Err(err).
			Msg("failed to start server")
	}
}
