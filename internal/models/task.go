package models

import (
	"fmt"
)

// TaskType описывает тип задачи в Kaiten
type TaskType struct {
	ID   int    `json:"id"`   // ID типа задачи
	Name string `json:"name"` // Название типа задачи
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

// Task описывает задачу для создания карточки
type Task struct {
	Type  string `json:"type"`  // Тип задачи (например, "delivery", "discovery" -- задачи данного типа не доступны)
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
