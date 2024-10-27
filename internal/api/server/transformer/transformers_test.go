//go:build unit

package transformer

import (
	"github.com/quizapp/internal/api/model"
	"github.com/quizapp/pkg/dto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransformQuestionTo(t *testing.T) {
	modelQuestions := []*model.Question{
		{
			ID:         "123",
			Text:       "question 123",
			Options:    []string{"a", "b", "c"},
			CorrectIdx: 1,
		},
	}

	expected := []dto.Question{
		{
			ID:      "123",
			Text:    "question 123",
			Options: []string{"a", "b", "c"},
		},
	}

	result := TransformQuestionsTo(modelQuestions)

	assert.Equal(t, expected, result)
}

func TestTransformQuizResultTo(t *testing.T) {
	modelQuiz := &model.QuizResult{
		TotalQuestions:   10,
		CorrectAnswers:   5,
		ScorePercentage:  50,
		ComparativeScore: 0,
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
	}

	expected := dto.QuizResult{
		TotalQuestions:   10,
		CorrectAnswers:   5,
		ScorePercentage:  50,
		ComparativeScore: 0,
	}

	result := TransformQuizResultTo(modelQuiz)

	assert.Equal(t, expected, result)
}
