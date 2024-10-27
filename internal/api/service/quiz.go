package service

import (
	"context"
	"errors"
	"sync"

	"github.com/quizapp/internal/api/model"
	"github.com/quizapp/pkg/dto"
)

var (
	ErrMissingAnswers       = errors.New("missing answers")
	ErrQuestionNotFound     = errors.New("question not found")
	ErrInvalidAnswerOptions = errors.New("invalid answer options")
)

//go:generate mockgen -destination=mocks/mock_questions_repository.go -package=mocks github.com/quizapp/internal/api/service QuestionsRepository
type QuestionsRepository interface {
	GetQuestions(ctx context.Context) ([]*model.Question, error)
}

type QuizService struct {
	questionsRepository QuestionsRepository
	scores              []float64 // Track scores for comparison
	sync.Mutex
}

func NewQuizService(questionsRepository QuestionsRepository) *QuizService {
	return &QuizService{
		questionsRepository: questionsRepository,
		scores:              []float64{},
	}
}

func (qs *QuizService) GetQuestions(ctx context.Context) ([]*model.Question, error) {
	return qs.questionsRepository.GetQuestions(ctx)
}

func (qs *QuizService) EvaluateAnswers(ctx context.Context, answers []dto.Answer) (*model.QuizResult, error) {
	questions, err := qs.GetQuestions(ctx)
	if err != nil {
		return nil, err
	}

	if len(answers) != len(questions) {
		return nil, ErrMissingAnswers
	}

	questionsWithKeyMap := sliceToMap(questions)

	correctCount := 0
	for _, answer := range answers {
		question, ok := questionsWithKeyMap[answer.QuestionID]
		if !ok {
			return nil, ErrQuestionNotFound
		}
		if answer.Value >= len(question.Options) {
			return nil, ErrInvalidAnswerOptions
		}
		if question.CorrectIdx == answer.Value {
			correctCount++
		}
	}

	totalQuestions := len(questions)
	score := float64(correctCount) / float64(totalQuestions) * 100

	qs.Lock()
	comparativeScore := calculateComparativeScore(score, qs.scores)
	qs.scores = append(qs.scores, score)
	qs.Unlock()

	return &model.QuizResult{
		TotalQuestions:   totalQuestions,
		CorrectAnswers:   correctCount,
		ScorePercentage:  score,
		ComparativeScore: comparativeScore,
	}, nil
}

// calculateComparativeScore returns the percentage of scores in allScores
// that are lower than userScore. If allScores is empty, it returns 0.
func calculateComparativeScore(userScore float64, allScores []float64) float64 {
	if len(allScores) == 0 {
		return 0
	}
	betterCount := 0
	for _, score := range allScores {
		if userScore > score {
			betterCount++
		}
	}
	return float64(betterCount) / float64(len(allScores)) * 100
}

func sliceToMap(questions []*model.Question) map[string]*model.Question {
	result := make(map[string]*model.Question)
	for _, question := range questions {
		result[question.ID] = question
	}
	return result
}
