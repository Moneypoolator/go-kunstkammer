package api

import (
	"fmt"
	"kunstkammer/internal/models"
)

func (kc *KaitenClient) getTaskTypes() ([]models.TaskType, error) {
	resp, err := kc.doRequestWithBody("GET", "/card-types", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch task types: %w", err)
	}

	var taskTypes []models.TaskType
	if err := kc.decodeResponse(resp, &taskTypes); err != nil {
		return nil, fmt.Errorf("failed to decode task types: %w", err)
	}

	return taskTypes, nil
}

func (kc *KaitenClient) GetTaskTypes() ([]models.TaskType, error) {
	return kc.getTaskTypes()
}
