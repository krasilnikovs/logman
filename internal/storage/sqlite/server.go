package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/krasilnikovm/logman/internal/entity"
	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"

// A ServerStorage contains methods for communication with Server entity
type ServerStorage struct {
	connStr string
}

func NewServerStorage(connStr string) *ServerStorage {
	return &ServerStorage{
		connStr: connStr,
	}
}

func (s *ServerStorage) Create(ctx context.Context, server *entity.Server) error {
	db, err := sql.Open(driverName, s.connStr)

	if err != nil {
		return fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"INSERT INTO servers (name, host, log_location_path, log_location_format, created_at, updated_at) VALUES(?,?,?,?,?,?)",
	)

	if err != nil {
		return fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		server.Name,
		server.Host,
		server.LogLocation.Path,
		server.LogLocation.Format,
		server.CreatedAt,
		server.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error during executing query: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("can not fetch last insert id: %w", err)
	}

	server.Id = int(id)

	return nil
}

// A GetById method return Server if no errors
// In case when Server is not found the method will return nil
func (s *ServerStorage) GetById(ctx context.Context, id int) (*entity.Server, error) {
	db, err := sql.Open(driverName, s.connStr)

	if err != nil {
		return nil, fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"SELECT id, name, host, log_location_path, log_location_format, created_at, updated_at FROM servers WHERE id = ?;",
	)

	if err != nil {
		return nil, fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var server entity.Server

	rows.Scan(
		&server.Id,
		&server.Name,
		&server.Host,
		&server.LogLocation.Path,
		&server.LogLocation.Format,
		&server.CreatedAt,
		&server.UpdatedAt,
	)

	return &server, nil
}

func (s *ServerStorage) DeleteById(ctx context.Context, id int) error {
	db, err := sql.Open(driverName, s.connStr)

	if err != nil {
		return fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"DELETE FROM servers WHERE id = ?;",
	)

	if err != nil {
		return fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	return nil
}
