package store

import (
	"testing"
	"time"

	"github.com/SuleymanyanArkadi/eventhub/internal/task"
	"github.com/google/uuid"
)

func TestMemoryStore_CreateGetListUpdate(t *testing.T) {
	s := NewMemoryStore()

	id := uuid.New().String()
	now := time.Now().UTC()
	t1 := &task.Task{
		ID:        id,
		Type:      "email",
		Payload:   `{"to":"a@example.com"}`,
		Status:    task.StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create
	if err := s.Create(t1); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	// Get
	got, err := s.Get(id)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if got.ID != id || got.Type != t1.Type {
		t.Fatalf("Get returned wrong task: %+v", got)
	}

	// List
	list, err := s.List(0, 10)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected list len 1; got %d", len(list))
	}

	// Update
	got.Status = task.StatusRunning
	if err := s.Update(got); err != nil {
		t.Fatalf("Update error: %v", err)
	}
	updated, _ := s.Get(id)
	if updated.Status != task.StatusRunning {
		t.Fatalf("expected status running; got %s", updated.Status)
	}
}
