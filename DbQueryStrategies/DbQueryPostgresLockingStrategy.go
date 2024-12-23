package DbQueryStrategies

// PostgresRowLockingStrategy for PostgreSQL operations with explicit row locking
type PostgresRowLockingStrategy struct {
	BasePostgresStrategy
}

func (prls *PostgresRowLockingStrategy) CompleteTask(id int) (int64, error) {
	// Add row locking before completing the task
	lockQuery := `SELECT id FROM TODO.tasks WHERE id = $1 FOR UPDATE`
	_, err := prls.Db.Exec(lockQuery, id)
	if err != nil {
		return 0, err
	}

	// Use the base method for the actual update
	return prls.BasePostgresStrategy.CompleteTask(id)
}

func (prls *PostgresRowLockingStrategy) DeleteTask(id int) (int64, error) {
	// Add row locking before deleting the task
	lockQuery := `SELECT id FROM TODO.tasks WHERE id = $1 FOR UPDATE`
	_, err := prls.Db.Exec(lockQuery, id)
	if err != nil {
		return 0, err
	}

	// Use the base method for the actual delete
	return prls.BasePostgresStrategy.DeleteTask(id)
}
