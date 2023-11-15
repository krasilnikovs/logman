package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/krasilnikovm/logman/internal/entity"
	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"

// A ServerStorage contains methods for communication with Server, LogInfo entities
type ServerStorage struct {
	dbPath string
}

func NewServerStorage(dbPath string) *ServerStorage {
	return &ServerStorage{
		dbPath: dbPath,
	}
}

// Create method create server in database
func (s *ServerStorage) Create(server *entity.Server) error {
	db, err := sql.Open(driverName, s.dbPath)

	if err != nil {
		return fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO servers (log_infos_id, name, host, created_at, updated_at) VALUES(?,?,?,?)")

	if err != nil {
		return fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(server.LogInfo.Id, server.Name, server.Host, server.CreatedAt, server.UpdatedAt)

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
