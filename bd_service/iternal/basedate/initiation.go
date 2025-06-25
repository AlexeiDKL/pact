package basedate

import (
	"database/sql"
	"fmt"

	"dkl.ru/pact/bd_service/iternal/config"
	"dkl.ru/pact/bd_service/iternal/logger"
	_ "github.com/lib/pq"
)

var (
	host            = config.Config.Bd_server.Host
	port            = config.Config.Bd_server.Port
	user            = config.Config.Bd_server.User
	password        = config.Config.Bd_server.Password
	dbname          = config.Config.Bd_server.Dbname
	sql_path        = config.Config.Bd_server.Path
	sql_delete_path = config.Config.Bd_server.DeletePath
)

var DB *sql.DB

/*
	инициализация
	проверяем наличие бд
	проверяем наличие таблиц
	инициализируем и возвращаем ссылку на бд
*/

func initVar() {
	host = config.Config.Bd_server.Host
	port = config.Config.Bd_server.Port
	user = config.Config.Bd_server.User
	password = config.Config.Bd_server.Password
	dbname = config.Config.Bd_server.Dbname
	sql_path = config.Config.Bd_server.Path
	sql_delete_path = config.Config.Bd_server.DeletePath
}

func Init() error {
	initVar()

	// 1. Проверяем существование базы
	if err := EnsureDatabaseExists(); err != nil {
		return err
	}

	// 2. Подключаемся к БД
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	// 3. Проверяем структуру таблиц, читаем из файла

	// if err := ExecuteDropTablesSQL(db, sql_delete_path); err != nil {
	// 	return fmt.Errorf("не удалось дропнуть Таблицы: %v", err)
	// }
	// return nil

	exist, _ := valideTable(db)

	if exist {
		return nil
	}

	if err := ExecuteSQLFile(db, sql_path); err != nil {
		return err
	}

	// if err := ExecuteDropTablesSQL(db, sql_delete_path); err != nil {
	// 	return fmt.Errorf("не удалось дропнуть Таблицы: %v", err)
	// }
	DB = db
	return nil
}

// Функция проверки существования БД
func EnsureDatabaseExists() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close()

	// Проверяем, существует ли наша БД
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s');", config.Config.Bd_server.Dbname)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}

	if !exists {
		logger.Logger.Info("База данных отсутствует, создаем...")
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", config.Config.Bd_server.Dbname))
		if err != nil {
			return fmt.Errorf("ошибка создания базы данных: %v", err)
		}
		logger.Logger.Info("База данных успешно создана.")
	} else {
		logger.Logger.Info("База данных существует.")
	}
	return nil
}
