package models

// Report описывает отчет по задачам в спринте
type Report struct {
	SprintID    int          `json:"sprint"`      // Номер спринта
	Responsible string       `json:"responsible"` // Email ответственного
	TotalTasks  int          `json:"total_tasks"` // Общее количество задач
	TotalHours  int          `json:"total_hours"` // Общее количество часов
	Tasks       []ReportTask `json:"tasks"`       // Список задач
}

// ReportTask описывает задачу в отчете
type ReportTask struct {
	ID             int    `json:"id"`             // ID задачи
	Title          string `json:"title"`          // Название задачи
	Type           string `json:"type"`           // Тип задачи
	Size           int    `json:"size"`           // Размер в часах
	Status         string `json:"status"`         // Статус задачи / колонка
	Description    string `json:"description"`    // Описание задачи
	SprintNumber   int    `json:"sprintNumber"`   // Номер спринта
	Responsible    string `json:"responsible"`    // Ответственный
	Column         string `json:"column"`         // Колонка
	Team           string `json:"team"`           // Команда
	ParentIDs      []int  `json:"parentIds"`      // Родительские id
	Role           string `json:"role"`           // Роль
	IsBlocked      bool   `json:"isBlocked"`      // Заблокирована
	TotalTime      int    `json:"totalTime"`      // Суммарное время
	SizeUnit       string `json:"sizeUnit"`       // Единица изменения размера
	ChildSizeSum   int    `json:"childSizeSum"`   // Суммарный размер дочерних карточек
	ChildSprintSum int    `json:"childSprintSum"` // Номер спринта (суммарное значение дочерних карточек)
}

// ScheduleFile описывает структуру JSON-файла
type ReportFile struct {
	Report Report `json:"report"`
}
