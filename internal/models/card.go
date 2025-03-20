package models

import (
	"encoding/json"
	"fmt"
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
