package api

import (
	"fmt"
)

// Column represents a Kaiten column
type Column struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetColumn returns information about a specific column
func (kc *KaitenClient) GetColumn(columnID int) (*Column, error) {
	resp, err := kc.doRequestWithBody("GET", fmt.Sprintf("/columns/%d", columnID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch column: %w", err)
	}

	var column Column
	if err := kc.decodeResponse(resp, &column); err != nil {
		return nil, fmt.Errorf("failed to decode column: %w", err)
	}

	return &column, nil
}

// GetColumns returns all columns for a specific board
func (kc *KaitenClient) GetColumns(boardID int) ([]Column, error) {
	resp, err := kc.doRequestWithBody("GET", fmt.Sprintf("/boards/%d/columns", boardID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch columns: %w", err)
	}

	var columns []Column
	if err := kc.decodeResponse(resp, &columns); err != nil {
		return nil, fmt.Errorf("failed to decode columns: %w", err)
	}

	return columns, nil
}
