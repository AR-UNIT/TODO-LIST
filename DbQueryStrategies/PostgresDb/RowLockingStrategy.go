package PostgresDb

// RowLockingStrategy for PostgreSQL operations with explicit row locking
type RowLockingStrategy struct {
	DbContext
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
/*
	IF WE HAVE ROW LOCKING, WE NEED TO HAVE RETRY, SO THAT DEADLOCK SCENARIOS CAN BE HANDLED
*/

func (prls *RowLockingStrategy) CompleteTask(id int) (int64, error) {
	// Add row locking before completing the task
	lockQuery := `SELECT id FROM TODO.tasks WHERE id = $1 FOR UPDATE`
	_, err := prls.Db.Exec(lockQuery, id)
	if err != nil {
		return 0, err
	}

	// Use the base method for the actual update
	return prls.DbContext.CompleteTask(id)
}

func (prls *RowLockingStrategy) DeleteTask(id int) (int64, error) {
	// Add row locking before deleting the task
	lockQuery := `SELECT id FROM TODO.tasks WHERE id = $1 FOR UPDATE`
	_, err := prls.Db.Exec(lockQuery, id)
	if err != nil {
		return 0, err
	}

	// Use the base method for the actual delete
	return prls.DbContext.DeleteTask(id)
}
