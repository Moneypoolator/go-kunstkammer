package api

import (
	"encoding/json"
	"fmt"
	"kunstkammer/internal/models"
	"net/http"
	"net/url"
)

// type User struct {
// 	ID        int    `json:"id"`
// 	FullName  string `json:"full_name"`
// 	Email     string `json:"email"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// }

// func PrintUser(user User) {
// 	fmt.Println("User Details:")
// 	fmt.Printf("ID: %d\n", user.ID)
// 	fmt.Printf("FullName: %s\n", user.FullName)
// 	fmt.Printf("Email: %s\n", user.Email)
// 	fmt.Printf("First Name: %s\n", user.FirstName)
// 	fmt.Printf("Last Name: %s\n", user.LastName)
// }

// // FindUserByEmail ищет пользователя по email в массиве
// func FindUserByEmail(users []User, email string) (*User, error) {
// 	for _, user := range users {
// 		if user.Email == email {
// 			return &user, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("user with email '%s' not found", email)
// }

func (kc *KaitenClient) GetUsers() ([]models.User, error) {
	resp, err := kc.doRequestWithBody("GET", "/users", nil)
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := kc.decodeResponse(resp, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (kc *KaitenClient) GetCurrentUser() (*models.User, error) {
	resp, err := kc.doRequest("GET", "/users/current", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUser возвращает информацию о конкретном пользователе по ID
func (kc *KaitenClient) GetUser(userID int) (*models.User, error) {
	resp, err := kc.doRequest("GET", fmt.Sprintf("/users/%d", userID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail возвращает пользователя по email
func (kc *KaitenClient) getUserByEmail(email string) (*models.User, error) {
	resp, err := kc.doRequestWithBody("GET", "/users?email="+url.QueryEscape(email), nil)
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := kc.decodeResponse(resp, &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	return models.FindUserByEmail(users, email)
}
func (kc *KaitenClient) GetUserIDByEmail(responsibleEmail string) (int, error) {
	user, err := kc.getUserByEmail(responsibleEmail)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (kc *KaitenClient) createUser(user *models.User) (*models.User, error) {
	resp, err := kc.doRequestWithBody("POST", "/users", user)
	if err != nil {
		return nil, err
	}

	var newUser models.User
	if err := kc.decodeResponse(resp, &newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (kc *KaitenClient) updateUser(userID int, user *models.User) (*models.User, error) {
	resp, err := kc.doRequestWithBody("PUT", fmt.Sprintf("/users/%d", userID), user)
	if err != nil {
		return nil, err
	}

	var updatedUser models.User
	if err := kc.decodeResponse(resp, &updatedUser); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// CreateUser создает нового пользователя
func (kc *KaitenClient) CreateUser(user *models.User) (*models.User, error) {
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

	var newUser models.User
	if err := json.NewDecoder(resp.Body).Decode(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}

// UpdateUser обновляет информацию о пользователе
func (kc *KaitenClient) UpdateUser(userID int, user *models.User) (*models.User, error) {
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

	var updatedUser models.User
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
