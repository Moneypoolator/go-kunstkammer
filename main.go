package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	baseURL = "https://kaiten.norsoft.ru/api/latest"
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

// FindUserByEmail ищет пользователя по email в массиве
func FindUserByEmail(users []User, email string) (*User, error) {
	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with email '%s' not found", email)
}

// TaskType описывает тип задачи в Kaiten
type TaskType struct {
	ID   int    `json:"id"`   // ID типа задачи
	Name string `json:"name"` // Название типа задачи
	// Добавьте другие поля, если они есть в API
}

// TaskType описывает тип задачи в Kaiten
type TaskIDType int

// Константы для типов задач
const (
	BugTaskType                 TaskIDType = 7
	UserStoryTaskType           TaskIDType = 5
	TaskDiscoveryTaskType       TaskIDType = 11
	CardTaskType                TaskIDType = 1
	FeatureTaskType             TaskIDType = 4
	ACCTaskType                 TaskIDType = 18
	ImprovementTaskType         TaskIDType = 15
	VulnerabilityTaskType       TaskIDType = 13
	IntegrationTasksTaskType    TaskIDType = 19
	TechDebtTaskType            TaskIDType = 9
	AdministrativeTasksTaskType TaskIDType = 12
	TaskDeliveryTaskType        TaskIDType = 6
	RequestTaskType             TaskIDType = 17
	EnablerTaskType             TaskIDType = 8
)

// String возвращает строковое представление типа задачи
func (tt TaskIDType) String() string {
	switch tt {
	case BugTaskType:
		return "Bug"
	case UserStoryTaskType:
		return "User Story"
	case TaskDiscoveryTaskType:
		return "Task Discovery"
	case CardTaskType:
		return "Card"
	case FeatureTaskType:
		return "Feature"
	case ACCTaskType:
		return "ACC"
	case ImprovementTaskType:
		return "Улучшение"
	case VulnerabilityTaskType:
		return "Уязвимость"
	case IntegrationTasksTaskType:
		return "Интеграционные задачи"
	case TechDebtTaskType:
		return "Техдолг"
	case AdministrativeTasksTaskType:
		return "Административные задачи"
	case TaskDeliveryTaskType:
		return "Task Delivery"
	case RequestTaskType:
		return "Обращение"
	case EnablerTaskType:
		return "Enabler"
	default:
		return "Unknown"
	}
}

// taskTypeMap сопоставляет имя типа задачи с его константой
var taskTypeMap = map[string]TaskIDType{
	"Bug":            BugTaskType,
	"User Story":     UserStoryTaskType,
	"Task Discovery": TaskDiscoveryTaskType,
	"discovery":      TaskDiscoveryTaskType,
	"Card":           CardTaskType,
	"Feature":        FeatureTaskType,
	"ACC":            ACCTaskType,
	"Улучшение":      ImprovementTaskType,
	"Уязвимость":     VulnerabilityTaskType,
	"Интеграционные задачи": IntegrationTasksTaskType,
	"Техдолг": TechDebtTaskType,
	"Административные задачи": AdministrativeTasksTaskType,
	"Task Delivery": TaskDeliveryTaskType,
	"delivery":      TaskDeliveryTaskType,
	"Обращение":     RequestTaskType,
	"Enabler":       EnablerTaskType,
}

// GetTaskTypeByName возвращает константу типа задачи по его имени
func GetTaskTypeByName(name string) (TaskIDType, error) {
	if taskType, exists := taskTypeMap[name]; exists {
		return taskType, nil
	}
	return 0, fmt.Errorf("task type '%s' not found", name)
}

