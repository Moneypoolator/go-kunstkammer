package main

import (
	"flag"
	"fmt"
	"kunstkammer/internal/api"
	"kunstkammer/internal/models"
	"kunstkammer/internal/utils"
	"strconv"
)

// // // const (
// // // 	baseURL = "https://kaiten.norsoft.ru/api/latest"
// // // )

// // // type User struct {
// // // 	ID        int    `json:"id"`
// // // 	FullName  string `json:"full_name"`
// // // 	Email     string `json:"email"`
// // // 	FirstName string `json:"first_name"`
// // // 	LastName  string `json:"last_name"`
// // // }

// // // func PrintUser(user User) {
// // // 	fmt.Println("User Details:")
// // // 	fmt.Printf("ID: %d\n", user.ID)
// // // 	fmt.Printf("FullName: %s\n", user.FullName)
// // // 	fmt.Printf("Email: %s\n", user.Email)
// // // 	fmt.Printf("First Name: %s\n", user.FirstName)
// // // 	fmt.Printf("Last Name: %s\n", user.LastName)
// // // }

// // // // FindUserByEmail ищет пользователя по email в массиве
// // // func FindUserByEmail(users []User, email string) (*User, error) {
// // // 	for _, user := range users {
// // // 		if user.Email == email {
// // // 			return &user, nil
// // // 		}
// // // 	}
// // // 	return nil, fmt.Errorf("user with email '%s' not found", email)
// // // }

// // // // TaskType описывает тип задачи в Kaiten
// // // type TaskType struct {
// // // 	ID   int    `json:"id"`   // ID типа задачи
// // // 	Name string `json:"name"` // Название типа задачи
// // // 	// Добавьте другие поля, если они есть в API
// // // }

// // // // TaskType описывает тип задачи в Kaiten
// // // type TaskIDType int

// // // // Константы для типов задач
// // // const (
// // // 	BugTaskType                 TaskIDType = 7
// // // 	UserStoryTaskType           TaskIDType = 5
// // // 	TaskDiscoveryTaskType       TaskIDType = 11
// // // 	CardTaskType                TaskIDType = 1
// // // 	FeatureTaskType             TaskIDType = 4
// // // 	ACCTaskType                 TaskIDType = 18
// // // 	ImprovementTaskType         TaskIDType = 15
// // // 	VulnerabilityTaskType       TaskIDType = 13
// // // 	IntegrationTasksTaskType    TaskIDType = 19
// // // 	TechDebtTaskType            TaskIDType = 9
// // // 	AdministrativeTasksTaskType TaskIDType = 12
// // // 	TaskDeliveryTaskType        TaskIDType = 6
// // // 	RequestTaskType             TaskIDType = 17
// // // 	EnablerTaskType             TaskIDType = 8
// // // )

// // // // String возвращает строковое представление типа задачи
// // // func (tt TaskIDType) String() string {
// // // 	switch tt {
// // // 	case BugTaskType:
// // // 		return "Bug"
// // // 	case UserStoryTaskType:
// // // 		return "User Story"
// // // 	case TaskDiscoveryTaskType:
// // // 		return "Task Discovery"
// // // 	case CardTaskType:
// // // 		return "Card"
// // // 	case FeatureTaskType:
// // // 		return "Feature"
// // // 	case ACCTaskType:
// // // 		return "ACC"
// // // 	case ImprovementTaskType:
// // // 		return "Улучшение"
// // // 	case VulnerabilityTaskType:
// // // 		return "Уязвимость"
// // // 	case IntegrationTasksTaskType:
// // // 		return "Интеграционные задачи"
// // // 	case TechDebtTaskType:
// // // 		return "Техдолг"
// // // 	case AdministrativeTasksTaskType:
// // // 		return "Административные задачи"
// // // 	case TaskDeliveryTaskType:
// // // 		return "Task Delivery"
// // // 	case RequestTaskType:
// // // 		return "Обращение"
// // // 	case EnablerTaskType:
// // // 		return "Enabler"
// // // 	default:
// // // 		return "Unknown"
// // // 	}
// // // }

