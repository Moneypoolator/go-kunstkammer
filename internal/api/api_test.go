package api

import (
	"encoding/json"
	"kunstkammer/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsers(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users := []models.User{
			{ID: 1, Email: "user1@example.com"},
			{ID: 2, Email: "user2@example.com"},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client := &KaitenClient{
		client:  ts.Client(),
		baseURL: ts.URL,
		token:   "test-token",
	}

	// Выполняем тест
	users, err := client.GetUsers()
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestGetUserByEmail(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users := []models.User{
			{ID: 1, Email: "user1@example.com"},
			{ID: 2, Email: "user2@example.com"},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client := &KaitenClient{
		client:  ts.Client(),
		baseURL: ts.URL,
		token:   "test-token",
	}

	// Выполняем тест
	user, err := client.GetUserIDByEmail("user1@example.com")
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if user != 1 {
		t.Errorf("Expected user email 'user1@example.com', got '%d'", user)
	}
}

func TestGetCards(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cards := []models.Card{
			{ID: 1, Title: "Card 1"},
			{ID: 2, Title: "Card 2"},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cards)
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client := &KaitenClient{
		client:  ts.Client(),
		baseURL: ts.URL,
		token:   "test-token",
	}

	// Выполняем тест
	cards, err := client.GetCards()
	if err != nil {
		t.Fatalf("Failed to get cards: %v", err)
	}

	if len(cards) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(cards))
	}
}

func TestCreateCard(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var card models.Card
		json.NewDecoder(r.Body).Decode(&card)
		card.ID = 1 // Присваиваем ID для имитации ответа сервера
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(card)
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client := &KaitenClient{
		client:  ts.Client(),
		baseURL: ts.URL,
		token:   "test-token",
	}

	// Выполняем тест
	newCard := &models.Card{Title: "New Card"}
	createdCard, err := client.CreateCard(newCard)
	if err != nil {
		t.Fatalf("Failed to create card: %v", err)
	}

	if createdCard.ID != 1 {
		t.Errorf("Expected card ID 1, got %d", createdCard.ID)
	}
}

func TestGetTaskTypes(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskTypes := []models.TaskType{
			{ID: 1, Name: "Bug"},
			{ID: 2, Name: "Feature"},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(taskTypes)
	}))
	defer ts.Close()

	// Создаем клиент с тестовым сервером
	client := &KaitenClient{
		client:  ts.Client(),
		baseURL: ts.URL,
		token:   "test-token",
	}

	// Выполняем тест
	taskTypes, err := client.GetTaskTypes()
	if err != nil {
		t.Fatalf("Failed to get task types: %v", err)
	}

	if len(taskTypes) != 2 {
		t.Errorf("Expected 2 task types, got %d", len(taskTypes))
	}
}
