package dao

import (
	"errors"
)

// Session is a 20 questions game session
type Session struct {
	ID        string      `json:"id"`
	Word      string      `json:"-"`
	Questions []*Question `json:"questions"`
	Answered  bool        `json:"answered"`
}

// Question is a question from a 20 questions game
type Question struct {
	ID        int    `json:"id"`
	SessionID string `json:"-"`
	Question  string `json:"question"`
	Answer    *bool  `json:"answer"`
}

// ErrSessionNotFound should be returned when the session is not found
var ErrSessionNotFound = errors.New("Session not found")

// ErrQuestionNotFound should be returned when a question is not found
var ErrQuestionNotFound = errors.New("Question not found")

// ErrNoMoreQuestions should be returned when attempting to save a questions where 20 have already been asked
var ErrNoMoreQuestions = errors.New("20 Questions have been asked already")

// DAO is the Data Access Object for the 20 questions game
type DAO interface {
	CreateSession(id string, word string) error
	UpdateSession(id string, answered bool) error
	SaveQuestion(sessionID string, question string) (int, error)
	SaveAnswer(sessionID string, questionID int, answer bool) error
	GetSession(id string) (*Session, error)
}
