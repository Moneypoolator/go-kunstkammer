package utils

import (
	"encoding/json"
	"fmt"
	"kunstkammer/internal/models"
	"os"
)

func loadFromJSON(filePath string, v interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

func LoadTasksFromJSON(filePath string) (*models.Schedule, error) {
	var scheduleFile models.ScheduleFile
	if err := loadFromJSON(filePath, &scheduleFile); err != nil {
		return nil, err
	}

	return &scheduleFile.Schedule, nil
}

func LoadReportFromJSON(filePath string) (*models.Report, error) {
	var reportFile models.ReportFile
	if err := loadFromJSON(filePath, &reportFile); err != nil {
		return nil, err
	}

	return &reportFile.Report, nil
}
