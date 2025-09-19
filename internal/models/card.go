package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Card struct {
	ID            int                    `json:"id"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	ColumnID      int                    `json:"column_id"`
	BoardID       int                    `json:"board_id"`
	LaneID        int                    `json:"lane_id"`              // Lane ID
	MemberIDs     []int                  `json:"member_ids"`           // Поле для идентификаторов участников карточки
	ParentID      int                    `json:"parent_id,omitempty"`  // ID родительской карточки (если есть)
	TypeID        int                    `json:"type_id"`              // ID типа задачи
	SizeText      string                 `json:"size_text"`            // Размер задачи
	ResponsibleID int                    `json:"responsible_id"`       // Responsible ID
	Properties    map[string]interface{} `json:"properties,omitempty"` // Пользовательские свойства
}

func PrintCard(card Card) {
	fmt.Println("Card Details:")
	fmt.Printf("ID: %d\n", card.ID)
	fmt.Printf("Title: %s\n", card.Title)
	fmt.Printf("Description: %s\n", card.Description)
	fmt.Printf("ColumnID: %d\n", card.ColumnID)
	fmt.Printf("BoardID: %d\n", card.BoardID)
	fmt.Println("MemberIDs: ", card.MemberIDs)
	fmt.Printf("ParentID: %d\n", card.ParentID)
	fmt.Printf("TypeID: %d\n", card.TypeID)
	fmt.Printf("SizeText: %s\n", card.SizeText)
	fmt.Printf("ResponsibleID: %d\n", card.ResponsibleID)

	// Вывод Properties в формате JSON
	fmt.Println("Properties:")
	if len(card.Properties) > 0 {
		propsJSON, err := json.MarshalIndent(card.Properties, "  ", "  ")
		if err != nil {
			fmt.Println("  Error formatting properties:", err)
		} else {
			fmt.Println(string(propsJSON))
		}
	} else {
		fmt.Println("  No properties")
	}
}

// CardUpdate описывает данные для обновления карточки
type CardUpdate struct {
	Title         *string                `json:"title,omitempty"`
	Description   *string                `json:"description,omitempty"`
	ColumnID      *int                   `json:"column_id,omitempty"`
	BoardID       *int                   `json:"board_id,omitempty"`
	LaneID        *int                   `json:"lane_id,omitempty"`
	MemberIDs     *[]int                 `json:"member_ids,omitempty"`
	ParentID      *int                   `json:"parent_id,omitempty"`
	TypeID        *int                   `json:"type_id,omitempty"`
	SizeText      *string                `json:"size_text,omitempty"`
	ResponsibleID *int                   `json:"responsible_id,omitempty"`
	OwnerID       *int                   `json:"owner_id,omitempty"`
	OwnerEmailID  *string                `json:"owner_email,omitempty"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}

func PrintCardsList(cards []Card, userID int) {
	fmt.Printf("Cards count=%d :\n", len(cards))
	fmt.Printf("Cards for user %d:\n", userID)
	for _, card := range cards {
		fmt.Printf("Card ID: %d, Title: %s, Description: %s, ColumnID: %d, BoardID: %d\n",
			card.ID, card.Title, card.Description, card.ColumnID, card.BoardID)
	}
}

// GetPropertyInt extracts an integer property by key
func (c *Card) GetPropertyInt(key string) (int, bool) {
	if c.Properties == nil {
		return 0, false
	}

	value, exists := c.Properties[key]
	if !exists {
		return 0, false
	}

	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case string:
		if intVal, err := strconv.Atoi(v); err == nil {
			return intVal, true
		}
	}
	return 0, false
}

// GetPropertyString extracts a string property by key
func (c *Card) GetPropertyString(key string) (string, bool) {
	if c.Properties == nil {
		return "", false
	}

	value, exists := c.Properties[key]
	if !exists {
		return "", false
	}

	switch v := value.(type) {
	case string:
		return v, true
	case int:
		return strconv.Itoa(v), true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	}
	return "", false
}

// GetSprintNumber extracts sprint number from property id_12
func (c *Card) GetSprintNumber() (int, bool) {
	return c.GetPropertyInt("id_12")
}

// GetRoleID extracts role ID from property id_19
func (c *Card) GetRoleID() (int, bool) {
	return c.GetPropertyInt("id_19")
}

// GetPropertyObject extracts an object/map property by key
func (c *Card) GetPropertyObject(key string) (map[string]interface{}, bool) {
	if c.Properties == nil {
		return nil, false
	}

	value, exists := c.Properties[key]
	if !exists || value == nil {
		return nil, false
	}

	if obj, ok := value.(map[string]interface{}); ok {
		return obj, true
	}
	return nil, false
}

// DefaultTeamPropertyKey is the default property key for team on a card.
// Set this to the actual property id, e.g., "id_21".
var DefaultTeamPropertyKey = "id_143"

// GetTeamIDFrom extracts team ID from a specific property key
func (c *Card) GetTeamIDFrom(propertyKey string) (int, bool) {
	return c.GetPropertyInt(propertyKey)
}

// GetTeamNameFrom extracts team name from a specific property key.
// Supports both object { name: string } and direct string values.
func (c *Card) GetTeamNameFrom(propertyKey string) (string, bool) {
	// Try object form first
	if obj, ok := c.GetPropertyObject(propertyKey); ok {
		if nameVal, exists := obj["value"]; exists {
			// if nameVal, exists := obj["name"]; exists {
			switch v := nameVal.(type) {
			case string:
				return v, true
			case int:
				return strconv.Itoa(v), true
			case float64:
				return strconv.FormatFloat(v, 'f', -1, 64), true
			}
		}
	}
	// Fallback to plain string/int/float
	return c.GetPropertyString(propertyKey)
}

// GetTeamID extracts team ID using DefaultTeamPropertyKey
func (c *Card) GetTeamID() (int, bool) {
	return c.GetTeamIDFrom(DefaultTeamPropertyKey)
}

// GetTeamName extracts team name using DefaultTeamPropertyKey
func (c *Card) GetTeamName() (string, bool) {
	return c.GetTeamNameFrom(DefaultTeamPropertyKey)
}
