package DbQueryStrategies

// PostgresRowLockingStrategy for PostgreSQL operations with explicit row locking
type PostgresRowLockingStrategy struct {
	BasePostgresStrategy
}

// ONLY ROW LOCKING ON WRITES TO DB
// THIS WILL NOT BLOCK ANY READS TO ROWS BEEN MODIFIED BY A CONCURRENT PROCESS USING A STANDARD SELECT
// THIS ROW LOCKING STRATEGY IMPL IS ALSO USING STANDARD SELECT, ONLY ROW LOCKING FOR MODIFICATIONS
// STALE READS COULD OCCUR if a row locked and modified by an update operation, is read using standard select
// in Read Committed Isolation level in PostgresSql(default):
/*
 Transaction A locks and updates a row.
 Transaction B reads the row while Transaction A still holds the lock but has not yet committed.
 If Transaction A commits or rolls back, Transaction B may have seen stale or inconsistent data that
	is no longer valid after the commit.

	Problem with ReadCommitted isolation level is non-repeatable reads,
		same row has different values when read at different points of time in the same transaction.
*/

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