// // // // taskTypeMap сопоставляет имя типа задачи с его константой
// // // var taskTypeMap = map[string]TaskIDType{
// // // 	"Bug":            BugTaskType,
// // // 	"User Story":     UserStoryTaskType,
// // // 	"Task Discovery": TaskDiscoveryTaskType,
// // // 	"discovery":      TaskDiscoveryTaskType,
// // // 	"Card":           CardTaskType,
// // // 	"Feature":        FeatureTaskType,
// // // 	"ACC":            ACCTaskType,
// // // 	"Улучшение":      ImprovementTaskType,
// // // 	"Уязвимость":     VulnerabilityTaskType,
// // // 	"Интеграционные задачи": IntegrationTasksTaskType,
// // // 	"Техдолг": TechDebtTaskType,
// // // 	"Административные задачи": AdministrativeTasksTaskType,
// // // 	"Task Delivery": TaskDeliveryTaskType,
// // // 	"delivery":      TaskDeliveryTaskType,
// // // 	"Обращение":     RequestTaskType,
// // // 	"Enabler":       EnablerTaskType,
// // // }

// // // // GetTaskTypeByName возвращает константу типа задачи по его имени
// // // func GetTaskTypeByName(name string) (TaskIDType, error) {
// // // 	if taskType, exists := taskTypeMap[name]; exists {
// // // 		return taskType, nil
// // // 	}
// // // 	return 0, fmt.Errorf("task type '%s' not found", name)
// // // }

// // type Card struct {
// // 	ID            int                    `json:"id"`
// // 	Title         string                 `json:"title"`
// // 	Description   string                 `json:"description"`
// // 	ColumnID      int                    `json:"column_id"`
// // 	BoardID       int                    `json:"board_id"`
// // 	LaneID        int                    `json:"lane_id"`              // Lane ID
// // 	MemberIDs     []int                  `json:"member_ids"`           // Поле для идентификаторов участников карточки
// // 	ParentID      int                    `json:"parent_id,omitempty"`  // ID родительской карточки (если есть)
// // 	TypeID        int                    `json:"type_id"`              // ID типа задачи
// // 	SizeText      string                 `json:"size_text"`            // Размер задачи
// // 	ResponsibleID int                    `json:"responsible_id"`       // Responsible ID
// // 	Properties    map[string]interface{} `json:"properties,omitempty"` // Пользовательские свойства
// // }

// // func PrintCard(card Card) {
// // 	fmt.Println("Card Details:")
// // 	fmt.Printf("ID: %d\n", card.ID)
// // 	fmt.Printf("Title: %s\n", card.Title)
// // 	fmt.Printf("Description: %s\n", card.Description)
// // 	fmt.Printf("ColumnID: %d\n", card.ColumnID)
// // 	fmt.Printf("BoardID: %d\n", card.BoardID)
// // 	fmt.Println("MemberIDs: ", card.MemberIDs)
// // 	fmt.Printf("ParentID: %d\n", card.ParentID)
// // 	fmt.Printf("TypeID: %d\n", card.TypeID)
// // 	fmt.Printf("SizeText: %s\n", card.SizeText)
// // 	fmt.Printf("ResponsibleID: %d\n", card.ResponsibleID)
// // }

// // type KaitenClient struct {
// // 	client  *http.Client
// // 	baseURL string
// // 	token   string
// // }

// // func createHTTPClient() *http.Client {
// // 	return &http.Client{
// // 		Transport: &http.Transport{
// // 			TLSClientConfig: &tls.Config{
// // 				InsecureSkipVerify: true, // Отключаем проверку сертификата
// // 			},
// // 		},
// // 	}
// // }

// // func createTLSConfig(caCertPath string) (*tls.Config, error) {
// // 	caCert, err := os.ReadFile(caCertPath)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
// // 	}

// // 	caCertPool := x509.NewCertPool()
// // 	caCertPool.AppendCertsFromPEM(caCert)

