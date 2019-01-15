package dao_test

import (
	"github.com/mnbbrown/20questions/dao"
	"testing"
)

func TestCreateSession(t *testing.T) {
	d := dao.NewMemoryDAO()
	err := d.CreateSession("id", "word")
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.GetSession("id")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSessionNotFound(t *testing.T) {
	d := dao.NewMemoryDAO()
	_, err := d.GetSession("not_found")
	if err != dao.ErrSessionNotFound {
		t.Fatal("Error should be not found")
	}
}

func TestAddQuestion(t *testing.T) {
	d := dao.NewMemoryDAO()
	err := d.CreateSession("id", "word")
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.SaveQuestion("id", "What is the question?")
	if err != nil {
		t.Fatal(err)
	}
	session, err := d.GetSession("id")
	if len(session.Questions) != 1 {
		t.Fatal("Questions did not increment")
	}
}

func TestAddQuestion_MoreThan20(t *testing.T) {
	d := dao.NewMemoryDAO()
	err := d.CreateSession("id", "word")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		_, err := d.SaveQuestion("id", "What is the question?")
		if err != nil {
			t.Fatal(err)
		}
	}
	_, err = d.SaveQuestion("id", "What is the question?")
	if err != dao.ErrNoMoreQuestions {
		t.Fatal("Accepted more than 20 answers")
	}
}
