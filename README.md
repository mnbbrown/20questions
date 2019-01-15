## 20questions

Simple 20 questions game that can be played through an API.

### Running the game

```bash
brew install httpie
docker run -p 8080:8080 mnbbrown/20questions
```

### Playing the game

#### Player 1 (who has the answer)

Request: `http http://127.0.0.1:8080/session word=pineapple`
Response: `{ "id": "uuid-uuid", questions: [], answered: false }`

#### Player 2 (the guesser)

a) Guess

Request: `http http://127.0.0.1:8080/session/uuid-uuid/questions guess=pineapple`
Response: `{ "correct": false }`

b) Ask a question

Request: `http http://127.0.0.1:8080/session/uuid-uuid/questions question="is is a fruit?"`
Response: empty
