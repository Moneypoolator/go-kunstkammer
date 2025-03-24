package models

// Schedule описывает расписание задач
type Report struct {
	SprintID    int    `json:"sprint"`      // Номер спринта
	Responsible string `json:"responsible"` // Email ответственного
}

// ScheduleFile описывает структуру JSON-файла
type ReportFile struct {
	Report Report `json:"report"`
}
