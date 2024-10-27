package repository

import (
	"context"

	"github.com/quizapp/internal/api/model"
)

type InMemoryQuestionsRepository struct {
	questionCollection map[string]model.Question
}

func NewInMemoryQuestionsRepository() *InMemoryQuestionsRepository {
	return &InMemoryQuestionsRepository{
		questionCollection: map[string]model.Question{
			"d920a7a3-f6d2-4ce9-97f1-f1f83fcaf49f": {ID: "d920a7a3-f6d2-4ce9-97f1-f1f83fcaf49f", Text: "What is 2 * 2?", Options: []string{"3", "4", "5"}, CorrectIdx: 1},
			"d01bc313-5182-4b1e-8657-d9c5f759d68e": {ID: "d01bc313-5182-4b1e-8657-d9c5f759d68e", Text: "What is the capital of Asturias?", Options: []string{"Oviedo", "Gijón", "Avilés"}, CorrectIdx: 0},
		},
	}
}

func (r *InMemoryQuestionsRepository) GetQuestions(_ context.Context) ([]*model.Question, error) {
	result := make([]*model.Question, 0, len(r.questionCollection))
	for _, v := range r.questionCollection {
		result = append(result, &v)
	}
	return result, nil
}
