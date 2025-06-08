package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"wallpaper-system/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к базе данных
	dsn := cfg.Database.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки подключения к базе данных: %v", err)
	}

	// Создание таблицы миграций если не существует
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("Ошибка создания таблицы миграций: %v", err)
	}

	switch command {
	case "up":
		if err := runMigrationsUp(db); err != nil {
			log.Fatalf("Ошибка выполнения миграций: %v", err)
		}
	case "down":
		if err := runMigrationsDown(db); err != nil {
			log.Fatalf("Ошибка отката миграций: %v", err)
		}
	case "status":
		if err := showMigrationStatus(db); err != nil {
			log.Fatalf("Ошибка получения статуса миграций: %v", err)
		}
	default:
		fmt.Printf("Неизвестная команда: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Использование:")
	fmt.Println("  go run cmd/migrate/main.go <команда>")
	fmt.Println("")
	fmt.Println("Команды:")
	fmt.Println("  up     - Выполнить все миграции")
	fmt.Println("  down   - Откатить последнюю миграцию")
	fmt.Println("  status - Показать статус миграций")
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func runMigrationsUp(db *sql.DB) error {
	// Получаем список выполненных миграций
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("ошибка получения списка выполненных миграций: %w", err)
	}

	// Получаем список файлов миграций
	migrationFiles, err := getMigrationFiles("up")
	if err != nil {
		return fmt.Errorf("ошибка получения файлов миграций: %w", err)
	}

	executed := 0
	for _, file := range migrationFiles {
		version := extractVersionFromFilename(file)

		// Проверяем, была ли уже выполнена эта миграция
		if appliedMigrations[version] {
			log.Printf("Миграция %s уже выполнена, пропускаем", version)
			continue
		}

		log.Printf("Выполнение миграции: %s", file)

		// Читаем содержимое файла миграции
		content, err := os.ReadFile(filepath.Join("migrations", file))
		if err != nil {
			return fmt.Errorf("ошибка чтения файла миграции %s: %w", file, err)
		}

		// Выполняем миграцию в транзакции
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("ошибка начала транзакции: %w", err)
		}

		// Выполняем SQL из файла миграции
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка выполнения миграции %s: %w", file, err)
		}

		// Записываем информацию о выполненной миграции
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка записи статуса миграции %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("ошибка подтверждения транзакции для миграции %s: %w", file, err)
		}

		log.Printf("Миграция %s выполнена успешно", version)
		executed++
	}

	if executed == 0 {
		log.Println("Новых миграций для выполнения нет")
	} else {
		log.Printf("Выполнено миграций: %d", executed)
	}

	return nil
}

func runMigrationsDown(db *sql.DB) error {
	// Получаем последнюю выполненную миграцию
	lastVersion, err := getLastAppliedMigration(db)
	if err != nil {
		return fmt.Errorf("ошибка получения последней миграции: %w", err)
	}

	if lastVersion == "" {
		log.Println("Нет миграций для отката")
		return nil
	}

	// Ищем соответствующий down файл
	downFile := ""
	files, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("ошибка чтения директории миграций: %w", err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), lastVersion) && strings.Contains(file.Name(), ".down.sql") {
			downFile = file.Name()
			break
		}
	}

	if downFile == "" {
		return fmt.Errorf("файл отката для миграции %s не найден", lastVersion)
	}

	log.Printf("Откат миграции: %s", downFile)

	// Читаем содержимое файла отката
	content, err := os.ReadFile(filepath.Join("migrations", downFile))
	if err != nil {
		return fmt.Errorf("ошибка чтения файла отката %s: %w", downFile, err)
	}

	// Выполняем откат в транзакции
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	// Выполняем SQL отката
	if _, err := tx.Exec(string(content)); err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка выполнения отката %s: %w", downFile, err)
	}

	// Удаляем запись о миграции
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", lastVersion); err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка удаления записи миграции %s: %w", lastVersion, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции отката %s: %w", downFile, err)
	}

	log.Printf("Откат миграции %s выполнен успешно", lastVersion)
	return nil
}

func showMigrationStatus(db *sql.DB) error {
	// Получаем список выполненных миграций
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("ошибка получения списка выполненных миграций: %w", err)
	}

	// Получаем список всех файлов миграций
	allMigrations, err := getMigrationFiles("up")
	if err != nil {
		return fmt.Errorf("ошибка получения файлов миграций: %w", err)
	}

	fmt.Println("Статус миграций:")
	fmt.Println("================")

	for _, file := range allMigrations {
		version := extractVersionFromFilename(file)
		status := "НЕ ВЫПОЛНЕНА"
		if appliedMigrations[version] {
			status = "ВЫПОЛНЕНА"
		}
		fmt.Printf("%-20s %s\n", version, status)
	}

	return nil
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	migrations := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		migrations[version] = true
	}

	return migrations, nil
}

func getLastAppliedMigration(db *sql.DB) (string, error) {
	var version string
	err := db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return version, err
}

func getMigrationFiles(direction string) ([]string, error) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), fmt.Sprintf(".%s.sql", direction)) {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Сортируем файлы по версии
	sort.Slice(migrationFiles, func(i, j int) bool {
		versionI := extractVersionFromFilename(migrationFiles[i])
		versionJ := extractVersionFromFilename(migrationFiles[j])

		// Извлекаем числовую часть версии для сортировки
		numI := extractNumberFromVersion(versionI)
		numJ := extractNumberFromVersion(versionJ)

		return numI < numJ
	})

	return migrationFiles, nil
}

func extractVersionFromFilename(filename string) string {
	// Извлекаем версию из имени файла (например, "001_initial_schema" из "001_initial_schema.up.sql")
	parts := strings.Split(filename, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return filename
}

func extractNumberFromVersion(version string) int {
	// Извлекаем число из начала версии (например, 1 из "001_initial_schema")
	parts := strings.Split(version, "_")
	if len(parts) > 0 {
		if num, err := strconv.Atoi(parts[0]); err == nil {
			return num
		}
	}
	return 0
}
