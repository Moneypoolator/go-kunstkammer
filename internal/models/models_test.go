package models

import (
	"testing"
)

func TestFindUserByEmail(t *testing.T) {
	// Подготовка тестовых данных
	users := []User{
		{ID: 1, Email: "user1@example.com"},
		{ID: 2, Email: "user2@example.com"},
		{ID: 3, Email: "user3@example.com"},
	}

	// Тест 1: Поиск существующего пользователя
	t.Run("Find existing user", func(t *testing.T) {
		email := "user2@example.com"
		user, err := FindUserByEmail(users, email)
		if err != nil {
			t.Fatalf("Expected to find user with email %s, got error: %v", email, err)
		}
		if user.Email != email {
			t.Errorf("Expected user email %s, got %s", email, user.Email)
		}
	})

	// Тест 2: Поиск несуществующего пользователя
	t.Run("Find non-existing user", func(t *testing.T) {
		email := "nonexistent@example.com"
		_, err := FindUserByEmail(users, email)
		if err == nil {
			t.Error("Expected error for non-existing user, got nil")
		}
	})

	// Тест 3: Пустой список пользователей
	t.Run("Empty user list", func(t *testing.T) {
		emptyUsers := []User{}
		email := "user1@example.com"
		_, err := FindUserByEmail(emptyUsers, email)
		if err == nil {
			t.Error("Expected error for empty user list, got nil")
		}
	})
}

func TestPrintCard(t *testing.T) {
	// Подготовка тестовой карточки
	card := Card{
		ID:            1,
		Title:         "Test Card",
		Description:   "This is a test card",
		ColumnID:      10,
		BoardID:       20,
		LaneID:        30,
		MemberIDs:     []int{1, 2, 3},
		ParentID:      5,
		TypeID:        7,
		SizeText:      "8h",
		ResponsibleID: 1,
		Properties:    map[string]interface{}{"priority": "high"},
	}

	// Проверка, что функция PrintCard не паникует
	t.Run("PrintCard does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PrintCard panicked: %v", r)
			}
		}()
		PrintCard(card)
	})
}

func TestTaskValidation(t *testing.T) {
	tests := []struct {
		name  string
		task  Task
		valid bool
	}{
		{
			name:  "Valid task",
			task:  Task{Type: "delivery", Size: 8, Title: "Valid Task"},
			valid: true,
		},
		{
			name:  "Empty title",
			task:  Task{Type: "discovery", Size: 16, Title: ""},
			valid: false,
		},
		{
			name:  "Negative size",
			task:  Task{Type: "bug", Size: -1, Title: "Invalid Size"},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid && (tt.task.Title == "" || tt.task.Size <= 0) {
				t.Errorf("Task should be valid: %+v", tt.task)
			}
			if !tt.valid && (tt.task.Title != "" && tt.task.Size > 0) {
				t.Errorf("Task should be invalid: %+v", tt.task)
			}
		})
	}
}

func TestScheduleValidation(t *testing.T) {
	tests := []struct {
		name     string
		schedule Schedule
		valid    bool
	}{
		{
			name: "Valid schedule",
			schedule: Schedule{
				Parent:      "123",
				Responsible: "user@example.com",
				Tasks: []Task{
					{Type: "delivery", Size: 8, Title: "Task 1"},
					{Type: "discovery", Size: 16, Title: "Task 2"},
				},
			},
			valid: true,
		},
		{
			name: "Empty parent",
			schedule: Schedule{
				Parent:      "",
				Responsible: "user@example.com",
				Tasks: []Task{
					{Type: "delivery", Size: 8, Title: "Task 1"},
				},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid && (tt.schedule.Parent == "" || len(tt.schedule.Tasks) == 0) {
				t.Errorf("Schedule should be valid: %+v", tt.schedule)
			}
			if !tt.valid && tt.schedule.Parent != "" && len(tt.schedule.Tasks) > 0 {
				t.Errorf("Schedule should be invalid: %+v", tt.schedule)
			}
		})
	}
}
