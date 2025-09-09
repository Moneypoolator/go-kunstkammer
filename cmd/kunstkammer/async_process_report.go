package main

import (
	"fmt"
	"kunstkammer/internal/api"
	"kunstkammer/internal/models"
	"kunstkammer/pkg/config"
	"log/slog"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func saveReportToExcel(report *models.Report) error {
	// Create a new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Error closing Excel file", "error", err)
		}
	}()

	// Set the sheet name
	sheetName := "Sprint Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %w", err)
	}
	f.SetActiveSheet(index)

	// Set column headers
	headers := []string{
		"ID", "Title", "Type", "Size (hours)", "Status", "Description",
		"SprintNumber", "Responsible", "Column", "Team", "ParentIds",
		"Role", "IsBlocked", "TotalTime", "SizeUnit", "ChildSizeSum", "ChildSprintSum",
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "Q", 20)

	// Write report summary
	f.SetCellValue(sheetName, "A3", "Sprint ID:")
	f.SetCellValue(sheetName, "B3", report.SprintID)
	f.SetCellValue(sheetName, "A4", "Responsible:")
	f.SetCellValue(sheetName, "B4", report.Responsible)
	f.SetCellValue(sheetName, "A5", "Total Tasks:")
	f.SetCellValue(sheetName, "B5", report.TotalTasks)
	f.SetCellValue(sheetName, "A6", "Total Hours:")
	f.SetCellValue(sheetName, "B6", report.TotalHours)

	// Write tasks data
	for i, task := range report.Tasks {
		row := i + 8 // Start from row 8 to leave space for summary

		// Format parent IDs
		parentIDsStr := ""
		if len(task.ParentIDs) > 0 {
			parts := make([]string, len(task.ParentIDs))
			for idx, v := range task.ParentIDs {
				parts[idx] = fmt.Sprintf("%d", v)
			}
			parentIDsStr = strings.Join(parts, ", ")
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), task.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), task.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), task.Type)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), task.Size)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), task.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), task.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), task.SprintNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), task.Responsible)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), task.Column)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), task.Team)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), parentIDsStr)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), task.Role)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), task.IsBlocked)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), task.TotalTime)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), task.SizeUnit)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), task.ChildSizeSum)
		f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), task.ChildSprintSum)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("sprint_report_%d_%s.xlsx", report.SprintID, timestamp)

	// Save the file
	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	slog.Info("Report saved to Excel file", "filename", filename)
	return nil
}

func AsyncProcessReport(env config.Config, token string, kaitenURL string, sprintID int, userEmail string, report *models.Report) error {
	if report == nil {
		return fmt.Errorf("report is nil")
	}

	client := api.CreateKaitenClient(token, kaitenURL)

	report, err := client.GetSprintReport(sprintID, userEmail)
	if err != nil {
		return fmt.Errorf("report generation error: %w", err)
	}

	// Print report summary to console
	fmt.Println("\nSprint Report Summary:")
	fmt.Printf("Sprint ID: %d\n", report.SprintID)
	fmt.Printf("Responsible: %s\n", report.Responsible)
	fmt.Printf("Total Tasks: %d\n", report.TotalTasks)
	fmt.Printf("Total Hours: %d\n", report.TotalHours)
	fmt.Println("\nTasks:")

	// Print individual tasks to console
	for i, task := range report.Tasks {
		fmt.Printf("\nTask #%d:\n", i+1)
		fmt.Printf("  ID: %d\n", task.ID)
		fmt.Printf("  Title: %s\n", task.Title)
		fmt.Printf("  Type: %s\n", task.Type)
		fmt.Printf("  Size: %d hours\n", task.Size)
		fmt.Printf("  Status: %s\n", task.Status)
		if task.Description != "" {
			fmt.Printf("  Description: %s\n", task.Description)
		}
	}

	// Save report to Excel file
	if err := saveReportToExcel(report); err != nil {
		return fmt.Errorf("failed to save report to Excel: %w", err)
	}

	return nil
}
