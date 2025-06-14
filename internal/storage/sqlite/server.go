package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/krasilnikovm/logman/internal/entity"
	_ "github.com/mattn/go-sqlite3"
)

const DriverName = "sqlite3"

// A ServerStorage contains methods for communication with Server entity
type ServerStorage struct {
	connStr string
}

func NewServerStorage(connStr string) *ServerStorage {
	return &ServerStorage{
		connStr: connStr,
	}
}

// A Create method creates new Server in database
func (s *ServerStorage) Create(ctx context.Context, server *entity.Server) error {
	db, err := sql.Open(DriverName, s.connStr)

	if err != nil {
		return fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"INSERT INTO servers (name, host, log_location_path, log_location_format, credential_id, created_at, updated_at) VALUES(?,?,?,?,?,?,?)",
	)

	if err != nil {
		return fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		server.Name,
		server.Host,
		server.LogFolderPath,
		server.LogFormat,
		server.CredentialId,
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
	db, err := sql.Open(DriverName, s.connStr)

	if err != nil {
		return nil, fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"SELECT id, name, host, log_location_path, log_location_format, credential_id, created_at, updated_at FROM servers WHERE id = ?;",
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
		&server.LogFolderPath,
		&server.LogFormat,
		&server.CredentialId,
		&server.CreatedAt,
		&server.UpdatedAt,
	)

	return &server, nil
}

// A DeleteById method deletes Server by id
func (s *ServerStorage) DeleteById(ctx context.Context, id int) error {
	db, err := sql.Open(DriverName, s.connStr)

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

// A GetList method returns list of Servers
func (s *ServerStorage) GetList(ctx context.Context, limit, page int) ([]entity.Server, error) {
	var servers []entity.Server

	if limit < 0 || page < 0 {
		return servers, fmt.Errorf("invalid input parameters")
	}

	db, err := sql.Open(DriverName, s.connStr)

	if err != nil {
		return servers, fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"SELECT id, name, host, log_location_path, log_location_format, credential_id, created_at, updated_at FROM servers ORDER BY id DESC LIMIT ? OFFSET ?;",
	)

	if err != nil {
		return servers, fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, limit, (page-1)*limit)

	if err != nil {
		return servers, fmt.Errorf("query execution failed: %w", err)
	}

	defer rows.Close()

	var server entity.Server

	for rows.Next() {
		rows.Scan(
			&server.Id,
			&server.Name,
			&server.Host,
			&server.LogFolderPath,
			&server.LogFormat,
			&server.CredentialId,
			&server.CreatedAt,
			&server.UpdatedAt,
		)

		servers = append(servers, server)
	}

	return servers, nil
}

// A Update method updates Server by id
func (s *ServerStorage) Update(ctx context.Context, server *entity.Server, id int) error {
	db, err := sql.Open(DriverName, s.connStr)

	if err != nil {
		return fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.PrepareContext(
		ctx,
		"UPDATE servers SET name = ?, host = ?, log_location_path = ?, log_location_format = ?, credential_id = ?, updated_at = ? WHERE id = ?;",
	)

	if err != nil {
		return fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		server.Name,
		server.Host,
		server.LogFolderPath,
		server.LogFormat,
		server.CredentialId,
		server.UpdatedAt,
		id,
	)

	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	return nil

}
