package dao

// MemoryDAO is an in memory version of the DAO
type MemoryDAO struct {
	sessions map[string]*Session
}

// NewMemoryDAO creates a new MemoryDAO
func NewMemoryDAO() *MemoryDAO {
	return &MemoryDAO{
		sessions: make(map[string]*Session),
	}
}

// CreateSession creates a session and sets up the question store
func (m *MemoryDAO) CreateSession(id string, word string) error {
	m.sessions[id] = &Session{
		ID:   id,
		Word: word,
	}
	return nil
}

// UpdateSession updates a session
func (m *MemoryDAO) UpdateSession(id string, answered bool) error {
	sess, ok := m.sessions[id]
	if !ok {
		return ErrSessionNotFound
	}
	sess.Answered = true
	return nil
}

// SaveQuestion saves a question and increments the question counter
func (m *MemoryDAO) SaveQuestion(sessionID string, q string) (int, error) {
	if sess, ok := m.sessions[sessionID]; ok {
		if len(sess.Questions) == 20 {
			return 0, ErrNoMoreQuestions
		}
		id := len(sess.Questions)
		sess.Questions = append(sess.Questions, &Question{
			ID:        id,
			SessionID: sessionID,
			Question:  q,
		})
		return id, nil
	}
	return 0, ErrSessionNotFound
}

// GetSession will retrieve a question
func (m *MemoryDAO) GetSession(id string) (*Session, error) {
	if sess, ok := m.sessions[id]; ok {
		return sess, nil
	}
	return nil, ErrSessionNotFound
}

// SaveAnswer saves the answer for a given question
func (m *MemoryDAO) SaveAnswer(sessionID string, questionIndex int, answer bool) error {
	if sess, ok := m.sessions[sessionID]; ok {
		if questionIndex > len(m.sessions) {
			return ErrQuestionNotFound
		}
		sess.Questions[questionIndex].Answer = &answer
		return nil
	}
	return ErrSessionNotFound
}
