package quiz

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quizapp/internal/api/model"
	"github.com/quizapp/internal/api/service"
	"github.com/quizapp/internal/api/service/mocks"
	"github.com/quizapp/pkg/dto"
	"github.com/stretchr/testify/assert"
)

// TestQuestions tests the Questions method of the Handler.
func TestQuestions(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuestionsRepository(ctrl)

	mockRepo.EXPECT().GetQuestions(gomock.Any()).Return([]*model.Question{
		{
			ID:   "2a71d7f5-3b8c-4d5c-b288-a69361ccfbd2",
			Text: "Question 1",
			Options: []string{
				"option A",
				"option B",
			},
			CorrectIdx: 0,
		},
	}, nil).Times(1)

	srv := service.NewQuizService(mockRepo)

	handler := NewHandler(srv)

	req, err := http.NewRequest("GET", "/questions", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.Questions(rr, req)

	// Checks
	assert.Equal(t, http.StatusOK, rr.Code)

	var questions []dto.Question
	err = json.Unmarshal(rr.Body.Bytes(), &questions)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(questions))
	assert.Equal(t, "Question 1", questions[0].Text)
}

func TestAnswers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuestionsRepository(ctrl)

	mockRepo.EXPECT().GetQuestions(gomock.Any()).Return([]*model.Question{
		{
			ID:   "2a71d7f5-3b8c-4d5c-b288-a69361ccfbd2",
			Text: "Question 1",
			Options: []string{
				"option A",
				"option B",
			},
			CorrectIdx: 0,
		},
		{
			ID:   "2ac2e2da-1fa1-4f0a-aa36-f80bae6277ee",
			Text: "Question 2",
			Options: []string{
				"option A",
				"option B",
				"option C",
			},
			CorrectIdx: 1,
		},
	}, nil).Times(1)

	srv := service.NewQuizService(mockRepo)

	handler := NewHandler(srv)

	answers := []dto.Answer{
		{QuestionID: "2a71d7f5-3b8c-4d5c-b288-a69361ccfbd2", Value: 0},
		{QuestionID: "2ac2e2da-1fa1-4f0a-aa36-f80bae6277ee", Value: 2},
	}
	body, err := json.Marshal(answers)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/answers", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Answers(rr, req)

	// Checks
	assert.Equal(t, http.StatusOK, rr.Code)

	var quizResult dto.QuizResult
	err = json.Unmarshal(rr.Body.Bytes(), &quizResult)
	assert.NoError(t, err)
	assert.Equal(t, 1, quizResult.CorrectAnswers)
	assert.Equal(t, 2, quizResult.TotalQuestions)
	assert.Equal(t, 50.0, quizResult.ScorePercentage)
	assert.Equal(t, 0.0, quizResult.ComparativeScore)
}
