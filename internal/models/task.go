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
// Если name пустая строка или nil, возвращает CardTaskType как тип по умолчанию
func GetTaskTypeByName(name string) (TaskIDType, error) {
	if name == "" {
		return TaskDeliveryTaskType, nil
	}
	if taskType, exists := taskTypeMap[name]; exists {
		return taskType, nil
	}
	return 0, fmt.Errorf("task type '%s' not found", name)
}

// GetTaskType безопасно извлекает тип задачи из структуры Task
func GetTaskType(task *Task) TaskIDType {
	if task == nil || task.Type == nil {
		return TaskDeliveryTaskType // используем Card как тип по умолчанию
	}
	taskType, err := GetTaskTypeByName(*task.Type)
	if err != nil {
		return TaskDeliveryTaskType // в случае ошибки также используем тип по умолчанию
	}
	return taskType
}

// Task описывает задачу для создания карточки
type Task struct {
	Type  *string `json:"type,omitempty"` // Тип задачи (например, "delivery", "discovery" -- задачи данного типа не доступны)
	Size  int     `json:"size"`           // Размер задачи (например, 8, 16)
	Title string  `json:"title"`          // Название задачи
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

// NewTask создает новую задачу с опциональным типом
func NewTask(title string, size int, taskType *string) Task {
	return Task{
		Title: title,
		Size:  size,
		Type:  taskType,
	}
}

// SetType устанавливает тип задачи
func (t *Task) SetType(taskType string) {
	t.Type = &taskType
}

// GetTypeString возвращает строковое представление типа задачи
func (t *Task) GetTypeString() string {
	if t.Type == nil {
		return TaskDeliveryTaskType.String()
	}
	return *t.Type
}

// Validate проверяет корректность задачи
func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("empty task title")
	}
	if t.Size <= 0 {
		return fmt.Errorf("invalid task size: %d", t.Size)
	}
	if t.Type != nil {
		if _, err := GetTaskTypeByName(*t.Type); err != nil {
			return fmt.Errorf("invalid task type: %s", *t.Type)
		}
	}
	return nil
}
