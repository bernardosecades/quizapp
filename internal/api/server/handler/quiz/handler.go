package quiz

import (
	"encoding/json"
	"net/http"

	"github.com/quizapp/internal/api/server/handler"

	"github.com/quizapp/internal/api/server/transformer"
	"github.com/quizapp/internal/api/service"
	"github.com/quizapp/pkg/dto"
)

type Handler struct {
	quizService *service.QuizService
}

// NewHandler initialize handler
func NewHandler(quizService *service.QuizService) *Handler {
	return &Handler{quizService: quizService}
}

// Questions handler to return all questions with all possible responses
func (qh *Handler) Questions(w http.ResponseWriter, r *http.Request) {

	questions, err := qh.quizService.GetQuestions(r.Context())
	if err != nil {
		handler.EncodeHTTPError(handler.NewHTTPError(err.Error(), http.StatusInternalServerError), w)
		return
	}
	err = json.NewEncoder(w).Encode(transformer.TransformQuestionsTo(questions))
	if err != nil {
		handler.EncodeHTTPError(handler.NewHTTPError(err.Error(), http.StatusInternalServerError), w)
		return
	}
}

// Answers handler to process answers and return statistics about the result
func (qh *Handler) Answers(w http.ResponseWriter, r *http.Request) {
	var answers []dto.Answer
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		handler.EncodeHTTPError(handler.NewHTTPError("invalid request", http.StatusBadRequest), w)
		return
	}

	quizResult, err := qh.quizService.EvaluateAnswers(r.Context(), answers)
	if err != nil {
		handler.EncodeHTTPError(handler.NewHTTPError(err.Error(), http.StatusBadRequest), w)
		return
	}

	err = json.NewEncoder(w).Encode(transformer.TransformQuizResultTo(quizResult))
	if err != nil {
		handler.EncodeHTTPError(handler.NewHTTPError(err.Error(), http.StatusInternalServerError), w)
		return
	}
}
