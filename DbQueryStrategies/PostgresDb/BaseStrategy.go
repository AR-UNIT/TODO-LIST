package PostgresDb

import (
	"TODO-LIST/commons"
	"database/sql"
	"fmt"
)

type DbContext struct {
	Db *sql.DB
}

// DbContext handles default database operations
type DefaultPostgresStrategy struct {
	DbContext
}

func (bps *DbContext) AddTask(taskInput *commons.TaskInputModel) (int, error) {
	query := `INSERT INTO TODO.tasks (description, completed) VALUES ($1, $2) RETURNING id`
	var id int
	fmt.Println("before applying insert task query")
	err := bps.Db.QueryRow(query, taskInput.Description, false).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// TODO: params is unsused in base strategy, need to update
func (bps *DbContext) ListTasks(params map[string]interface{}) ([]commons.Task, error) {
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

func (bps *DbContext) CompleteTask(id int) (int64, error) {
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

func (bps *DbContext) DeleteTask(id int) (int64, error) {
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
