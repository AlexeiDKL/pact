package basedate

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"dkl.ru/pact/bd_service/iternal/config"
	"dkl.ru/pact/bd_service/iternal/logger"
	// "dkl.ru/pact/bd_service/iternal/core"
	// Remove the import of logger to avoid import cycle
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
		logger.Logger.Info("База отсутствует, создаем...")
		if _, err := db.Exec("CREATE DATABASE " + d.DbName); err != nil {
			return fmt.Errorf("создание БД: %w", err)
		}
		logger.Logger.Info("База данных создана.")
	}
	return nil
}

func (d *Database) GetLatestVersionByLanguage(ctx context.Context, language string) (string, error) {
	// select left join между таблицами languages и versions по language_id,
	// где short_name = language order by version.version
	// чтобы получить последнюю версию для указанного языка
	if language == "" {
		return "", fmt.Errorf("язык не указан")
	}
	query := `
		SELECT v.version
		FROM languages l
		LEFT JOIN versions v ON l.id = v.language_id
		WHERE l.short_name = $1
		ORDER BY v.version DESC
		LIMIT 1;
	`
	var version string
	if err := d.DB.QueryRowContext(ctx, query, language).Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("нет версий для языка %s", language)
		}
		return "", fmt.Errorf("ошибка запроса последней версии для языка %s: %w", language, err)
	}
	logger.Logger.Debug("Последняя версия для языка", "language", language, "version", version)
	return version, nil
}

type ResultGetPathFtCLastVersionByLanguage struct {
	Full    string
	Content string
}

func (d *Database) GetPathFileClasVersionByLanguage(ctx context.Context, language string) (ResultGetPathFtCLastVersionByLanguage, error) {
	var res ResultGetPathFtCLastVersionByLanguage

	if language == "" {
		return res, fmt.Errorf("язык не указан")
	}
	query := `
    SELECT
        ft.file_path AS full_text_path,
        ct.file_path AS contents_path
    FROM version AS v
        JOIN language AS l ON v.language_id = l.id
        JOIN file AS ft ON v.full_text_id = ft.id
        JOIN file AS ct ON v.contents_id = ct.id
    WHERE l.short_name = $1
    ORDER BY v.version DESC
    LIMIT 1;
    `
	if err := d.DB.QueryRowContext(ctx, query, language).Scan(&res.Full, &res.Content); err != nil {
		if err == sql.ErrNoRows {
			return res, fmt.Errorf("нет файлов для языка %s", language)
		}
		return res, fmt.Errorf("ошибка запроса пути последней версии для языка %s: %w", language, err)
	}

	logger.Logger.Debug("Путь к последней версии для языка", "language", language, "file_path", res.Full, "content_path", res.Content)
	return res, nil
}

// Пакетные типы для запросов/ответов не нужны — это утилита внутри Database.
func (d *Database) CheckUpdates(ctx context.Context, language, currentVersion string) (bool, error) {
	// 1. Логируем структурировано
	logger.Logger.Debug("CheckUpdates start",
		"language", language,
		"current_version", currentVersion,
	)

	// 2. Получаем latestVersion из БД
	latestVersion, err := d.GetLatestVersionByLanguage(ctx, language)
	if err != nil {
		logger.Logger.Error("DB error on GetLatestVersionByLanguage",
			"language", language,
			"error", err,
		)
		return false, err
	}
	if latestVersion == "" {
		// Если версий нет, просто обновлять нечего
		logger.Logger.Info("No versions found for language", "language", language)
		return false, nil
	}
	// 3. Сравниваем версии
	bdVersion, err := ParseVersion(latestVersion)
	if err != nil {
		logger.Logger.Error("Error parsing version",
			"language", language,
			"version", latestVersion,
			"error", err,
		)
		return false, fmt.Errorf("ошибка парсинга версии %s для языка %s: %w", latestVersion, language, err)
	}
	currentVersioni, err := ParseVersion(currentVersion)
	if err != nil {
		logger.Logger.Error("Error parsing current version",
			"language", language,
			"version", currentVersion,
			"error", err,
		)
		return false, fmt.Errorf("ошибка парсинга текущей версии %s для языка%s: %w", currentVersion, language, err)
	}
	if bdVersion > currentVersioni {
		logger.Logger.Info(fmt.Sprintf("Доступно обновление для языка %s: текущая версия %s, новая версия %d", language, currentVersion, bdVersion))
		return true, nil
	}
	return false, nil
}

func ParseVersion(version string) (int, error) {
	// версия это timestamp
	// например 1753026232
	// преобразуем в int
	if version == "" {
		return 0, fmt.Errorf("пустая версия")
	}
	var ver int
	_, err := fmt.Sscanf(version, "%d", &ver)
	if err != nil {
		return 0, fmt.Errorf("неверный формат версии: %s", version)
	}
	if ver <= 0 {
		return 0, fmt.Errorf("версия должна быть положительным числом: %d", ver)
	}
	return ver, nil
}