// // 	return &tls.Config{
// // 		RootCAs: caCertPool,
// // 	}, nil
// // }

// // func createHTTPClientWithCertificate(crtFileName string) *http.Client {
// // 	tlsConfig, err := createTLSConfig(crtFileName)
// // 	if err != nil {
// // 		log.Fatalf("Failed to create TLS config: %v", err)
// // 	}

// // 	client := &http.Client{
// // 		Transport: &http.Transport{
// // 			TLSClientConfig: tlsConfig,
// // 		},
// // 	}
// // 	return client
// // }

// // func CreateKaitenClient(token string) *KaitenClient {
// // 	client := createHTTPClient()
// // 	return &KaitenClient{
// // 		client:  client, //&http.Client{},
// // 		baseURL: baseURL,
// // 		token:   token,
// // 	}
// // }

// // func (kc *KaitenClient) doRequest(method, path string, body []byte) (*http.Response, error) {
// // 	req, err := http.NewRequest(method, kc.baseURL+path, bytes.NewBuffer(body))
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	req.Header.Set("Accept", "application/json")
// // 	req.Header.Set("Authorization", "Bearer "+kc.token)
// // 	req.Header.Set("Content-Type", "application/json")

// // 	return kc.client.Do(req)
// // }

// // func (kc *KaitenClient) doRequestWithBody(method, path string, body interface{}) (*http.Response, error) {
// // 	var buf bytes.Buffer
// // 	if body != nil {
// // 		if err := json.NewEncoder(&buf).Encode(body); err != nil {
// // 			return nil, fmt.Errorf("failed to encode request body: %w", err)
// // 		}
// // 	}

// // 	req, err := http.NewRequest(method, kc.baseURL+path, &buf)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to create request: %w", err)
// // 	}

// // 	req.Header.Set("Accept", "application/json")
// // 	req.Header.Set("Authorization", "Bearer "+kc.token)
// // 	req.Header.Set("Content-Type", "application/json")

// // 	return kc.client.Do(req)
// // }

// // func (kc *KaitenClient) decodeResponse(resp *http.Response, v interface{}) error {
// // 	defer resp.Body.Close()
// // 	if resp.StatusCode != http.StatusOK {
// // 		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// // 	}
// // 	return json.NewDecoder(resp.Body).Decode(v)
// // }

// // func (kc *KaitenClient) GetUsers() ([]User, error) {
// // 	resp, err := kc.doRequestWithBody("GET", "/users", nil)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	var users []User
// // 	if err := kc.decodeResponse(resp, &users); err != nil {
// // 		return nil, err
// // 	}

// // 	return users, nil
// // }

// // func (kc *KaitenClient) GetCurrentUser() (*User, error) {
// // 	resp, err := kc.doRequest("GET", "/users/current", nil)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer resp.Body.Close()

// // 	if resp.StatusCode != http.StatusOK {
// // 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// // 	}

// // 	var user User
// // 	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
// // 		return nil, err
// // 	}

// // 	return &user, nil
// // }

// // // GetUser возвращает информацию о конкретном пользователе по ID
// // func (kc *KaitenClient) GetUser(userID int) (*User, error) {
// // 	resp, err := kc.doRequest("GET", fmt.Sprintf("/users/%d", userID), nil)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer resp.Body.Close()

// // 	if resp.StatusCode != http.StatusOK {
// // 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// // 	}

// // 	var user User
// // 	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
// // 		return nil, err
// // 	}

// // 	return &user, nil
// // }

// // // GetUserCards возвращает список карточек для пользователя с указанным идентификатором
// // func (kc *KaitenClient) GetUserCards(userID int) ([]Card, error) {
// // 	resp, err := kc.doRequest("GET", fmt.Sprintf("/users/%d/cards", userID), nil)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer resp.Body.Close()

// // 	if resp.StatusCode != http.StatusOK {
// // 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// // 	}

// // 	var cards []Card
// // 	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
// // 		return nil, err
// // 	}

// // 	return cards, nil
// // }

