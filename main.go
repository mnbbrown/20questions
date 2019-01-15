package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	chiMiddle "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/mnbbrown/20questions/dao"
	"github.com/rs/cors"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// Service is the root service for 20 questions
type Service struct {
	d dao.DAO
}

// NewService creates a new root service
func NewService(d dao.DAO) *Service {
	return &Service{d}
}

// Routes returns the root service router
func (s Service) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/session/{id}", s.getSession)
	r.Post("/session", s.createSession)
	r.Post("/session/{id}/questions", s.createQuestion)
	r.Post("/session/{id}/questions/{questionID}", s.answerQuestion)
	return r
}

type getSessionResponse struct {
	*dao.Session
}

func (gs *getSessionResponse) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func (s Service) getSession(rw http.ResponseWriter, req *http.Request) {
	sessionID := chi.URLParam(req, "id")
	sess, err := s.d.GetSession(sessionID)
	if err != nil {
		if err == dao.ErrSessionNotFound {
			render.Render(rw, req, ErrNotFound(err))
			return
		}
		render.Render(rw, req, &APIError{Status: 500, Message: "Internal Server Error"})
		return
	}
	render.Render(rw, req, &getSessionResponse{Session: sess})
}

type createSessionRequest struct {
	Word string `json:"word"`
}

func (csr *createSessionRequest) Bind(req *http.Request) error {
	return nil
}

type createSessionResponse struct {
	*dao.Session
}

func (csr *createSessionResponse) Render(rw http.ResponseWriter, req *http.Request) error {
	return nil
}

func (s Service) createSession(rw http.ResponseWriter, req *http.Request) {
	var requestBody createSessionRequest
	if err := render.Bind(req, &requestBody); err != nil {
		log.Error(err)
		render.Render(rw, req, &APIError{Status: 500, Message: "Internal Server Error"})
		return
	}
	id := uuid.NewV4().String()
	if err := s.d.CreateSession(id, requestBody.Word); err != nil {
		log.Error(err)
		render.Render(rw, req, &APIError{Status: 500, Message: "Internal Server Error"})
		return
	}
	respBody := &createSessionResponse{&dao.Session{
		ID:        id,
		Word:      requestBody.Word,
		Questions: make([]*dao.Question, 0),
	}}
	render.Render(rw, req, respBody)
}

type createQuestionRequest struct {
	Question string `json:"question"`
	Guess    string `json:"guess"`
}

func (c *createQuestionRequest) Bind(r *http.Request) error {
	if c.Question != "" && c.Guess != "" {
		return errors.New("guess and question can't both have values")
	}
	return nil
}

type createQuestionResponse struct {
	*dao.Question
	Correct bool `json:"correct"`
}

func (c *createQuestionResponse) Render(rw http.ResponseWriter, req *http.Request) error {
	return nil
}

func (s Service) createQuestion(rw http.ResponseWriter, req *http.Request) {
	var requestBody createQuestionRequest
	if err := render.Bind(req, &requestBody); err != nil {
		log.Error(err)
		render.Render(rw, req, &APIError{Status: 400, Message: err.Error()})
		return
	}
	sessionID := chi.URLParam(req, "id")
	sess, err := s.d.GetSession(sessionID)
	if err != nil {
		render.Render(rw, req, ErrNotFound(err))
		return
	}
	if sess.Answered {
		render.Render(rw, req, &APIError{Status: 400, Message: "Bad Request"})
		return
	}
	if requestBody.Guess != "" {
		if sess.Word == requestBody.Guess {
			err = s.d.UpdateSession(sessionID, true)
			if err != nil {
				render.Render(rw, req, ErrNotFound(err))
				return
			}
			render.Render(rw, req, &createQuestionResponse{
				Correct: true,
			})
			return
		}
	}

	var id int
	if id, err = s.d.SaveQuestion(sessionID, requestBody.Question); err != nil {
		log.Error(err)
		render.Render(rw, req, &APIError{Status: 500, Message: "Internal Server Error"})
		return
	}
	render.Render(rw, req, &createQuestionResponse{Question: &dao.Question{
		Question: requestBody.Question,
		ID:       id,
	}, Correct: false})
}

type answerQuestionRequest struct {
	Answer bool `json:"answer"`
}

func (a *answerQuestionRequest) Bind(req *http.Request) error {
	return nil
}

func (s Service) answerQuestion(rw http.ResponseWriter, req *http.Request) {
	var requestBody answerQuestionRequest
	if err := render.Bind(req, &requestBody); err != nil {
		log.Error(err)
		render.Render(rw, req, &APIError{Status: 400, Message: "Bad Request"})
		return
	}

	sessionID := chi.URLParam(req, "id")
	questionID, err := strconv.ParseInt(chi.URLParam(req, "questionID"), 10, 32)
	if err != nil {
		log.Error(err)
		render.Render(rw, req, &APIError{Status: 400, Message: "Bad Request"})
		return
	}
	err = s.d.SaveAnswer(sessionID, int(questionID), requestBody.Answer)
	if err != nil {
		if err == dao.ErrSessionNotFound || err == dao.ErrQuestionNotFound {
			render.Render(rw, req, ErrNotFound(err))
			return
		}
		render.Render(rw, req, &APIError{Status: 500, Message: "Internal Server Error"})
		return
	}
}

func main() {
	d := dao.NewMemoryDAO()
	s := NewService(d)

	r := chi.NewRouter()
	r.Use(chiMiddle.RequestID)
	r.Use(chiMiddle.RealIP)
	r.Use(chiMiddle.Logger)
	r.Use(chiMiddle.Recoverer)
	r.Use(chiMiddle.Timeout(60 * time.Second))
	r.Mount("/", s.Routes())
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	log.Printf("listening on %d", 8080)
	log.Errorln(http.ListenAndServe(fmt.Sprintf(":%d", 8080), c.Handler(r)))
}
