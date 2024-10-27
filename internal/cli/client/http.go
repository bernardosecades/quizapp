package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/quizapp/pkg/dto"
)

const (
	uriQuestions = "/v1/questions"
	uriAnswers   = "/v1/answers"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetAllQuestions(ctx context.Context) ([]dto.Question, error) {
	reqURL, err := url.JoinPath(c.baseURL, uriQuestions)
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, uriQuestions)
	}

	var questions []dto.Question
	if err := json.NewDecoder(resp.Body).Decode(&questions); err != nil {
		return nil, fmt.Errorf("failed to decode questions response: %w", err)
	}
	return questions, nil
}

func (c *Client) SendAnswers(ctx context.Context, answers []dto.Answer) (dto.QuizResult, error) {
	body, err := json.Marshal(answers)
	if err != nil {
		return dto.QuizResult{}, fmt.Errorf("failed to marshal answers: %w", err)
	}

	reqURL, err := url.JoinPath(c.baseURL, uriAnswers)
	if err != nil {
		return dto.QuizResult{}, fmt.Errorf("failed to join path: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return dto.QuizResult{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return dto.QuizResult{}, fmt.Errorf("failed to send answers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.QuizResult{}, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, uriAnswers)
	}

	var result dto.QuizResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return dto.QuizResult{}, fmt.Errorf("failed to decode quiz result: %w", err)
	}
	return result, nil

}
