package main

import (
	"flag"
	"fmt"
	"kunstkammer/internal/api"
	"kunstkammer/internal/models"
	"kunstkammer/internal/utils"
	"kunstkammer/pkg/config"
	"strconv"
)

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

	var env config.Config

	// Загрузка конфигурации (если указан файл конфигурации)
	if configFile != "" {
		fmt.Printf("Loading configuration from: %s\n", configFile)

		// Загружаем конфигурацию
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			return
		}

		// Используем конфигурацию
		fmt.Printf("Token: %s\n", cfg.Token)
		fmt.Printf("BaseURL: %s\n", cfg.BaseURL)
		fmt.Printf("LogLevel: %s\n", cfg.LogLevel)

		env = *cfg
	} else {
		fmt.Println("Empty config file name")
		return
	}

	client := api.CreateKaitenClient(env.Token, env.BaseURL)

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

		// Заполняем карточку
		//TOOD: Вынести магические числа в конфиг файл
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

func Ptr[T any](value T) *T {
	return &value
}
