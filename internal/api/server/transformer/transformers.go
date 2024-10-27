package transformer

import (
	"github.com/quizapp/internal/api/model"
	"github.com/quizapp/pkg/dto"
)

func TransformQuestionsTo(questions []*model.Question) []dto.Question {
	result := make([]dto.Question, 0, len(questions))
	for _, question := range questions {
		result = append(result, dto.Question{
			ID:      question.ID,
			Text:    question.Text,
			Options: question.Options,
		})
	}
	return result
}

func TransformQuizResultTo(quizResult *model.QuizResult) dto.QuizResult {
	return dto.QuizResult{
		TotalQuestions:   quizResult.TotalQuestions,
		CorrectAnswers:   quizResult.CorrectAnswers,
		ScorePercentage:  quizResult.ScorePercentage,
		ComparativeScore: quizResult.ComparativeScore,
	}
}
