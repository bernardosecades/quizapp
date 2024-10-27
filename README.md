# Overview of the QuizApp

- API Server:
Provides endpoints for health check, fetching questions, and submitting answers.
- CLI Application:
Interacts with the API server to start a quiz, retrieve questions, and send user answers.

## API 

Run API:

`go run cmd/api/main.go `

### Health Endpoint

Request:

`curl http://localhost:8082/health`

Response:

```json
{
"status": "Ok"
}
```

### Get Questions Endpoint

Request:

`curl http://localhost:8082/v1/questions` 

Response:

```json
[
  {
    "id": "d920a7a3-f6d2-4ce9-97f1-f1f83fcaf49f",
    "text": "What is 2 * 2?",
    "options": [
      "3",
      "4",
      "5"
    ]
  },
  {
    "id": "d01bc313-5182-4b1e-8657-d9c5f759d68e",
    "text": "What is the capital of Asturias?",
    "options": [
      "Oviedo",
      "Gijón",
      "Avilés"
    ]
  }
]
```

### Submit Answers Endpoint

Request:

`curl http://localhost:8082/v1/answers \
--header 'Content-Type: application/json' \
--data '[
    {
        "questionID": "d01bc313-5182-4b1e-8657-d9c5f759d68e",
        "value": 0
    },
    {
        "questionID": "d920a7a3-f6d2-4ce9-97f1-f1f83fcaf49f",
        "value": 1
    }
]'`

Response:

```json
{
  "totalQuestions": 2,
  "correctAnswers": 2,
  "scorePercentage": 100,
  "comparativeScore": 0
}
```

## CLI

Run CLI application:

`go run cmd/cli/main.go start-quiz`

## Tech

- Execute test: `make test-unit`
- See coverage: `make coverage`
- Run Swagger UI - open api v3: `make run-openapi-ui`
- You have docker files to API and CLI so if you dont have installed go 1.23 you can create containers in your machine
- We are using github actions to execute: linter and tests before upload docker image in docker hub account 
