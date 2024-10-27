package command

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/quizapp/pkg/dto"
	"github.com/spf13/cobra"
)

// QuizClient defines the methods needed for a quiz client
//
//go:generate mockgen -destination=mocks/mock_quiz_client.go -package=mocks github.com/quizapp/internal/cli/command QuizClient
type QuizClient interface {
	GetAllQuestions(ctx context.Context) ([]dto.Question, error)
	SendAnswers(ctx context.Context, answers []dto.Answer) (dto.QuizResult, error)
}

// StartQuiz initializes the quiz command with a client dependency
func StartQuiz(client QuizClient) *cobra.Command {
	return &cobra.Command{
		Use:   "start-quiz",
		Short: "Start the quiz",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runQuiz(cmd, client)
		},
	}
}

// runQuiz executes the quiz-taking logic
func runQuiz(cmd *cobra.Command, client QuizClient) error {
	ctx := cmd.Context()
	questions, err := client.GetAllQuestions(ctx)
	if err != nil {
		return reportError(cmd, "error fetching questions", err)
	}

	answers, err := gatherAnswers(cmd, questions)
	if err != nil {
		return reportError(cmd, "error gathering answers", err)
	}

	result, err := client.SendAnswers(ctx, answers)
	if err != nil {
		return reportError(cmd, "error sending answers", err)
	}

	printResult(cmd, result)

	return nil
}

// gatherAnswers reads user input for each question and validates the answer
func gatherAnswers(cmd *cobra.Command, questions []dto.Question) ([]dto.Answer, error) {
	answers := make([]dto.Answer, 0, len(questions))
	reader := bufio.NewReader(cmd.InOrStdin())
	for _, q := range questions {
		reportInfo(cmd, "Question: %s: %s\n", q.ID, q.Text)
		for i, opt := range q.Options {
			reportInfo(cmd, "  %d: %s\n", i, opt)
		}

		val, err := promptAnswer(cmd, reader, len(q.Options))
		if err != nil {
			return nil, err
		}

		answers = append(answers, dto.Answer{
			QuestionID: q.ID,
			Value:      val,
		})
	}

	return answers, nil
}

// promptAnswer prompts the user for input and validates the response
func promptAnswer(cmd *cobra.Command, reader *bufio.Reader, numOptions int) (int, error) {
	for {
		reportInfo(cmd, "Enter the number of your answer: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		val, err := strconv.Atoi(input)
		if err == nil && val >= 0 && val < numOptions {
			return val, nil
		}
		reportInfo(cmd, "Invalid input, please enter a positive integer.")
	}
}

// printResult prints the quiz result summary
func printResult(cmd *cobra.Command, result dto.QuizResult) {
	reportInfo(cmd, "You got %d out of %d correct!\n", result.CorrectAnswers, result.TotalQuestions)
	reportInfo(cmd, "Score: %.2f%%\n", result.ScorePercentage)
	reportInfo(cmd, "You were better than %.2f%% of all quizzers\n", result.ComparativeScore)
}

func reportError(cmd *cobra.Command, message string, err error) error {
	fmt.Fprintf(cmd.OutOrStdout(), "%s: %v", message, err)
	return err
}

func reportInfo(cmd *cobra.Command, message string, args ...interface{}) {
	fmt.Fprintf(cmd.OutOrStdout(), message, args...)
}
