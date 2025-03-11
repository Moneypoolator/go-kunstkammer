package main

import (
	"fmt"
	"kunstkammer/internal/api"
	"kunstkammer/internal/models"
	"kunstkammer/internal/utils"
	"log/slog"
	"strconv"
)

func ProcessTasks(token string, kaitenURL string, schedule *models.Schedule) error {

	client := api.CreateKaitenClient(token, kaitenURL)

	// Получение данных текущего пользователя
	currentUser, err := client.GetCurrentUser()
	if err != nil {
		slog.Error("Getting current user", "error", err)
		return err
	}

	// Преобразуем email ответственного в ID пользователя
	responsibleID, err := client.GetUserIDByEmail(schedule.Responsible)
	if responsibleID == 0 || err != nil {
		slog.Warn("Getting responsible user ID:", "error", err)
		responsibleID = currentUser.ID
	}

	// Преобразуем ID родительской карточки из строки в число
	parentID, err := strconv.Atoi(schedule.Parent)
	if err != nil {
		slog.Error("Parsing parent ID", "error", err)
		return err
	}

	parentCard, err := client.GetCard(parentID)
	if err != nil {
		slog.Error("Parsing parent card by ID", "error", err)
		return err
	}

	parentCardWorkCode := "XXX.XX"
	if len(parentCard.Title) > 0 {
		workCode, err := utils.ExtractWorkCode(parentCard.Title)
		if err != nil {
			slog.Warn("Extract Work Code", "error", err)
		} else {
			parentCardWorkCode = workCode
			slog.Debug("Work code", "code", workCode)
		}
	} else {
		slog.Warn("Parent card title is empty")
	}

	// Создаем карточки для каждой задачи
	for _, task := range schedule.Tasks {
		// Определяем тип задачи
		taskTypeID, err := models.GetTaskTypeByName(task.Type)
		if err != nil {
			slog.Error("Getting task type ID", "task_type", task.Type, "error", err)
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
			slog.Error("Creating card", "Title", task.Title, "error", err)
			continue
		} else {
			slog.Info("Created card", "Title", createdCard.Title, "ID", createdCard.ID)
		}

		// cardService := &CardService{client: client}
		// createdCard, err := cardService.Create(card)
		// if err != nil {
		// 	fmt.Println("Error creating card:", err)
		// }

		if createdCard.TypeID == int(models.TaskDeliveryTaskType) || createdCard.TypeID == int(models.TaskDiscoveryTaskType) {

			titleUpdate := fmt.Sprintf("[CAD]:TS.%s.%d. %s", parentCardWorkCode, createdCard.ID, createdCard.Title)
			updateData := &models.CardUpdate{
				Title: utils.StringPtr(titleUpdate),
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
				slog.Error("Updating card", "Title", titleUpdate, "error", err)
			} else {
				slog.Info("Updated card", "Title", titleUpdate, "ID", createdCard.ID)
			}
		}

		err = client.AddChindrenToCard(parentID, createdCard.ID)
		if err != nil {
			slog.Error("Adding children to card", "error", err)
			//return
		}

		// Добавляем тег к карточке
		err = client.AddTagToCard(createdCard.ID, "ГГИС")
		if err != nil {
			slog.Error("Adding tag to card", "error", err)
			//return
		}

		// Добавляем тег к карточке
		err = client.AddTagToCard(createdCard.ID, "C++")
		if err != nil {
			slog.Error("Adding tag to card:", "error", err)
			//return
		}

	}

	return nil
}