type Card struct {
	ID            int                    `json:"id"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	ColumnID      int                    `json:"column_id"`
	BoardID       int                    `json:"board_id"`
	LaneID        int                    `json:"lane_id"`              // Lane ID
	MemberIDs     []int                  `json:"member_ids"`           // Поле для идентификаторов участников карточки
	ParentID      int                    `json:"parent_id,omitempty"`  // ID родительской карточки (если есть)
	TypeID        int                    `json:"type_id"`              // ID типа задачи
	SizeText      string                 `json:"size_text"`            // Размер задачи
	ResponsibleID int                    `json:"responsible_id"`       // Responsible ID
	Properties    map[string]interface{} `json:"properties,omitempty"` // Пользовательские свойства
}

func PrintCard(card Card) {
	fmt.Println("Card Details:")
	fmt.Printf("ID: %d\n", card.ID)
	fmt.Printf("Title: %s\n", card.Title)
	fmt.Printf("Description: %s\n", card.Description)
	fmt.Printf("ColumnID: %d\n", card.ColumnID)
	fmt.Printf("BoardID: %d\n", card.BoardID)
	fmt.Println("MemberIDs: ", card.MemberIDs)
	fmt.Printf("ParentID: %d\n", card.ParentID)
	fmt.Printf("TypeID: %d\n", card.TypeID)
	fmt.Printf("SizeText: %s\n", card.SizeText)
	fmt.Printf("ResponsibleID: %d\n", card.ResponsibleID)
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

	user, err := FindUserByEmail(users, email)
	if user == nil || err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return user, nil
}

func (kc *KaitenClient) GetUserIDByEmail(responsibleEmail string) (int, error) {
	user, err := kc.GetUserByEmail(responsibleEmail)
	if user == nil || err != nil {
		return 0, err
	}

	return user.ID, nil
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

// GetTaskTypes возвращает список типов задач из Kaiten
func (kc *KaitenClient) GetTaskTypes() ([]TaskType, error) {
	// Выполняем GET-запрос к эндпоинту /api/latest/card-types
	resp, err := kc.doRequest("GET", "/card-types", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch task types: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Декодируем ответ в список типов задач
	var taskTypes []TaskType
	if err := json.NewDecoder(resp.Body).Decode(&taskTypes); err != nil {
		return nil, fmt.Errorf("failed to decode task types: %w", err)
	}

	return taskTypes, nil
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

	if resp.StatusCode != http.StatusOK { // || resp.StatusCode != http.StatusCreated
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var createdCard Card
	if err := json.NewDecoder(resp.Body).Decode(&createdCard); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdCard, nil
}

// CardUpdate описывает данные для обновления карточки
type CardUpdate struct {
	Title         *string                `json:"title,omitempty"`
	Description   *string                `json:"description,omitempty"`
	ColumnID      *int                   `json:"column_id,omitempty"`
	BoardID       *int                   `json:"board_id,omitempty"`
	LaneID        *int                   `json:"lane_id,omitempty"`
	MemberIDs     *[]int                 `json:"member_ids,omitempty"`
	ParentID      *int                   `json:"parent_id,omitempty"`
	TypeID        *int                   `json:"type_id,omitempty"`
	SizeText      *string                `json:"size_text,omitempty"`
	ResponsibleID *int                   `json:"responsible_id,omitempty"`
	OwnerID       *int                   `json:"owner_id,omitempty"`
	OwnerEmailID  *string                `json:"owner_email,omitempty"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}

// UpdateCard обновляет карточку с указанным ID
func (kc *KaitenClient) UpdateCard(cardID int, update CardUpdate) error {
	// Преобразуем данные обновления в JSON
	body, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal update data: %w", err)
	}

	// Создаем PATCH-запрос
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/cards/%d", kc.baseURL, cardID), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+kc.token)
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := kc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
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

// TagRequest описывает данные для добавления тега
type TagRequest struct {
	Name string `json:"name"` // Имя тега
}

// AddTagToCard добавляет тег к карточке с указанным ID
func (kc *KaitenClient) AddTagToCard(cardID int, tagName string) error {
	// Создаем данные для добавления тега
	tagRequest := TagRequest{
		Name: tagName,
	}

	// Преобразуем данные в JSON
	body, err := json.Marshal(tagRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal tag data: %w", err)
	}

	// Создаем POST-запрос
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cards/%d/tags", kc.baseURL, cardID), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Authorization", "Bearer "+kc.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := kc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Add children Request
type AddChildrenRequest struct {
	CardID int `json:"card_id"`
}

// AddTagToCard добавляет тег к карточке с указанным ID
func (kc *KaitenClient) AddChindrenToCard(cardID int, childrenCardID int) error {
	// Создаем данные для добавления тега
	requestData := AddChildrenRequest{
		CardID: childrenCardID,
	}

	// Преобразуем данные в JSON
	body, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to marshal tag data: %w", err)
	}

	// Создаем POST-запрос
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cards/%d/children", kc.baseURL, cardID), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Authorization", "Bearer "+kc.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := kc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Task описывает задачу для создания карточки
type Task struct {
	Type  string `json:"type"`  // Тип задачи (например, "delivery", "discovery")
	Size  int    `json:"size"`  // Размер задачи (например, 8, 16)
	Title string `json:"title"` // Название задачи
}

// Schedule описывает расписание задач
type Schedule struct {
	Parent      string `json:"parent"`      // ID родительской карточки
	Responsible string `json:"responsible"` // Email ответственного
	Tasks       []Task `json:"tasks"`       // Список задач
}

// ScheduleFile описывает структуру JSON-файла
type ScheduleFile struct {
	Schedule Schedule `json:"schedule"`
}

// LoadTasksFromJSON загружает задачи из JSON-файла и возвращает описание для создания карточек
func LoadTasksFromJSON(filePath string) (*Schedule, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var scheduleFile ScheduleFile
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&scheduleFile); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &scheduleFile.Schedule, nil
}

