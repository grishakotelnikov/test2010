package repository

import (
	"context"
	"database/sql"
	"fmt"
	"project1/internal/models"
)

type Repository interface {
	Deposit(ctx context.Context, id int, amount int) error
	Transfer(ctx context.Context, fromId int, toId int, amount int) error
	GetLastTransactions(ctx context.Context, userId int) ([]models.Transaction, error)
}

type PostgresRep struct {
	db *sql.DB
}

func NewPostgresRep(db *sql.DB) *PostgresRep {
	return &PostgresRep{db: db}
}

func (pg *PostgresRep) PopUsersIfEmpty(ctx context.Context) error {
	var count int
	err := pg.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("ошибка проверки countUsers: %v", err)
	}

	if count > 0 {
		return nil
	}

	for i := 0; i < 100; i++ {
		firstName := fmt.Sprintf("User%d", i+1)
		lastName := fmt.Sprintf("Lastname%d", i+1)
		balance := 10000

		_, err := pg.db.ExecContext(ctx, `
			INSERT INTO users (first_name, last_name, balance)
			VALUES ($1, $2, $3)
		`, firstName, lastName, balance)
		if err != nil {
			return fmt.Errorf("ошибка вставки пользователя %d: %v", i+1, err)
		}
	}

	return nil
}

func (pg *PostgresRep) Deposit(ctx context.Context, id int, amount int) error {
	var exists bool
	err := pg.db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка про проверке сущетсвования пользователя: %v", err)
	}
	if !exists {
		return fmt.Errorf("пользователь с id %d не существует", id)
	}

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (user_id, type, amount) VALUES ($1, 'deposit', $2)", id, amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg *PostgresRep) Transfer(ctx context.Context, fromId int, toId int, amount int) error {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", fromId).Scan(&exists)
	if err != nil || !exists {
		tx.Rollback()
		return fmt.Errorf("пользователь с id %d не существует", fromId)
	}

	err = tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", toId).Scan(&exists)
	if err != nil || !exists {
		tx.Rollback()
		return fmt.Errorf("пользователь с id %d не существует", toId)
	}

	var balance int
	err = tx.QueryRowContext(ctx, "SELECT balance FROM users WHERE id = $1", fromId).Scan(&balance)
	if err != nil {
		tx.Rollback()
		return err
	}

	if balance < amount {
		tx.Rollback()
		return fmt.Errorf("недостаточно средств")
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", amount, fromId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (user_id, type, amount, to_id) VALUES ($1, 'transfer', $2, $3)", fromId, amount, toId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO transactions (user_id, type, amount, from_id) VALUES ($1, 'receive', $2, $3)", toId, amount, fromId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg *PostgresRep) GetLastTransactions(ctx context.Context, userId int) ([]models.Transaction, error) {
	var exists bool
	err := pg.db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", userId).Scan(&exists)
	if err != nil || !exists {
		return nil, fmt.Errorf("пользователь с id %d не существует", userId)
	}

	rows, err := pg.db.QueryContext(ctx, `
		SELECT id, user_id, type, amount, from_id, to_id, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 10
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.UserId, &t.Type, &t.Amount, &t.FromId, &t.ToId, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
