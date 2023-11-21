package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/krasilnikovm/logman/internal/entity"
)

type CredentialStorage struct {
	connStr string
}

func NewCredentialStorage(connStr string) *CredentialStorage {
	return &CredentialStorage{
		connStr: connStr,
	}
}

func (c *CredentialStorage) Create(ctx context.Context, credential *entity.Credential) error {
	db, err := sql.Open(DriverName, c.connStr)

	if err != nil {
		return fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"INSERT INTO credentials (name, path, created_at, updated_at) VALUES(?, ?,?,?)",
	)

	if err != nil {
		return fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		credential.Name,
		credential.Path,
		credential.CreatedAt,
		credential.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error during executing query: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("error during getting last insert id: %w", err)
	}

	credential.Id = int(id)

	return nil
}

func (c *CredentialStorage) GetById(ctx context.Context, id int) (*entity.Credential, error) {
	db, err := sql.Open(DriverName, c.connStr)

	if err != nil {
		return nil, fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"SELECT id, name, path, created_at, updated_at FROM credentials WHERE id = ?",
	)

	if err != nil {
		return nil, fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	var credential entity.Credential

	err = row.Scan(
		&credential.Id,
		&credential.Name,
		&credential.Path,
		&credential.CreatedAt,
		&credential.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error during scanning row: %w", err)
	}

	return &credential, nil
}

func (c *CredentialStorage) GetList(ctx context.Context, page, limit int) ([]*entity.Credential, error) {
	db, err := sql.Open(DriverName, c.connStr)

	if err != nil {
		return nil, fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	rows, err := db.QueryContext(
		ctx,
		"SELECT id, name, path, created_at, updated_at FROM credentials ORDER BY id DESC LIMIT ? OFFSET ?;",
		limit,
		(page-1)*limit,
	)

	if err != nil {
		return nil, fmt.Errorf("error during executing query: %w", err)
	}

	defer rows.Close()

	credentials := make([]*entity.Credential, 0, limit)

	for rows.Next() {
		credential := &entity.Credential{}

		err := rows.Scan(
			&credential.Id,
			&credential.Name,
			&credential.Path,
			&credential.CreatedAt,
			&credential.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error during scanning row: %w", err)
		}

		credentials = append(credentials, credential)
	}

	return credentials, nil
}

func (c *CredentialStorage) DeleteById(ctx context.Context, id int) error {
	db, err := sql.Open(DriverName, c.connStr)

	if err != nil {
		return fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"DELETE FROM credentials WHERE id = ?;",
	)

	if err != nil {
		return fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)

	if err != nil {
		return fmt.Errorf("error during executing query: %w", err)
	}

	return nil
}

func (c *CredentialStorage) Update(ctx context.Context, credential *entity.Credential) error {
	db, err := sql.Open(DriverName, c.connStr)

	if err != nil {
		return fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"UPDATE credentials SET name = ?, path = ?, updated_at = ? WHERE id = ?;",
	)

	if err != nil {
		return fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		credential.Name,
		credential.Path,
		credential.UpdatedAt,
		credential.Id,
	)

	if err != nil {
		return fmt.Errorf("error during executing query: %w", err)
	}

	return nil
}
