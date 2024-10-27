//go:build unit

package service_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/quizapp/internal/api/model"
	"github.com/quizapp/internal/api/service"
	"github.com/quizapp/internal/api/service/mocks"
	"github.com/quizapp/pkg/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuizService_Success_EvaluateAnswers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	question1 := "2a71d7f5-3b8c-4d5c-b288-a69361ccfbd2"
	question2 := "f1684e1a-138c-4568-bcf5-57b8eb8b0c59"

	mockRepo := mocks.NewMockQuestionsRepository(ctrl)

	mockRepo.EXPECT().GetQuestions(gomock.Any()).Return([]*model.Question{
		{
			ID: question1,
			Options: []string{
				"option A",
				"option B",
			},
			CorrectIdx: 0,
		},
		{
			ID: question2,
			Options: []string{
				"option A",
				"option B",
				"option C",
			},
			CorrectIdx: 2,
		},
	}, nil).Times(2)

	srv := service.NewQuizService(mockRepo)

	answersUser1 := []dto.Answer{
		{
			QuestionID: question1,
			Value:      0, // Correct
		},
		{
			QuestionID: question2,
			Value:      1, // Incorrect
		},
	}
	result, err := srv.EvaluateAnswers(context.Background(), answersUser1)

	assert.Nil(t, err)
	assert.Equal(t, 1, result.CorrectAnswers)
	assert.Equal(t, 50.0, result.ScorePercentage)
	assert.Equal(t, 0.0, result.ComparativeScore)
	assert.Equal(t, 2, result.TotalQuestions)

	answersUser2 := []dto.Answer{
		{
			QuestionID: question1,
			Value:      0, // Correct
		},
		{
			QuestionID: question2,
			Value:      2, // Incorrect
		},
	}

	result, err = srv.EvaluateAnswers(context.Background(), answersUser2)

	assert.Nil(t, err)
	assert.Equal(t, 2, result.CorrectAnswers)
	assert.Equal(t, 100.0, result.ScorePercentage)
	assert.Equal(t, 100.0, result.ComparativeScore)
	assert.Equal(t, 2, result.TotalQuestions)
}

func TestQuizService_Error_EvaluateAnswers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test data for questions
	question1 := "2a71d7f5-3b8c-4d5c-b288-a69361ccfbd2"
	question2 := "f1684e1a-138c-4568-bcf5-57b8eb8b0c59"

	mockRepo := mocks.NewMockQuestionsRepository(ctrl)
	mockRepo.EXPECT().GetQuestions(gomock.Any()).Return([]*model.Question{
		{
			ID: question1,
			Options: []string{
				"option A",
				"option B",
			},
			CorrectIdx: 0,
		},
		{
			ID: question2,
			Options: []string{
				"option A",
				"option B",
			},
			CorrectIdx: 1,
		},
	}, nil).AnyTimes()

	srv := service.NewQuizService(mockRepo)

	// Define test cases
	tests := []struct {
		name    string
		answers []dto.Answer
		err     error
	}{
		{
			name: "ErrMissingAnswers",
			answers: []dto.Answer{
				{QuestionID: question1, Value: 0},
			},
			err: service.ErrMissingAnswers,
		},
		{
			name: "ErrInvalidAnswerOptions",
			answers: []dto.Answer{
				{QuestionID: question1, Value: 0},
				{QuestionID: question2, Value: 8}, // invalid option
			},
			err: service.ErrInvalidAnswerOptions,
		},
		{
			name: "ErrQuestionNotFound",
			answers: []dto.Answer{
				{QuestionID: question1, Value: 0},
				{QuestionID: "id-does-not-exist", Value: 0},
			},
			err: service.ErrQuestionNotFound,
		},
	}

	// Run each test case as a subtest
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := srv.EvaluateAnswers(context.Background(), tc.answers)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestQuizService_GetQuestions(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuestionsRepository(ctrl)

	mockRepo.EXPECT().GetQuestions(gomock.Any()).Return([]*model.Question{{CorrectIdx: 3}, {CorrectIdx: 1}}, nil).Times(1)

	srv := service.NewQuizService(mockRepo)

	questions, err := srv.GetQuestions(context.Background())

	assert.Len(t, questions, 2)
	assert.Nil(t, err)
}
