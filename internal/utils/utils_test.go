package utils

import (
	//"kunstkammer/internal/utils"
	"path/filepath"
	"testing"
)

func TestLoadTasksFromJSON(t *testing.T) {
	// Путь к тестовому файлу
	testFilePath := filepath.Join("testdata", "test_tasks.json")

	// Загружаем задачи из JSON-файла
	schedule, err := LoadTasksFromJSON(testFilePath)
	if err != nil {
		t.Fatalf("Failed to load tasks from JSON: %v", err)
	}

	// Проверяем, что данные загружены корректно
	if schedule.Parent != "123" {
		t.Errorf("Expected parent ID '123', got '%s'", schedule.Parent)
	}

	if schedule.Responsible != "user@example.com" {
		t.Errorf("Expected responsible email 'user@example.com', got '%s'", schedule.Responsible)
	}

	if len(schedule.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(schedule.Tasks))
	}

	// Проверяем первую задачу
	task1 := schedule.Tasks[0]
	if task1.Type == nil || *task1.Type != "delivery" {
		t.Errorf("Expected task type 'delivery', got '%v'", task1.Type)
	}
	if task1.Size != 8 {
		t.Errorf("Expected task size 8, got %d", task1.Size)
	}
	if task1.Title != "Task 1" {
		t.Errorf("Expected task title 'Task 1', got '%s'", task1.Title)
	}

	// Проверяем вторую задачу
	task2 := schedule.Tasks[1]
	if task2.Type == nil || *task2.Type != "discovery" {
		t.Errorf("Expected task type 'discovery', got '%v'", task2.Type)
	}
	if task2.Size != 16 {
		t.Errorf("Expected task size 16, got %d", task2.Size)
	}
	if task2.Title != "Task 2" {
		t.Errorf("Expected task title 'Task 2', got '%s'", task2.Title)
	}
}

func TestLoadTasksFromJSON_FileNotFound(t *testing.T) {
	// Путь к несуществующему файлу
	nonExistentFilePath := filepath.Join("testdata", "non_existent.json")

	// Пытаемся загрузить задачи из несуществующего файла
	_, err := LoadTasksFromJSON(nonExistentFilePath)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestExtractWorkCode(t *testing.T) {
	tests := []struct {
		name        string
		parentTitle string
		expProduct  string
		expWorkCode string
		expectError bool
	}{
		{
			name:        "Valid CAD title",
			parentTitle: "[CAD]:TS.FEATURE.123. Task Title",
			expProduct:  "CAD",
			expWorkCode: "TS.FEATURE.123",
			expectError: false,
		},
		{
			name:        "Valid MGM title",
			parentTitle: "[MGM]:Widget.789.101 Some Task",
			expProduct:  "MGM",
			expWorkCode: "Widget.789.101",
			expectError: false,
		},
		{
			name:        "Invalid title - no match",
			parentTitle: "Invalid Title",
			expProduct:  "",
			expWorkCode: "",
			expectError: true,
		},
		{
			name:        "Empty title",
			parentTitle: "",
			expProduct:  "",
			expWorkCode: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, workCode, err := ExtractWorkCode(tt.parentTitle)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if product != tt.expProduct {
					t.Errorf("Expected product '%s', got '%s'", tt.expProduct, product)
				}
				if workCode != tt.expWorkCode {
					t.Errorf("Expected work code '%s', got '%s'", tt.expWorkCode, workCode)
				}
			}
		})
	}
}