// // // GetUserByEmail возвращает пользователя по email
// // func (kc *KaitenClient) getUserByEmail(email string) (*User, error) {
// // 	resp, err := kc.doRequestWithBody("GET", "/users?email="+url.QueryEscape(email), nil)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	var users []User
// // 	if err := kc.decodeResponse(resp, &users); err != nil {
// // 		return nil, err
// // 	}

// // 	if len(users) == 0 {
// // 		return nil, fmt.Errorf("user with email %s not found", email)
// // 	}

// // 	return FindUserByEmail(users, email)
// // }
// // func (kc *KaitenClient) GetUserIDByEmail(responsibleEmail string) (int, error) {
// // 	user, err := kc.getUserByEmail(responsibleEmail)
// // 	if err != nil {
// // 		return 0, err
// // 	}

// // 	return user.ID, nil
// // }

// // // GetUserCards возвращает список карточек, в которых участвует пользователь с указанным идентификатором
// // // limit - количество карточек на странице
// // // offset - смещение относительно начала списка
// // func (kc *KaitenClient) getUserCards(userID int, limit int, offset int) ([]Card, error) {
// // 	params := url.Values{}
// // 	params.Add("member_ids", fmt.Sprintf("%d", userID))
// // 	params.Add("limit", fmt.Sprintf("%d", limit))
// // 	params.Add("offset", fmt.Sprintf("%d", offset))
// // 	params.Add("condition", "1")

// // 	resp, err := kc.doRequestWithBody("GET", "/cards?"+params.Encode(), nil)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	var cards []Card
// // 	if err := kc.decodeResponse(resp, &cards); err != nil {
// // 		return nil, err
// // 	}

// // 	return cards, nil
// // }

// // func (kc *KaitenClient) GetUserCardsByMemberIDs(userID int, limit int, offset int) ([]Card, error) {
// // 	return kc.getUserCards(userID, limit, offset)
// // }

// // func (kc *KaitenClient) GetAllUserCards(userID int) ([]Card, error) {
// // 	var allCards []Card
// // 	limit := 100
// // 	offset := 0

// // 	for {
// // 		cards, err := kc.getUserCards(userID, limit, offset)
// // 		if err != nil {
// // 			return nil, err
// // 		}

// // 		if len(cards) == 0 {
// // 			break
// // 		}

// // 		allCards = append(allCards, cards...)
// // 		offset += limit
// // 	}

// // 	return allCards, nil
// // }

// // func (kc *KaitenClient) createUser(user *User) (*User, error) {
// // 	resp, err := kc.doRequestWithBody("POST", "/users", user)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	var newUser User
// // 	if err := kc.decodeResponse(resp, &newUser); err != nil {
// // 		return nil, err
// // 	}

// // 	return &newUser, nil
// // }

// // func (kc *KaitenClient) updateUser(userID int, user *User) (*User, error) {
// // 	resp, err := kc.doRequestWithBody("PUT", fmt.Sprintf("/users/%d", userID), user)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	var updatedUser User
// // 	if err := kc.decodeResponse(resp, &updatedUser); err != nil {
// // 		return nil, err
// // 	}

// // 	return &updatedUser, nil
// // }

// // // CreateUser создает нового пользователя
// // func (kc *KaitenClient) CreateUser(user *User) (*User, error) {
// // 	body, err := json.Marshal(user)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	resp, err := kc.doRequest("POST", "/users", body)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer resp.Body.Close()

// // 	if resp.StatusCode != http.StatusCreated {
// // 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// // 	}

// // 	var newUser User
// // 	if err := json.NewDecoder(resp.Body).Decode(&newUser); err != nil {
// // 		return nil, err
// // 	}

// // 	return &newUser, nil
// // }

// // // UpdateUser обновляет информацию о пользователе
// // func (kc *KaitenClient) UpdateUser(userID int, user *User) (*User, error) {
// // 	body, err := json.Marshal(user)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	resp, err := kc.doRequest("PUT", fmt.Sprintf("/users/%d", userID), body)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer resp.Body.Close()

