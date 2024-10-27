package model

import "time"

type Question struct {
	ID         string
	Text       string
	Options    []string
	CorrectIdx int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type QuizResult struct {
	TotalQuestions   int
	CorrectAnswers   int
	ScorePercentage  float64
	ComparativeScore float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
