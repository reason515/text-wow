package repository

import (
	"database/sql"
	"fmt"

	"text-wow/internal/database"
)

// TransactionFunc 事务函数类型
type TransactionFunc func(*sql.Tx) error

// WithTransaction 执行事务
// 如果函数返回错误，事务会自动回滚
func WithTransaction(fn TransactionFunc) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTransactionResult 执行事务并返回结果
func WithTransactionResult[T any](fn func(*sql.Tx) (T, error)) (T, error) {
	var zero T
	tx, err := database.DB.Begin()
	if err != nil {
		return zero, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	result, err := fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return zero, fmt.Errorf("transaction error: %w, rollback error: %v", err, rbErr)
		}
		return zero, err
	}

	if err := tx.Commit(); err != nil {
		return zero, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}



















