func main() {

	// Определение флагов командной строки
	tasksFile := flag.String("tasks", "", "Path to the tasks JSON file (required)")
	configFile := flag.String("config", "config.json", "Path to the configuration file (optional, default: config.json)")

	// Парсинг аргументов командной строки
	flag.Parse()

	// Проверка обязательного аргумента
	if *tasksFile == "" {
		fmt.Println("Error: The 'tasks' flag is required.")
		flag.Usage() // Вывод справки по использованию
		//os.Exit(1)
		*tasksFile = "tasks.json"
	}

	// Загрузка конфигурации (если указан файл конфигурации)
	if *configFile != "" {
		fmt.Printf("Loading configuration from: %s\n", *configFile)
		// Здесь можно добавить логику для загрузки конфигурации
	} else {
		fmt.Println("Using default configuration.")
	}

	client := CreateKaitenClient("9ecb4b54-508a-4d1e-ad99-c1c4a04847bb")

	// Создаем кастомный HTTP-клиент с отключенной проверкой сертификата
	// client = &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{
	// 			InsecureSkipVerify: true, // Отключаем проверку сертификата
	// 		},
	// 	},
	// }

	// // Загрузите сертификат CA
	// caCert, err := ioutil.ReadFile("path/to/your/ca.crt")
	// if err != nil {
	// 	log.Fatalf("Failed to read CA certificate: %v", err)
	// }

	// // Создайте пул сертификатов и добавьте CA
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	// // Создайте кастомный TLS-конфиг
	// tlsConfig := &tls.Config{
	// 	RootCAs: caCertPool, // Используем кастомный пул CA
	// }

	// // Создайте кастомный HTTP-клиент
	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: tlsConfig,
	// 	},
	// }

	// // Получение списка пользователей
	// users, err := client.GetUsers()
	// if err != nil {
	// 	fmt.Println("Error getting users:", err)
	// 	return
	// }
	// fmt.Println("Users:", users)

	// Получение данных текущего пользователя
	currentUser, err := client.GetCurrentUser()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	// // Печать данных текущего пользователя
	// PrintUser(*currentUser)
	// fmt.Println("-----------------------------------------")

	// // Получение списка карточек для пользователя с ID
	// userID := currentUser.ID
	// cards, err := client.GetUserCards(userID)
	// if err != nil {
	// 	fmt.Println("Error getting user cards:", err)
	// 	return
	// }

	// PrintCardsList(cards, userID)
	// fmt.Println("-----------------------------------------")

	// limit := 100
	// offset := 0

	// fmt.Printf("Get card by pages: limit=%d, offset=%d\n", limit, offset)

	// cardsByMemberIDs, err := client.GetUserCardsByMemberIDs(userID, limit, offset)
	// if err != nil {
	// 	fmt.Println("Error getting user cards:", err)
	// 	return
	// }

	// PrintCardsList(cardsByMemberIDs, userID)
	// fmt.Println("-----------------------------------------")

	// allUserCards, err := client.GetAllUserCards(userID)
	// if err != nil {
	// 	fmt.Println("Error getting all user cards:", err)
	// 	return
	// }

	// // Печать списка всех карточек
	// fmt.Printf("All cards for user %d:\n", 11)
	// fmt.Printf("Cards count=%d :\n", len(allUserCards))
	// fmt.Println("Cards:", allUserCards)
	// fmt.Println("-----------------------------------------")

	// Загрузка задач из JSON-файла
	schedule, err := LoadTasksFromJSON(*tasksFile)
	if err != nil {
		fmt.Println("Error loading tasks from JSON file:", err)
		return
	}

	// Вывод задач
	fmt.Printf("Parent: %s\n", schedule.Parent)
	fmt.Printf("Responsible: %s\n", schedule.Responsible)

	// Преобразуем email ответственного в ID пользователя
	responsibleID, err := client.GetUserIDByEmail(schedule.Responsible)
	if responsibleID == 0 || err != nil {
		fmt.Println("Error getting responsible user ID:", err)
		responsibleID = currentUser.ID
		// return
	}

	// Преобразуем ID родительской карточки из строки в число
	parentID, err := strconv.Atoi(schedule.Parent)
	if err != nil {
		fmt.Println("Error parsing parent ID:", err)
		return
	}

	// Создаем карточки для каждой задачи
	for _, task := range schedule.Tasks {
		// Определяем тип задачи
		taskTypeID, err := GetTaskTypeByName(task.Type)
		if err != nil {
			fmt.Printf("Error getting task type ID for '%s': %v\n", task.Type, err)
			continue
		}

		// Создаем карточку
		card := &Card{
			ID:            0,
			Title:         task.Title,
			BoardID:       192,
			ColumnID:      776,
			LaneID:        1275,
			TypeID:        int(taskTypeID),
			SizeText:      fmt.Sprintf("%d ч", task.Size),
			ParentID:      parentID,
			MemberIDs:     []int{responsibleID},
			ResponsibleID: responsibleID,

			Properties: map[string]interface{}{
				"id_19": "1", // Строка
			},
		}

		PrintCard(*card)

		// Создаем карточку в Kaiten
		createdCard, err := client.CreateCard(card)
		if err != nil {
			fmt.Printf("Error creating card '%s': %v\n", task.Title, err)
			continue
		} else {
			fmt.Printf("Created card: %s (ID: %d)\n", createdCard.Title, createdCard.ID)
		}

		if createdCard.TypeID == int(TaskDeliveryTaskType) || createdCard.TypeID == int(TaskDiscoveryTaskType) {

			titleUpdate := fmt.Sprintf("[CAD]:TS.%s.%d. %s", "XX.XX", createdCard.ID, createdCard.Title)
			updateData := &CardUpdate{
				Title: stringPtr(titleUpdate),
				// BoardID:      intPtr(192),
				// ColumnID:     intPtr(776),
				// LaneID:       intPtr(1275),
				// TypeID:       intPtr(int(createdCard.TypeID)),
				// OwnerID:      intPtr(responsibleID),
				// OwnerEmailID: stringPtr(schedule.Responsible),
				// Properties: map[string]interface{}{
				// 	"id_19": "1", // Строка
				// },
			}

			err = client.UpdateCard(createdCard.ID, *updateData)
			if err != nil {
				fmt.Printf("Error updating card '%s': %v\n", titleUpdate, err)
			} else {
				fmt.Printf("Updated card: %s (ID: %d)\n", titleUpdate, createdCard.ID)
			}
		}

		err = client.AddChindrenToCard(parentID, createdCard.ID)
		if err != nil {
			fmt.Println("Error adding children to card:", err)
			//return
		}

		// Добавляем тег к карточке
		err = client.AddTagToCard(createdCard.ID, "ГГИС")
		if err != nil {
			fmt.Println("Error adding tag to card:", err)
			//return
		}

		// Добавляем тег к карточке
		err = client.AddTagToCard(createdCard.ID, "C++")
		if err != nil {
			fmt.Println("Error adding tag to card:", err)
			//return
		}

	}

	// // Получаем список типов задач
	// taskTypes, err := client.GetTaskTypes()
	// if err != nil {
	// 	fmt.Println("Error fetching task types:", err)
	// 	return
	// }

	// // Выводим список типов задач
	// fmt.Println("Task Types:")
	// for _, taskType := range taskTypes {
	// 	fmt.Printf("ID: %d, Name: %s\n", taskType.ID, taskType.Name)
	// }

	// // Пример использования переименованных констант
	// taskType := BugTaskType
	// fmt.Printf("Task Type: %s (ID: %d)\n", taskType, taskType)

	// // Итерация по всем типам задач
	// taskTypesWithIds := []TaskIDType{
	// 	BugTaskType,
	// 	UserStoryTaskType,
	// 	TaskDiscoveryTaskType,
	// 	CardTaskType,
	// 	FeatureTaskType,
	// 	ACCTaskType,
	// 	ImprovementTaskType,
	// 	VulnerabilityTaskType,
	// 	IntegrationTasksTaskType,
	// 	TechDebtTaskType,
	// 	AdministrativeTasksTaskType,
	// 	TaskDeliveryTaskType,
	// 	RequestTaskType,
	// 	EnablerTaskType,
	// }

	// fmt.Println("All Task Types:")
	// for _, tt := range taskTypesWithIds {
	// 	fmt.Printf("ID: %d, Name: %s\n", tt, tt)
	// }

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

// Вспомогательные функции для создания указателей
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
