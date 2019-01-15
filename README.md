## 20questions

Simple 20 questions game that can be played through an API.

### Running the game

```bash
brew install httpie
docker run -p 8080:8080 mnbbrown/20questions
```

### Playing the game

#### Player 1 (starting a new game)

Request: `http POST http://127.0.0.1:8080/session word=pineapple`
Response: 200 `{ "id": "uuid-uuid", questions: [], answered: false }`

#### Player 2 (guessing incorrectly and asking a question)

a) Guess

Request: `http POST http://127.0.0.1:8080/session/uuid-uuid/questions guess=apple`
Response: 200 `{ "correct": false }`

b) Ask a question

Request: `http POST http://127.0.0.1:8080/session/uuid-uuid/questions question="is is a fruit?"`
Response: 204

#### Player 1 (responding to the question)

a) list questions

Request: `http GET http://127.0.0.1:8080/session/uuid-uuid guess=pineapple`
Response: 200 `{ "id": "uuid-uuid", "questions": [ { id: 1, "question": "is it a fruit?", "answer": null } ], answered: false }`

b) respond to question

Request: `http POST http://127.0.0.1:8080/session/uuid-uuid/questions/1 answer=false`
Response: 204

#### Player 2 (guessing correctly)

Request: `http POST http://127.0.0.1:8080/session/uuid-uuid/questions guess=pineapple`
Response: 200 `{ "correct": true }`
