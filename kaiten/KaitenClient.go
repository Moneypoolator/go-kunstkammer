package kaiten

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL = "https://kaiten.nordev.ru/api/latest"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	
}

type Card struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ColumnID    int    `json:"column_id"`
	BoardID     int    `json:"board_id"`
	
}


type KaitenClient struct {
	client  *http.Client
	baseURL string
	token   string
}

func NewKaitenClient(token string) *KaitenClient {
	return &KaitenClient{
		client:  &http.Client{},
		baseURL: baseURL,
		token:   token,
	}
}

func (kc *KaitenClient) doRequest(method, path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, kc.baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+kc.token)
	req.Header.Set("Content-Type", "application/json")

	return kc.client.Do(req)
}

// GetUsers возвращает список всех пользователей
func (kc *KaitenClient) GetUsers() ([]User, error) {
	resp, err := kc.doRequest("GET", "/users", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUser возвращает информацию о конкретном пользователе по ID
func (kc *KaitenClient) GetUser(userID int) (*User, error) {
	resp, err := kc.doRequest("GET", fmt.Sprintf("/users/%d", userID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser создает нового пользователя
func (kc *KaitenClient) CreateUser(user *User) (*User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	resp, err := kc.doRequest("POST", "/users", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var newUser User
	if err := json.NewDecoder(resp.Body).Decode(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

// UpdateUser обновляет информацию о пользователе
func (kc *KaitenClient) UpdateUser(userID int, user *User) (*User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	resp, err := kc.doRequest("PUT", fmt.Sprintf("/users/%d", userID), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var updatedUser User
	if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// DeleteUser удаляет пользователя по ID
func (kc *KaitenClient) DeleteUser(userID int) error {
	resp, err := kc.doRequest("DELETE", fmt.Sprintf("/users/%d", userID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}


// GetCards возвращает список всех карт
func (kc *KaitenClient) GetCards() ([]Card, error) {
	resp, err := kc.doRequest("GET", "/cards", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cards []Card
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCard возвращает информацию о конкретной карте по ID
func (kc *KaitenClient) GetCard(cardID int) (*Card, error) {
	resp, err := kc.doRequest("GET", fmt.Sprintf("/cards/%d", cardID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var card Card
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		return nil, err
	}

	return &card, nil
}

// CreateCard создает новую карту
func (kc *KaitenClient) CreateCard(card *Card) (*Card, error) {
	body, err := json.Marshal(card)
	if err != nil {
		return nil, err
	}

	resp, err := kc.doRequest("POST", "/cards", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var newCard Card
	if err := json.NewDecoder(resp.Body).Decode(&newCard); err != nil {
		return nil, err
	}

	return &newCard, nil
}

// UpdateCard обновляет информацию о карте
func (kc *KaitenClient) UpdateCard(cardID int, card *Card) (*Card, error) {
	body, err := json.Marshal(card)
	if err != nil {
		return nil, err
	}

	resp, err := kc.doRequest("PUT", fmt.Sprintf("/cards/%d", cardID), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var updatedCard Card
	if err := json.NewDecoder(resp.Body).Decode(&updatedCard); err != nil {
		return nil, err
	}

	return &updatedCard, nil
}

// DeleteCard удаляет карту по ID
func (kc *KaitenClient) DeleteCard(cardID int) error {
	resp, err := kc.doRequest("DELETE", fmt.Sprintf("/cards/%d", cardID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
