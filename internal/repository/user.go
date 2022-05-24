package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"weedy/internal/entity"
	"weedy/internal/helper"
)

type UserRepo struct {
	DB *sql.DB
}

type ListParams struct {
	UserIDs []string
	Email   string
}

func (r *UserRepo) List(ctx context.Context, params *ListParams) ([]*entity.User, error) {
	sqlFieldMap, err := helper.GetSQLFieldMapFromEntity(&entity.User{})
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT %s FROM users 
	`, strings.Join(sqlFieldMap.Keys, ","))

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		u := &entity.User{}
		sqlFieldMap, err := helper.GetSQLFieldMapFromEntity(u)
		if err != nil {
			return nil, err
		}

		err = rows.Scan(sqlFieldMap.Values...)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) Upsert(ctx context.Context, e *entity.User) error {
	sqlFieldMap, err := helper.GetSQLFieldMapFromEntity(e)
	if err != nil {
		return err
	}

	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now

	stmt := fmt.Sprintf(`
		INSERT INTO users (%s) VALUES (%s) 
		ON CONFLICT(user_id) DO 
		UPDATE SET last_name=$2,first_name=$3,email=$4,updated_at=now(),deleted_at=$5
	`, strings.Join(sqlFieldMap.Keys, ","), helper.GeneratePlaceholder(len(sqlFieldMap.Keys)))
	_, err = r.DB.ExecContext(ctx, stmt, sqlFieldMap.Values...)
	return err
}
