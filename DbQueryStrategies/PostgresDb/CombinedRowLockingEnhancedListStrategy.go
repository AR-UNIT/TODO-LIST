package PostgresDb

import (
	"TODO-LIST/commons"
)

// CombinedRowLockingEnhancedListStrategy combines both row locking and enhanced listing strategies
type CombinedRowLockingEnhancedListStrategy struct {
	RowLockingStrategy   RowLockingStrategy
	EnhancedListStrategy EnhancedListStrategy
}

func (crls *CombinedRowLockingEnhancedListStrategy) AddTask(task commons.Task) (int, error) {
	return crls.RowLockingStrategy.AddTask(task)
}

func (crls *CombinedRowLockingEnhancedListStrategy) ListTasks(params map[string]interface{}) ([]commons.Task, error) {
	return crls.EnhancedListStrategy.ListTasks(params)
}

func (crls *CombinedRowLockingEnhancedListStrategy) CompleteTask(id int) (int64, error) {
	return crls.RowLockingStrategy.CompleteTask(id)
}

func (crls *CombinedRowLockingEnhancedListStrategy) DeleteTask(id int) (int64, error) {
	return crls.RowLockingStrategy.DeleteTask(id)
}