// // 	if resp.StatusCode != http.StatusOK {
// // 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// // 	}

// // 	var updatedUser User
// // 	if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
// // 		return nil, err
// // 	}

// // 	return &updatedUser, nil
// // }

// // // DeleteUser удаляет пользователя по ID
// // func (kc *KaitenClient) DeleteUser(userID int) error {
// // 	resp, err := kc.doRequest("DELETE", fmt.Sprintf("/users/%d", userID), nil)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer resp.Body.Close()

// // 	if resp.StatusCode != http.StatusNoContent {
// // 		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// // 	}

// // 	return nil
// // }

// func (kc *KaitenClient) getTaskTypes() ([]TaskType, error) {
// 	resp, err := kc.doRequestWithBody("GET", "/card-types", nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch task types: %w", err)
// 	}

// 	var taskTypes []TaskType
// 	if err := kc.decodeResponse(resp, &taskTypes); err != nil {
// 		return nil, fmt.Errorf("failed to decode task types: %w", err)
// 	}

// 	return taskTypes, nil
// }

// func (kc *KaitenClient) GetTaskTypes() ([]TaskType, error) {
// 	return kc.getTaskTypes()
// }

// // GetCards возвращает список всех карт
// func (kc *KaitenClient) GetCards() ([]Card, error) {
// 	resp, err := kc.doRequest("GET", "/cards", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	var cards []Card
// 	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
// 		return nil, err
// 	}

// 	return cards, nil
// }

// func (kc *KaitenClient) getCard(cardID int) (*Card, error) {
// 	resp, err := kc.doRequestWithBody("GET", fmt.Sprintf("/cards/%d", cardID), nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var card Card
// 	if err := kc.decodeResponse(resp, &card); err != nil {
// 		return nil, err
// 	}

// 	return &card, nil
// }
// func (kc *KaitenClient) GetCard(cardID int) (*Card, error) {
// 	return kc.getCard(cardID)
// }

// // CardUpdate описывает данные для обновления карточки
// type CardUpdate struct {
// 	Title         *string                `json:"title,omitempty"`
// 	Description   *string                `json:"description,omitempty"`
// 	ColumnID      *int                   `json:"column_id,omitempty"`
// 	BoardID       *int                   `json:"board_id,omitempty"`
// 	LaneID        *int                   `json:"lane_id,omitempty"`
// 	MemberIDs     *[]int                 `json:"member_ids,omitempty"`
// 	ParentID      *int                   `json:"parent_id,omitempty"`
// 	TypeID        *int                   `json:"type_id,omitempty"`
// 	SizeText      *string                `json:"size_text,omitempty"`
// 	ResponsibleID *int                   `json:"responsible_id,omitempty"`
// 	OwnerID       *int                   `json:"owner_id,omitempty"`
// 	OwnerEmailID  *string                `json:"owner_email,omitempty"`
// 	Properties    map[string]interface{} `json:"properties,omitempty"`
// }

// func (kc *KaitenClient) createCard(card *Card) (*Card, error) {
// 	resp, err := kc.doRequestWithBody("POST", "/cards", card)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var createdCard Card
// 	if err := kc.decodeResponse(resp, &createdCard); err != nil {
// 		return nil, err
// 	}

// 	return &createdCard, nil
// }

// func (kc *KaitenClient) updateCard(cardID int, update CardUpdate) error {
// 	resp, err := kc.doRequestWithBody("PATCH", fmt.Sprintf("/cards/%d", cardID), update)
// 	if err != nil {
// 		return err
// 	}

// 	return kc.decodeResponse(resp, nil)
// }

// func (kc *KaitenClient) CreateCard(card *Card) (*Card, error) {
// 	return kc.createCard(card)
// }

// func (kc *KaitenClient) UpdateCard(cardID int, update CardUpdate) error {
// 	return kc.updateCard(cardID, update)
// }

// func (kc *KaitenClient) updateCardProperties(cardID int, properties map[string]interface{}) error {
// 	updateData := CardUpdate{
// 		Properties: properties,
// 	}

