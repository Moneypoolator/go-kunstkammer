package main

import (
	"fmt"
	"kunstkammer/internal/api"
	"kunstkammer/internal/models"
	"kunstkammer/internal/utils"
	"log/slog"
	"strconv"
	"sync"
)

func AsyncProcessTasks(token string, kaitenURL string, schedule *models.Schedule) error {

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

	// Канал для передачи ошибок
	errorsChannel := make(chan error, len(schedule.Tasks))
	// WaitGroup для ожидания завершения всех горутин
	var tasksCreationWaitGroup sync.WaitGroup

	// Создаем карточки для каждой задачи асинхронно
	for _, task := range schedule.Tasks {
		tasksCreationWaitGroup.Add(1) // Увеличиваем счетчик WaitGroup

		go func(task models.Task) {
			defer tasksCreationWaitGroup.Done() // Уменьшаем счетчик WaitGroup при завершении горутины

			// Определяем тип задачи
			taskTypeID, err := models.GetTaskTypeByName(task.Type)
			if err != nil {
				slog.Error("Getting task type ID", "task_type", task.Type, "error", err)
				errorsChannel <- err
				return
			}

			// Заполняем карточку
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

			// Создаем карточку в Kaiten
			createdCard, err := client.CreateCard(card)
			if err != nil {
				slog.Error("Creating card", "Title", task.Title, "error", err)
				errorsChannel <- err
				return
			} else {
				slog.Info("Created card", "Title", createdCard.Title, "ID", createdCard.ID)
			}

			// Обновляем заголовок карточки, если это необходимо
			if createdCard.TypeID == int(models.TaskDeliveryTaskType) || createdCard.TypeID == int(models.TaskDiscoveryTaskType) {
				titleUpdate := fmt.Sprintf("[CAD]:TS.%s.%d. %s", parentCardWorkCode, createdCard.ID, createdCard.Title)
				updateData := &models.CardUpdate{
					Title: utils.StringPtr(titleUpdate),
				}

				err = client.UpdateCard(createdCard.ID, *updateData)
				if err != nil {
					slog.Error("Updating card", "Title", titleUpdate, "error", err)
					errorsChannel <- err
				} else {
					slog.Info("Updated card", "Title", titleUpdate, "ID", createdCard.ID)
				}
			}

			// Добавляем карточку как дочернюю
			err = client.AddChindrenToCard(parentID, createdCard.ID)
			if err != nil {
				slog.Error("Adding children to card", "error", err)
				errorsChannel <- err
			}

			// Добавляем теги к карточке
			tags := []string{"ГГИС", "C++"}
			for _, tag := range tags {
				err = client.AddTagToCard(createdCard.ID, tag)
				if err != nil {
					slog.Error("Adding tag to card", "tag", tag, "error", err)
					errorsChannel <- err
				}
			}
		}(task)
	}

	// Ожидаем завершения всех горутин
	tasksCreationWaitGroup.Wait()
	close(errorsChannel) // Закрываем канал ошибок

	// Собираем все ошибки
	var errors []error
	for err := range errorsChannel {
		errors = append(errors, err)
	}

	// Если были ошибки, возвращаем их
	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors while processing tasks", len(errors))
	}

	return nil
}
