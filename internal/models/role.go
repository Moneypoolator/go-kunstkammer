package models

// RoleType описывает тип роли к которой принадлежит карточка в Kaiten
type RoleType int

// Константы для роли задач
const (
	CppRole         RoleType = 1
	BackendRole     RoleType = 2
	FrontendRole    RoleType = 3
	AnalystRole     RoleType = 8
	UIUXRole        RoleType = 9
	DevOpsRole      RoleType = 11
	WriterRole      RoleType = 12
	ApprobationRole RoleType = 20
)

// String возвращает строковое представление типа задачи
func (rt RoleType) String() string {
	switch rt {
	case CppRole:
		return "C++"
	case BackendRole:
		return "Backend"
	case FrontendRole:
		return "Frontend"
	case AnalystRole:
		return "Analyst"
	case UIUXRole:
		return "UI/UX"
	case DevOpsRole:
		return "DevOps"
	case WriterRole:
		return "Writer"
	case ApprobationRole:
		return "Внедрение"

	default:
		return "Unknown"
	}
}