// 	return kc.updateCard(cardID, updateData)
// }

// func (kc *KaitenClient) UpdateCardProperties(cardID int, properties map[string]interface{}) error {
// 	return kc.updateCardProperties(cardID, properties)
// }

// // DeleteCard удаляет карту по ID
// func (kc *KaitenClient) DeleteCard(cardID int) error {
// 	resp, err := kc.doRequest("DELETE", fmt.Sprintf("/cards/%d", cardID), nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusNoContent {
// 		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	return nil
// }

// // TagRequest описывает данные для добавления тега
// type TagRequest struct {
// 	Name string `json:"name"` // Имя тега
// }

// func (kc *KaitenClient) addTagToCard(cardID int, tagName string) error {
// 	tagRequest := TagRequest{Name: tagName}

// 	resp, err := kc.doRequestWithBody("POST", fmt.Sprintf("/cards/%d/tags", cardID), tagRequest)
// 	if err != nil {
// 		return err
// 	}

// 	return kc.decodeResponse(resp, nil)
// }

// func (kc *KaitenClient) AddTagToCard(cardID int, tagName string) error {
// 	return kc.addTagToCard(cardID, tagName)
// }

// // Add children Request
// type AddChildrenRequest struct {
// 	CardID int `json:"card_id"`
// }

// func (kc *KaitenClient) addChildrenToCard(cardID int, childrenCardID int) error {
// 	requestData := AddChildrenRequest{CardID: childrenCardID}

// 	resp, err := kc.doRequestWithBody("POST", fmt.Sprintf("/cards/%d/children", cardID), requestData)
// 	if err != nil {
// 		return err
// 	}

// 	return kc.decodeResponse(resp, nil)
// }

// func (kc *KaitenClient) AddChindrenToCard(cardID int, childrenCardID int) error {
// 	return kc.addChildrenToCard(cardID, childrenCardID)
// }

// // Task описывает задачу для создания карточки
// type Task struct {
// 	Type  string `json:"type"`  // Тип задачи (например, "delivery", "discovery")
// 	Size  int    `json:"size"`  // Размер задачи (например, 8, 16)
// 	Title string `json:"title"` // Название задачи
// }

// // Schedule описывает расписание задач
// type Schedule struct {
// 	Parent      string `json:"parent"`      // ID родительской карточки
// 	Responsible string `json:"responsible"` // Email ответственного
// 	Tasks       []Task `json:"tasks"`       // Список задач
// }

// // ScheduleFile описывает структуру JSON-файла
// type ScheduleFile struct {
// 	Schedule Schedule `json:"schedule"`
// }

// func loadFromJSON(filePath string, v interface{}) error {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to open file: %w", err)
// 	}
// 	defer file.Close()

// 	decoder := json.NewDecoder(file)
// 	if err := decoder.Decode(v); err != nil {
// 		return fmt.Errorf("failed to decode JSON: %w", err)
// 	}

// 	return nil
// }

// func LoadTasksFromJSON(filePath string) (*Schedule, error) {
// 	var scheduleFile ScheduleFile
// 	if err := loadFromJSON(filePath, &scheduleFile); err != nil {
// 		return nil, err
// 	}

// 	return &scheduleFile.Schedule, nil
// }

// // ExtractWorkCode извлекает код работы из названия родительской карточки
// func ExtractWorkCode(parentTitle string) (string, error) {

// 	// Регулярное выражение для поиска кода работы
// 	// Шаблон: номер фичи (цифры или символы), точка, номер пользовательской истории (цифры или символы)
// 	re := regexp.MustCompile(`\[CAD\]:[A-Za-z]+\.([^.\s]+\.[^.\s]+)`)
// 	match := re.FindStringSubmatch(parentTitle)

// 	if len(match) < 2 {
// 		return "", fmt.Errorf("work code not found in title: %s", parentTitle)
// 	}

// 	return match[1], nil
// }

