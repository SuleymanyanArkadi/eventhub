package store

import (
	"errors"
	"sync"
	"time"

	"github.com/SuleymanyanArkadi/eventhub/internal/task"
)

var (
	ErrNotFound = errors.New("not found")
)

type MemoryStore struct {
	mu    sync.RWMutex
	tasks map[string]*task.Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make(map[string]*task.Task),
	}
}

func (s *MemoryStore) Create(t *task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[t.ID] = t
	return nil
}

func (s *MemoryStore) Get(id string) (*task.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if t, ok := s.tasks[id]; ok {
		// возвращаем копию чтобы внешняя модификация не ломала хранилище
		cpy := *t
		return &cpy, nil
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) List(offset, limit int) ([]*task.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*task.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		cpy := *t
		res = append(res, &cpy)
	}
	// Простая пагинация (порядок не стабилен — но для примера нормально)
	start := offset
	if start > len(res) {
		start = len(res)
	}
	end := start + limit
	if end > len(res) {
		end = len(res)
	}
	return res[start:end], nil
}

func (s *MemoryStore) Update(t *task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[t.ID]; !ok {
		return ErrNotFound
	}
	// заменяем запись
	t.UpdatedAt = time.Now().UTC()
	s.tasks[t.ID] = t
	return nil
}
