package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config представляет структуру конфигурации приложения.
type Config struct {
	Token    string   `json:"token"`     // Токен для доступа к API Kaiten
	BaseURL  string   `json:"base_url"`  // Базовый URL API Kaiten
	LogLevel string   `json:"log_level"` // Уровень логирования (debug, info, error)
	BoardID  int      `json:"board_id"`  // по-умолчанию  192,
	ColumnID int      `json:"column_id"` // по-умолчанию  776,
	LaneID   int      `json:"lane_id"`   // по-умолчанию  1275
	Tags     []string `json:"tags"`      // Список задач
}

// LoadConfig загружает конфигурацию из JSON-файла или переменных окружения.
// Можно создать переменные среды окружения:
// export KAITEN_TOKEN="your-kaiten-api-token"
// export KAITEN_BASE_URL="https://kaiten.norsoft.ru/api/latest"
// export KAITEN_LOG_LEVEL="debug"
func LoadConfig(filePath string) (*Config, error) {
	var cfg Config

	// Пытаемся загрузить конфигурацию из файла, если он указан
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&cfg); err != nil {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}
	}

	// Переменные окружения имеют приоритет над значениями из файла
	if token := os.Getenv("KAITEN_TOKEN"); token != "" {
		cfg.Token = token
	}
	if baseURL := os.Getenv("KAITEN_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}
	if logLevel := os.Getenv("KAITEN_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	// Проверяем обязательные поля
	if cfg.Token == "" {
		return nil, fmt.Errorf("token is required in config")
	}
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base_url is required in config")
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info" // Устанавливаем значение по умолчанию
	}

	return &cfg, nil
}
