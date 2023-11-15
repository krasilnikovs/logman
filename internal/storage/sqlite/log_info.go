package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/krasilnikovm/logman/internal/entity"
)

// A LogInfoStorage contains methods for communication with LogInfo entity
type LogInfoStorage struct {
	dbPath string
}

func NewLogInfoStorage(dbPath string) *LogInfoStorage {
	return &LogInfoStorage{
		dbPath: dbPath,
	}
}

func (l *LogInfoStorage) Create(ctx context.Context, logInfo entity.LogInfo) error {
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

	result, err := stmt.ExecContext(ctx, logInfo.Name, logInfo.Location, logInfo.Format, logInfo.CreatedAt, logInfo.UpdatedAt)

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

// A GetById method return LogInfo if no errors
// In case when LogInfo is not found the method will return nil
func (l *LogInfoStorage) GetById(ctx context.Context, id int) (*entity.LogInfo, error) {
	db, err := sql.Open(driverName, l.dbPath)

	if err != nil {
		return nil, fmt.Errorf("can not open sqlite connection: %w", err)
	}

	defer db.Close()

	stmt, err := db.Prepare("SELECT id, name, location, format, created_at, updated_at FROM log_infos WHERE id=?")

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

	var logInfo entity.LogInfo

	rows.Scan(
		&logInfo.Id,
		&logInfo.Name,
		&logInfo.Location,
		&logInfo.Format,
		&logInfo.CreatedAt,
		&logInfo.UpdatedAt,
	)

	return &logInfo, nil
}
