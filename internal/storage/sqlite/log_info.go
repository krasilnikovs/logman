package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/krasilnikovm/logman/internal/entity"
)

type LogInfoStorage struct {
	dbPath string
}

func NewLogInfoStorage(dbPath string) *LogInfoStorage {
	return &LogInfoStorage{
		dbPath: dbPath,
	}
}

func (l *LogInfoStorage) Create(logInfo entity.LogInfo) error {
	db, err := sql.Open(driverName, l.dbPath)

	if err != nil {
		return fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO log_infos (name, location, format, created_at, updated_at) VALUES(?,?,?,?,?)")

	if err != nil {
		return fmt.Errorf("error during preparing query: %w", err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(logInfo.Name, logInfo.Location, logInfo.Format, logInfo.CreatedAt, logInfo.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error during executing query: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("can not fetch last insert id: %w", err)
	}

	logInfo.Id = int(id)

	return nil
}
