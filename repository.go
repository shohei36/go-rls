package main

import (
	"context"
	"database/sql"
)

type userRepository struct {
	db TenantDB
}

func NewUserRepository(db TenantDB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FetchByID(ctx context.Context, tenantID, userID string) (*User, error) {
	tx, ok := GetTx(ctx)
	if !ok {
		conn, err := r.db.Conn(ctx, tenantID)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		tx = conn
	}
	query := "select id, name, gender, age from users where id = $1"
	row := tx.QueryRowContext(ctx, query, userID)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Gender,
		&user.Age,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, tenantID string, u User) (rowsAffected int64, err error) {
	tx, ok := GetTx(ctx)
	if !ok {
		conn, err := r.db.Conn(ctx, tenantID)
		if err != nil {
			return 0, err
		}
		defer conn.Close()
		tx = conn
	}
	query := "update users set name = $1, gender = $2, age = $3 where id = $4"
	result, err := tx.ExecContext(ctx, query, u.Name, u.Gender, u.Age, u.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
