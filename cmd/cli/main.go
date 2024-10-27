package main

import (
	"context"
	"log"
	"time"

	"github.com/quizapp/internal/cli/client"
	"github.com/quizapp/internal/cli/command"
	"github.com/quizapp/pkg/env"
	"github.com/spf13/cobra"
)

const (
	timeoutAPIRequest = 5 * time.Second
)

func main() {
	ctx := context.Background()

	baseURL := env.GetEnv("BASE_URL", "http://localhost:8082")
	quizClient := client.NewClient(baseURL, timeoutAPIRequest)
	rootCmd := &cobra.Command{Use: "quiz-cli"}

	cmdStartQuiz := command.StartQuiz(quizClient)
	cmdStartQuiz.SetContext(ctx)

	rootCmd.AddCommand(command.StartQuiz(quizClient))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
