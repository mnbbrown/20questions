package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/mnbbrown/20questions/dao"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createService() chi.Router {
	return NewService(dao.NewMemoryDAO()).Routes()
}

func newPost(t *testing.T, url string, body interface{}) *http.Request {
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func Test_CreateSession(t *testing.T) {
	reqB := map[string]string{
		"word": "test_word",
	}
	req := newPost(t, "/session", reqB)
	rr := httptest.NewRecorder()
	createService().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("Invalid API response. Code not 200 - %v", rr.Code)
	}
	respB := map[string]interface{}{}
	if err := json.Unmarshal(rr.Body.Bytes(), &respB); err != nil {
		t.Fatal(err)
	}
	if _, ok := respB["word"]; ok {
		t.Fatal("Word should not be returned")
	}
	if _, ok := respB["id"]; !ok {
		t.Fatal("No session ID")
	}
}

func Test_CreateQuestion(t *testing.T) {
	d := dao.NewMemoryDAO()
	service := NewService(d).Routes()
	reqB := map[string]string{
		"question": "test question?",
	}
	_ = d.CreateSession("test_session", "word")
	req := newPost(t, fmt.Sprintf("/session/%s/questions", "test_session"), reqB)
	rr := httptest.NewRecorder()
	service.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("Invalid API response. Code not 200 - %v", rr.Code)
	}
}

func Test_CreateQuestionWhereAnswered(t *testing.T) {
	d := dao.NewMemoryDAO()
	service := NewService(d).Routes()
	reqB := map[string]string{
		"question": "test question?",
	}
	_ = d.CreateSession("test_session", "word")
	_ = d.UpdateSession("test_session", true)
	req := newPost(t, fmt.Sprintf("/session/%s/questions", "test_session"), reqB)
	rr := httptest.NewRecorder()
	service.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Fatalf("Invalid API response. Code not 400 - %v", rr.Code)
	}
}

func Test_CorrectGuess(t *testing.T) {
	d := dao.NewMemoryDAO()
	service := NewService(d).Routes()
	reqB := map[string]string{
		"guess": "word",
	}
	_ = d.CreateSession("test_session", "word")
	req := newPost(t, fmt.Sprintf("/session/%s/questions", "test_session"), reqB)
	rr := httptest.NewRecorder()
	service.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("Invalid API response. Code not 200 - %v", rr.Code)
	}
	respB := map[string]interface{}{}
	if err := json.Unmarshal(rr.Body.Bytes(), &respB); err != nil {
		t.Fatal(err)
	}
	if result, ok := respB["correct"]; !ok {
		t.Fatal("correct not returned")
	} else {
		if asBool, ok := result.(bool); !ok {
			t.Fatal("not a bool")
		} else {
			if asBool != true {
				t.Fatal("incorrect response")
			}
		}
	}
}

func Test_IncorrectGuess(t *testing.T) {
	d := dao.NewMemoryDAO()
	service := NewService(d).Routes()
	reqB := map[string]string{
		"guess": "wrong word",
	}
	_ = d.CreateSession("test_session", "word")
	req := newPost(t, fmt.Sprintf("/session/%s/questions", "test_session"), reqB)
	rr := httptest.NewRecorder()
	service.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("Invalid API response. Code not 200 - %v", rr.Code)
	}
	respB := map[string]interface{}{}
	if err := json.Unmarshal(rr.Body.Bytes(), &respB); err != nil {
		t.Fatal(err)
	}
	if result, ok := respB["correct"]; !ok {
		t.Fatal("correct not returned")
	} else {
		if asBool, ok := result.(bool); !ok {
			t.Fatal("not a bool")
		} else {
			if asBool != false {
				t.Fatal("incorrect response")
			}
		}
	}
}

func Test_AnswerQuestion(t *testing.T) {
	d := dao.NewMemoryDAO()
	service := NewService(d).Routes()
	_ = d.CreateSession("test_session", "word")
	id, _ := d.SaveQuestion("test_session", "what is the question?")

	reqB := map[string]bool{
		"answer": true,
	}
	req := newPost(t, fmt.Sprintf("/session/%s/questions/%v", "test_session", id), reqB)
	rr := httptest.NewRecorder()
	service.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Log(rr.Body.String())
		t.Fatalf("Invalid API response. Code not 200 - %v", rr.Code)
	}
}
