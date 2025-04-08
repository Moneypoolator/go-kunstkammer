package main

import (
	"flag"
	"kunstkammer/internal/utils"
	"kunstkammer/pkg/config"
	"log/slog"
	"os"
)

func parseFlags() (string, string, string) {
	tasksFile := flag.String("tasks", "", "Path to the tasks JSON file for task creation mode")
	reportFile := flag.String("report", "", "Path to the report JSON file for report import mode")
	configFile := flag.String("config", "config.json", "Path to the configuration file (optional, default: config.json)")

	flag.Parse()

	if *tasksFile == "" && *reportFile == "" {
		slog.Error("Wrong command line arguments")
		flag.Usage()
	}

	return *tasksFile, *configFile, *reportFile
}

func main() {

	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	tasksFile, configFile, reportFile := parseFlags()

	var env config.Config

	// Загрузка конфигурации (если указан файл конфигурации)
	if configFile != "" {
		slog.Info("Loading configuration from", "file", configFile)

		// Загружаем конфигурацию
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			slog.Error("Failed to load config", "error", err)
			return
		}

		// Используем конфигурацию
		slog.Debug("Config loaded", "Token", cfg.Token)
		slog.Debug("Config loaded", "BaseURL", cfg.BaseURL)
		slog.Debug("Config loaded", "LogLevel", cfg.LogLevel)
		slog.Debug("Config loaded", "BoardID", cfg.BoardID)
		slog.Debug("Config loaded", "ColumnID", cfg.ColumnID)
		slog.Debug("Config loaded", "LaneID", cfg.LaneID)

		for _, tag := range cfg.Tags {
			slog.Debug("Config loaded", "Tag", tag)
		}

		env = *cfg
	} else {
		slog.Error("Empty config file name")
		return
	}

	if tasksFile != "" {

		// Загрузка задач из JSON-файла
		schedule, err := utils.LoadTasksFromJSON(tasksFile)
		if err != nil {
			slog.Error("Loading tasks from JSON file", "error", err)
			return
		}

		// Вывод задач
		slog.Debug("Tasks loaded", "Parent", schedule.Parent)
		slog.Debug("Tasks loaded", "Responsible", schedule.Responsible)

		err = AsyncProcessTasks(env, env.Token, env.BaseURL, schedule)
		if err != nil {
			slog.Error("Process Tasks", "error", err)
		}
	}

	if reportFile != "" {
		slog.Info("Loading configuration from", "file", configFile)

		report, err := utils.LoadReportFromJSON(reportFile)
		if err != nil {
			slog.Error("Loading report from JSON file", "error", err)
			return
		}

		if report != nil {
			slog.Info("Start report process")

			slog.Debug("Report loaded", "SprintID", report.SprintID)
			slog.Debug("Report loaded", "Responsible", report.Responsible)

		}
	}
}
