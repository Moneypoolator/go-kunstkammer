package api

import (
	"fmt"
	"kunstkammer/internal/models"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

// columnCache caches column information to avoid repeated API calls
var columnCache = struct {
	sync.RWMutex
	columns map[int]string
}{
	columns: make(map[int]string),
}

// getColumnName returns the name of a column, using cache if available
func (kc *KaitenClient) getColumnName(columnID int) (string, error) {
	// Check cache first
	columnCache.RLock()
	if name, exists := columnCache.columns[columnID]; exists {
		columnCache.RUnlock()
		return name, nil
	}
	columnCache.RUnlock()

	// If not in cache, fetch from API
	column, err := kc.GetColumn(columnID)
	if err != nil {
		return "", fmt.Errorf("failed to get column name: %w", err)
	}

	// Update cache
	columnCache.Lock()
	columnCache.columns[columnID] = column.Name
	columnCache.Unlock()

	return column.Name, nil
}

// GetSprintTasks returns all tasks for a given sprint and user
func (kc *KaitenClient) GetSprintTasks(sprintID int, userEmail string) ([]models.Card, error) {
	// First, get the user ID by email
	userID, err := kc.GetUserIDByEmail(userEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Prepare query parameters
	params := url.Values{}
	params.Add("member_ids", fmt.Sprintf("%d", userID))
	params.Add("limit", "100") // Maximum number of cards per request

	var allCards []models.Card
	offset := 0

	for {
		// Update offset for pagination
		params.Set("offset", fmt.Sprintf("%d", offset))

		// Make the API request
		resp, err := kc.doRequestWithBody("GET", "/cards?"+params.Encode(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch sprint tasks: %w", err)
		}

		var cards []models.Card
		if err := kc.decodeResponse(resp, &cards); err != nil {
			return nil, fmt.Errorf("failed to decode sprint tasks: %w", err)
		}

		// If no more cards, break the loop
		if len(cards) == 0 {
			break
		}

		// Filter cards by sprint ID in Properties
		for _, card := range cards {
			if card.Properties != nil {
				if value, exists := card.Properties["id_12"]; exists {
					// Convert sprint ID to int for comparison
					var cardSprintID int
					switch v := value.(type) {
					case float64:
						cardSprintID = int(v)
					case int:
						cardSprintID = v
					case string:
						if id, err := strconv.Atoi(v); err == nil {
							cardSprintID = id
						}
					}

					if cardSprintID == sprintID {
						allCards = append(allCards, card)
					}
				}
			}
		}

		// Move to the next page
		offset += len(cards)
	}

	return allCards, nil
}

// GetSprintReport generates a report for a given sprint and user
func (kc *KaitenClient) GetSprintReport(sprintID int, userEmail string) (*models.Report, error) {
	cards, err := kc.GetSprintTasks(sprintID, userEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to get sprint tasks: %w", err)
	}

	// Create a report with the sprint ID and user email
	report := &models.Report{
		SprintID:    sprintID,
		Responsible: userEmail,
		TotalTasks:  len(cards),
		Tasks:       make([]models.ReportTask, 0, len(cards)),
	}

	// Process each card and add it to the report
	for _, card := range cards {
		// Extract size from SizeText (format: "X ч")
		size := 0
		if card.SizeText != "" {
			sizeStr := strings.TrimSuffix(card.SizeText, " ч")
			if sizeInt, err := strconv.Atoi(sizeStr); err == nil {
				size = sizeInt
			}
		}

		// Get task type name
		taskType := models.TaskIDType(card.TypeID).String()

		// Get column name as status
		status, err := kc.getColumnName(card.ColumnID)
		if err != nil {
			status = "Unknown" // Use "Unknown" if we can't get the column name
		}

		// Create report task
		reportTask := models.ReportTask{
			ID:          card.ID,
			Title:       card.Title,
			Type:        taskType,
			Size:        size,
			Status:      status,
			Description: card.Description,
		}

		report.Tasks = append(report.Tasks, reportTask)
		report.TotalHours += size
	}

	return report, nil
}
