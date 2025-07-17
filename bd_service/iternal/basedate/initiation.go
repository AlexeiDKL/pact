package basedate

import (
	"database/sql"
	"fmt"
	"log/slog"

	"dkl.ru/pact/bd_service/iternal/config"
)

type Database struct {
	DB         *sql.DB
	Host       string
	Port       int
	User       string
	Password   string
	DbName     string
	SchemaPath string
	DropPath   string
	Logger     *slog.Logger
}

func New(cfg config.BDServer, log *slog.Logger) (*Database, error) {
	db := &Database{
		Host:       cfg.Host,
		Port:       cfg.Port,
		User:       cfg.User,
		Password:   cfg.Password,
		DbName:     cfg.Dbname,
		SchemaPath: cfg.Path,
		DropPath:   cfg.DeletePath,
		Logger:     log,
	}

	if err := db.ensureDatabaseExists(); err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.DbName)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	db.DB = conn

	// // Проверка таблиц и инициализация схемы
	// ok, err := db.validateTables()
	// if err != nil {
	// 	return nil, err
	// }
	// if !ok {
	// 	if err := ExecuteSQLFile(conn, db.SchemaPath); err != nil {
	// 		return nil, err
	// 	}
	// }

	return db, nil
}

func (d *Database) ensureDatabaseExists() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", d.Host, d.Port, d.User, d.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("подключение к postgres: %w", err)
	}
	defer db.Close()

	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s');", d.DbName)
	if err := db.QueryRow(query).Scan(&exists); err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}

	if !exists {
		d.Logger.Info("База отсутствует, создаем...")
		if _, err := db.Exec("CREATE DATABASE " + d.DbName); err != nil {
			return fmt.Errorf("создание БД: %w", err)
		}
		d.Logger.Info("База данных создана.")
	}
	return nil
}
