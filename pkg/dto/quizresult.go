package dto

type QuizResult struct {
	TotalQuestions   int     `json:"totalQuestions"`
	CorrectAnswers   int     `json:"correctAnswers"`
	ScorePercentage  float64 `json:"scorePercentage"`
	ComparativeScore float64 `json:"comparativeScore"`
}
