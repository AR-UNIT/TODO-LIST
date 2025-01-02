package PostgresDb

import (
	"TODO-LIST/commons"
	"database/sql"
	"fmt"
	"strings"
)

type EnhancedListStrategy struct {
	DbContext
}

// ListTasks implements pagination, filtering, and sorting using a map of parameters
func (pes *EnhancedListStrategy) ListTasks(params map[string]interface{}) ([]commons.Task, error) {
	// Extract parameters using helper functions
	page := extractIntParam(params, "page", 1)
	limit := extractIntParam(params, "limit", 10)
	status := extractStringParam(params, "status", "")
	sort := extractStringParam(params, "sort", "")

	// Base query
	query := `SELECT id, description, completed FROM TODO.tasks`
	countQuery := `SELECT COUNT(*) FROM TODO.tasks`

	// Build filters
	filters, args := buildFilters(status)
	if len(filters) > 0 {
		filterString := " WHERE " + strings.Join(filters, " AND ")
		query += filterString
		countQuery += filterString
	}

	// Add sorting
	query += buildSortClause(sort)

	// Add pagination
	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	// Execute the query to retrieve tasks
	rows, err := pes.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the results into a slice of tasks
	return parseTasks(rows)
}

// Helper function to extract integer parameters
func extractIntParam(params map[string]interface{}, key string, defaultValue int) int {
	if val, ok := params[key].(int); ok {
		return val
	}
	return defaultValue
}

// Helper function to extract string parameters
func extractStringParam(params map[string]interface{}, key string, defaultValue string) string {
	if val, ok := params[key].(string); ok {
		return val
	}
	return defaultValue
}

// Helper function to build filters
func buildFilters(status string) ([]string, []interface{}) {
	filters := []string{}
	args := []interface{}{}

	if status != "" {
		switch status {
		case "completed":
			filters = append(filters, fmt.Sprintf("completed = $%d", len(args)+1))
			args = append(args, true)
		case "pending":
			filters = append(filters, fmt.Sprintf("completed = $%d", len(args)+1))
			args = append(args, false)
		}
	}
	return filters, args
}

// Helper function to build sort clause
func buildSortClause(sort string) string {
	if sort == "" {
		return ""
	}

	sortParts := strings.Split(sort, ":")
	column := sortParts[0]
	order := "ASC"
	if len(sortParts) > 1 && strings.ToUpper(sortParts[1]) == "DESC" {
		order = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", column, order)
}

// Helper function to parse tasks from rows
func parseTasks(rows *sql.Rows) ([]commons.Task, error) {
	tasks := []commons.Task{}
	for rows.Next() {
		var task commons.Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
