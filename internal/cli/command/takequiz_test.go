//go:build unit

package command_test

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quizapp/internal/cli/command"
	"github.com/quizapp/internal/cli/command/mocks"
	"github.com/quizapp/pkg/dto"
	"github.com/stretchr/testify/assert"
)

func TestStartQuiz_Success_Execute(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockQuizClient(ctrl)

	// Define expected questions
	expectedQuestions := []dto.Question{
		{ID: "1", Text: "What is the best fighter?", Options: []string{"Illia Topuria", "Max Holloway"}},
		{ID: "2", Text: "What is the capital of Asturias?", Options: []string{"Avilés", "Gijón", "Oviedo"}},
	}

	expectedAnswers := []dto.Answer{
		{QuestionID: "1", Value: 0}, // Answering with option 1 for first question
		{QuestionID: "2", Value: 2}, // Answering with option 2 for second question
	}

	expectedResult := dto.QuizResult{
		CorrectAnswers:   2,
		TotalQuestions:   2,
		ScorePercentage:  100.0,
		ComparativeScore: 100.0,
	}

	mockClient.EXPECT().GetAllQuestions(gomock.Any()).Return(expectedQuestions, nil).Times(1)
	mockClient.EXPECT().SendAnswers(gomock.Any(), expectedAnswers).Return(expectedResult, nil).Times(1)

	// Simulate user input
	userInput := "0\n2\n"                        // Simulate selecting option 1 for the first question and option 2 for the second question
	reader := bytes.NewReader([]byte(userInput)) // Create a new bytes.Reader from the user input

	// Execute the command
	cmd := command.StartQuiz(mockClient)
	cmd.SetIn(reader) // Set the input for the command

	var outputBuffer bytes.Buffer
	cmd.SetOut(&outputBuffer)

	err := cmd.Execute()
	assert.NoError(t, err)

	// Check the output
	expectedOutput := "You got 2 out of 2 correct!\nScore: 100.00%\nYou were better than 100.00% of all quizzers\n" // Adjust based on what your command prints
	actualOutput := outputBuffer.String()

	assert.Contains(t, actualOutput, expectedOutput, "Output did not match expected result.")
}
