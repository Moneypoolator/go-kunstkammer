package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	baseURL = "https://kaiten.nordev.ru/api/latest"
)

type User struct {
	ID        int    `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func PrintUser(user User) {
	fmt.Println("User Details:")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("FullName: %s\n", user.FullName)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("First Name: %s\n", user.FirstName)
	fmt.Printf("Last Name: %s\n", user.LastName)
}

type Card struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ColumnID    int    `json:"column_id"`
	BoardID     int    `json:"board_id"`
	MemberIDs   []int  `json:"member_ids"`          // Поле для идентификаторов участников карточки
	ParentID    int    `json:"parent_id,omitempty"` // ID родительской карточки (если есть)
}

func PrintCard(card Card) {
	fmt.Println("Card Details:")
	fmt.Printf("ID: %d\n", card.ID)
	fmt.Printf("Title: %s\n", card.Title)
	fmt.Printf("Description: %s\n", card.Description)
	fmt.Printf("ColumnID: %d\n", card.ColumnID)
	fmt.Printf("BoardID: %d\n", card.BoardID)
	fmt.Println("MemberIDs: ", card.MemberIDs)
}

type KaitenClient struct {
	client  *http.Client
	baseURL string
	token   string
}

func CreateKaitenClient(token string) *KaitenClient {
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

func (kc *KaitenClient) GetCurrentUser() (*User, error) {
	resp, err := kc.doRequest("GET", "/users/current", nil)
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

// GetUserCards возвращает список карточек для пользователя с указанным идентификатором
func (kc *KaitenClient) GetUserCards(userID int) ([]Card, error) {
	resp, err := kc.doRequest("GET", fmt.Sprintf("/users/%d/cards", userID), nil)
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

// GetUserByEmail возвращает пользователя по email
func (kc *KaitenClient) GetUserByEmail(email string) (*User, error) {
	resp, err := kc.doRequest("GET", "/users?email="+url.QueryEscape(email), nil)
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

	if len(users) == 0 {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	return &users[0], nil
}

// GetUserCards возвращает список карточек, в которых участвует пользователь с указанным идентификатором
// limit - количество карточек на странице
// offset - смещение относительно начала списка
func (kc *KaitenClient) GetUserCardsByMemberIDs(userID int, limit int, offset int) ([]Card, error) {
	// Формируем параметры запроса
	params := url.Values{}
	params.Add("member_ids", fmt.Sprintf("%d", userID)) // Фильтр по пользователю
	params.Add("limit", fmt.Sprintf("%d", limit))       // Количество карточек на странице
	params.Add("offset", fmt.Sprintf("%d", offset))     // Смещение
	params.Add("condition", "1")                        // странное поле, которое должно быть заполнено 1

	// Выполняем запрос к эндпоинту /cards с параметрами
	resp, err := kc.doRequest("GET", "/cards?"+params.Encode(), nil)
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

// GetAllUserCards возвращает все карточки, в которых участвует пользователь с указанным идентификатором
func (kc *KaitenClient) GetAllUserCards(userID int) ([]Card, error) {
	var allCards []Card
	limit := 100 // Количество карточек на странице
	offset := 0  // Начальное смещение

	for {
		// Формируем параметры запроса
		params := url.Values{}
		params.Add("member_ids", fmt.Sprintf("%d", userID)) // Фильтр по пользователю
		params.Add("limit", fmt.Sprintf("%d", limit))       // Количество карточек на странице
		params.Add("offset", fmt.Sprintf("%d", offset))     // Смещение
		params.Add("condition", "1")                        // странное поле, которое должно быть заполнено 1

		// Выполняем запрос к эндпоинту /cards с параметрами
		resp, err := kc.doRequest("GET", "/cards?"+params.Encode(), nil)
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

		// Если карточек нет, завершаем цикл
		if len(cards) == 0 {
			break
		}

		// Добавляем полученные карточки в общий список
		allCards = append(allCards, cards...)

		// Увеличиваем offset для следующей страницы
		offset += limit
	}

	return allCards, nil
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

// ChildTasks описывает задачи для создания внутри родительской карточки
type ChildTasks struct {
	ParentTitle      string   // Название родительской карточки или часть названия
	ResponsibleEmail string   // Email ответственного пользователя (если есть)
	ParentID         int      // ID родительской карточки (если есть)
	TasksTitles      []string // Список задач, которые должны быть созданы внутри родительской карточки
}

// LoadTasksFromFile загружает задачи из текстового файла и возвращает описание для создания карточек
func LoadTasksFromFile(filePath string) ([]ChildTasks, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var tasks []ChildTasks
	var currentParent *ChildTasks

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") {
			// Родительская задача по ID
			parentIDStr := strings.TrimSpace(line[1:])
			parentID, err := strconv.Atoi(parentIDStr)
			if err != nil {
				return nil, fmt.Errorf("invalid parent ID: %s", parentIDStr)
			}

			// Создаем новую родительскую задачу
			currentParent = &ChildTasks{
				ParentID:    parentID,
				TasksTitles: []string{},
			}
			tasks = append(tasks, *currentParent)
		} else if strings.HasPrefix(line, "@") {
			// Ответственный
			email := strings.TrimSpace(line[1:])
			if len(tasks) == 0 {
				return nil, fmt.Errorf("responsible email specified without a parent task")
			}
			// Обновляем последнюю родительскую задачу
			tasks[len(tasks)-1].ResponsibleEmail = email
		} else if strings.HasPrefix(line, "\t") {
			// Название задачи
			taskTitle := strings.TrimSpace(line)
			if taskTitle == "" {
				continue
			}

			if len(tasks) == 0 {
				return nil, fmt.Errorf("task specified without a parent task")
			}
			// Добавляем задачу к последней родительской задаче
			tasks[len(tasks)-1].TasksTitles = append(tasks[len(tasks)-1].TasksTitles, taskTitle)
		} else {
			// Родительская задача по названию
			// if len(tasks) > 0 && tasks[len(tasks)-1].ParentID == 0 && tasks[len(tasks)-1].ParentTitle != "" {
			// Если текущая родительская задача уже существует (по названию), добавляем к её названию
			// 	tasks[len(tasks)-1].ParentTitle = line
			// } else {
			// Создаем новую родительскую задачу по названию
			currentParent = &ChildTasks{
				ParentTitle: line,
				TasksTitles: []string{},
			}
			tasks = append(tasks, *currentParent)
			// }
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return tasks, nil
}

func main() {
	client := CreateKaitenClient("9ecb4b54-508a-4d1e-ad99-c1c4a04847bb")

	// Получение списка пользователей
	users, err := client.GetUsers()
	if err != nil {
		fmt.Println("Error getting users:", err)
		return
	}
	fmt.Println("Users:", users)

	// Получение данных текущего пользователя
	currentUser, err := client.GetCurrentUser()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	// Печать данных текущего пользователя
	PrintUser(*currentUser)
	fmt.Println("-----------------------------------------")

	// Получение списка карточек для пользователя с ID
	userID := currentUser.ID
	cards, err := client.GetUserCards(userID)
	if err != nil {
		fmt.Println("Error getting user cards:", err)
		return
	}

	PrintCardsList(cards, userID)
	fmt.Println("-----------------------------------------")

	limit := 100
	offset := 0

	fmt.Printf("Get card by pages: limit=%d, offset=%d\n", limit, offset)

	cardsByMemberIDs, err := client.GetUserCardsByMemberIDs(userID, limit, offset)
	if err != nil {
		fmt.Println("Error getting user cards:", err)
		return
	}

	PrintCardsList(cardsByMemberIDs, userID)
	fmt.Println("-----------------------------------------")

	allUserCards, err := client.GetAllUserCards(userID)
	if err != nil {
		fmt.Println("Error getting all user cards:", err)
		return
	}

	// Печать списка всех карточек
	fmt.Printf("All cards for user %d:\n", 11)
	fmt.Printf("Cards count=%d :\n", len(allUserCards))
	fmt.Println("Cards:", allUserCards)
	fmt.Println("-----------------------------------------")

	// Загрузка задач из файла
	tasks, err := LoadTasksFromFile("tasks.txt")
	if err != nil {
		fmt.Println("Error loading tasks from file:", err)
		return
	}

	// Вывод задач
	fmt.Println("Tasks loaded:")
	for _, task := range tasks {
		fmt.Printf("Parent: %s (ID: %d, Responsible: %s)\n", task.ParentTitle, task.ParentID, task.ResponsibleEmail)
		for _, title := range task.TasksTitles {
			fmt.Printf("  - %s\n", title)
		}
	}

	// // Создание нового пользователя
	// newUser := &kaiten.User{
	// 	Email:     "newuser@example.com",
	// 	FirstName: "John",
	// 	LastName:  "Doe",
	// }
	// createdUser, err := client.CreateUser(newUser)
	// if err != nil {
	// 	fmt.Println("Error creating user:", err)
	// 	return
	// }
	// fmt.Println("Created user:", createdUser)

	// // Обновление пользователя
	// createdUser.FirstName = "Jane"
	// updatedUser, err := client.UpdateUser(createdUser.ID, createdUser)
	// if err != nil {
	// 	fmt.Println("Error updating user:", err)
	// 	return
	// }
	// fmt.Println("Updated user:", updatedUser)

	// // Удаление пользователя
	// err = client.DeleteUser(updatedUser.ID)
	// if err != nil {
	// 	fmt.Println("Error deleting user:", err)
	// 	return
	// }
	// fmt.Println("User deleted")

	// // Получение списка карт
	// allCards, err := client.GetCards()
	// if err != nil {
	// 	fmt.Println("Error getting cards:", err)
	// 	return
	// }
	// fmt.Printf("Cards count=%d :\n", len(allCards))
	// fmt.Println("Cards:", allCards)

	// // Создание новой карты
	// newCard := &kaiten.Card{
	// 	Title:       "New Card",
	// 	Description: "This is a new card",
	// 	ColumnID:    123,
	// 	BoardID:     456,
	// }
	// createdCard, err := client.CreateCard(newCard)
	// if err != nil {
	// 	fmt.Println("Error creating card:", err)
	// 	return
	// }
	// fmt.Println("Created card:", createdCard)

	// // Обновление карты
	// createdCard.Title = "Updated Card"
	// updatedCard, err := client.UpdateCard(createdCard.ID, createdCard)
	// if err != nil {
	// 	fmt.Println("Error updating card:", err)
	// 	return
	// }
	// fmt.Println("Updated card:", updatedCard)

	// // Удаление карты
	// err = client.DeleteCard(updatedCard.ID)
	// if err != nil {
	// 	fmt.Println("Error deleting card:", err)
	// 	return
	// }
	// fmt.Println("Card deleted")

}

func PrintCardsList(cards []Card, userID int) {
	fmt.Printf("Cards count=%d :\n", len(cards))
	fmt.Printf("Cards for user %d:\n", userID)
	for _, card := range cards {
		fmt.Printf("Card ID: %d, Title: %s, Description: %s, ColumnID: %d, BoardID: %d\n",
			card.ID, card.Title, card.Description, card.ColumnID, card.BoardID)
	}
}