type CRUD interface {
	Create(interface{}) (interface{}, error)
	Update(int, interface{}) error
	Delete(int) error
}

type CardService struct {
	client *api.KaitenClient
}

func (cs *CardService) Create(item interface{}) (interface{}, error) {
	card, ok := item.(*models.Card)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}
	return cs.client.CreateCard(card)
}

func (cs *CardService) Update(id int, item interface{}) error {
	update, ok := item.(models.CardUpdate)
	if !ok {
		return fmt.Errorf("invalid type")
	}
	return cs.client.UpdateCard(id, update)
}

func (cs *CardService) Delete(id int) error {
	return cs.client.DeleteCard(id)
}

func parseFlags() (string, string) {
	tasksFile := flag.String("tasks", "", "Path to the tasks JSON file (required)")
	configFile := flag.String("config", "config.json", "Path to the configuration file (optional, default: config.json)")

	flag.Parse()

	if *tasksFile == "" {
		fmt.Println("Error: The 'tasks' flag is required.")
		flag.Usage()
		*tasksFile = "tasks.json"
	}

	return *tasksFile, *configFile
}

func main() {

	tasksFile, configFile := parseFlags()

	// Загрузка конфигурации (если указан файл конфигурации)
	if configFile != "" {
		fmt.Printf("Loading configuration from: %s\n", configFile)
	} else {
		fmt.Println("Using default configuration.")
	}

	client := api.CreateKaitenClient("9ecb4b54-508a-4d1e-ad99-c1c4a04847bb")

	// Получение данных текущего пользователя
	currentUser, err := client.GetCurrentUser()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	// Загрузка задач из JSON-файла
	schedule, err := utils.LoadTasksFromJSON(tasksFile)
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

	parentCard, err := client.GetCard(parentID)
	if err != nil {
		fmt.Println("Error parsing parent card by ID:", err)
		return
	}

	parentCardWorkCode := "XXX.XX"
	if len(parentCard.Title) > 0 {
		workCode, err := utils.ExtractWorkCode(parentCard.Title)
		if err != nil {
			fmt.Println("Error:", err)
			// return
		} else {
			parentCardWorkCode = workCode
		}
		//fmt.Printf("Work code: %s\n", workCode)
	} else {
		fmt.Println("Error parent card title is empty")
	}

	// Создаем карточки для каждой задачи
	for _, task := range schedule.Tasks {
		// Определяем тип задачи
		taskTypeID, err := models.GetTaskTypeByName(task.Type)
		if err != nil {
			fmt.Printf("Error getting task type ID for '%s': %v\n", task.Type, err)
			continue
		}

		// Создаем карточку
		card := &models.Card{
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

		//PrintCard(*card)

		// Создаем карточку в Kaiten
		createdCard, err := client.CreateCard(card)
		if err != nil {
			fmt.Printf("Error creating card '%s': %v\n", task.Title, err)
			continue
		} else {
			fmt.Printf("Created card: %s (ID: %d)\n", createdCard.Title, createdCard.ID)
		}

		// cardService := &CardService{client: client}
		// createdCard, err := cardService.Create(card)
		// if err != nil {
		// 	fmt.Println("Error creating card:", err)
		// }

		if createdCard.TypeID == int(models.TaskDeliveryTaskType) || createdCard.TypeID == int(models.TaskDiscoveryTaskType) {

			titleUpdate := fmt.Sprintf("[CAD]:TS.%s.%d. %s", parentCardWorkCode, createdCard.ID, createdCard.Title)
			updateData := &models.CardUpdate{
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

}

// func PrintCardsList(cards []Card, userID int) {
// 	fmt.Printf("Cards count=%d :\n", len(cards))
// 	fmt.Printf("Cards for user %d:\n", userID)
// 	for _, card := range cards {
// 		fmt.Printf("Card ID: %d, Title: %s, Description: %s, ColumnID: %d, BoardID: %d\n",
// 			card.ID, card.Title, card.Description, card.ColumnID, card.BoardID)
// 	}
// }

// Вспомогательные функции для создания указателей
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
