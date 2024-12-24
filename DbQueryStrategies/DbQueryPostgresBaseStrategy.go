package DbQueryStrategies

import (
	"TODO-LIST/commons"
	"database/sql"
	"fmt"
)

type BasePostgresStrategy struct {
	Db *sql.DB
}

// BasePostgresStrategy handles common database operations

func (bps *BasePostgresStrategy) AddTask(task commons.Task) (int, error) {
	query := `INSERT INTO TODO.tasks (description, completed) VALUES ($1, $2) RETURNING id`
	var id int
	err := bps.Db.QueryRow(query, task.Description, false).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// TODO: params is unsused in base strategy, need to update
func (bps *BasePostgresStrategy) ListTasks(params map[string]interface{}) ([]commons.Task, error) {
	rows, err := bps.Db.Query(`SELECT id, description, completed FROM TODO.tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (bps *BasePostgresStrategy) CompleteTask(id int) (int64, error) {
	query := `UPDATE TODO.tasks SET completed = true WHERE id = $1`
	res, err := bps.Db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return 0, fmt.Errorf("task not found")
	}
	return rowsAffected, nil
}

func (bps *BasePostgresStrategy) DeleteTask(id int) (int64, error) {
	query := `DELETE FROM TODO.tasks WHERE id = $1`
	res, err := bps.Db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return 0, fmt.Errorf("task not found")
	}
	return rowsAffected, nil
}
