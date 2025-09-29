package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/zandomed/sync-playlist-api/internal/config"
)

func main() {
	// Obtener configuración singleton
	cfg := config.Get()

	// Conectar a la base de datos
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Crear tabla de migraciones si no existe
	if err := createMigrationsTable(db); err != nil {
		log.Fatal("Error creating migrations table:", err)
	}

	// Obtener comando
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		if err := runMigrations(db); err != nil {
			log.Fatal("Error running migrations:", err)
		}
		fmt.Println("✅ All migrations applied successfully")
	case "down":
		if err := rollbackMigration(db); err != nil {
			log.Fatal("Error rolling back migration:", err)
		}
		fmt.Println("✅ Migration rolled back successfully")
	case "status":
		if err := showMigrationStatus(db); err != nil {
			log.Fatal("Error showing migration status:", err)
		}
	default:
		fmt.Println("Usage: go run scripts/migrate.go [up|down|status]")
		os.Exit(1)
	}
}

func createMigrationsTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)`
	_, err := db.Exec(query)
	return err
}

func runMigrations(db *sqlx.DB) error {
	// Leer archivos de migración
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return err
	}

	// Obtener migraciones ya aplicadas
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	// Aplicar migraciones pendientes
	for _, file := range migrationFiles {
		version := getVersionFromFilename(file)
		if _, exists := appliedMigrations[version]; exists {
			fmt.Printf("⏭️  Skipping %s (already applied)\n", file)
			continue
		}

		fmt.Printf("🔄 Applying %s...\n", file)
		if err := applyMigration(db, file, version); err != nil {
			return fmt.Errorf("failed to apply %s: %w", file, err)
		}
		fmt.Printf("✅ Applied %s\n", file)
	}

	return nil
}

func rollbackMigration(db *sqlx.DB) error {
	// Obtener la última migración aplicada
	var lastVersion string
	query := `SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1`
	err := db.Get(&lastVersion, query)
	if err != nil {
		return fmt.Errorf("no migrations to rollback: %w", err)
	}

	fmt.Printf("🔄 Rolling back migration %s...\n", lastVersion)

	// Buscar el archivo de rollback .down.sql
	downFile, err := findDownMigrationFile(lastVersion)
	if err != nil {
		return fmt.Errorf("failed to find down migration for %s: %w", lastVersion, err)
	}

	if downFile == "" {
		fmt.Printf("⚠️  No .down.sql file found for migration %s\n", lastVersion)
		fmt.Printf("⚠️  Only removing migration record. Manual cleanup may be required.\n")

		// Solo remover el registro
		deleteQuery := `DELETE FROM schema_migrations WHERE version = $1`
		_, err = db.Exec(deleteQuery, lastVersion)
		return err
	}

	// Leer contenido del archivo .down.sql
	downSQL, err := os.ReadFile(filepath.Join("migrations", downFile))
	if err != nil {
		return fmt.Errorf("failed to read down migration file: %w", err)
	}

	// Iniciar transacción
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	// Ejecutar rollback SQL
	if _, err := tx.Exec(string(downSQL)); err != nil {
		return fmt.Errorf("failed to execute rollback: %w", err)
	}

	// Remover registro de migración
	deleteQuery := `DELETE FROM schema_migrations WHERE version = $1`
	if _, err := tx.Exec(deleteQuery, lastVersion); err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit transacción
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback transaction: %w", err)
	}

	fmt.Printf("✅ Successfully rolled back migration %s\n", lastVersion)
	return nil
}

func showMigrationStatus(db *sqlx.DB) error {
	// Obtener todas las migraciones
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return err
	}

	// Obtener migraciones aplicadas
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	for _, file := range migrationFiles {
		version := getVersionFromFilename(file)
		status := "❌ Pending"
		if appliedAt, exists := appliedMigrations[version]; exists {
			status = fmt.Sprintf("✅ Applied (%s)", appliedAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("%-40s %s\n", file, status)
	}

	return nil
}

func getMigrationFiles() ([]string, error) {
	migrationsDir := "migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// Solo incluir archivos .up.sql para evitar duplicados
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, file.Name())
		} else if strings.HasSuffix(file.Name(), ".sql") && !strings.Contains(file.Name(), ".up.") && !strings.Contains(file.Name(), ".down.") {
			// También incluir archivos .sql regulares que no tengan .up. o .down.
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

func getAppliedMigrations(db *sqlx.DB) (map[string]time.Time, error) {
	var migrations []struct {
		Version   string    `db:"version"`
		AppliedAt time.Time `db:"applied_at"`
	}

	query := `SELECT version, applied_at FROM schema_migrations`
	if err := db.Select(&migrations, query); err != nil {
		return nil, err
	}

	appliedMigrations := make(map[string]time.Time)
	for _, migration := range migrations {
		appliedMigrations[migration.Version] = migration.AppliedAt
	}

	return appliedMigrations, nil
}

func applyMigration(db *sqlx.DB, filename, version string) error {
	// Determinar el archivo correcto a usar
	var migrationFile string
	if strings.HasSuffix(filename, ".up.sql") {
		migrationFile = filename
	} else {
		// Para archivos .sql regulares, verificar si existe un .up.sql correspondiente
		upFile := strings.Replace(filename, ".sql", ".up.sql", 1)
		upPath := filepath.Join("migrations", upFile)
		if _, err := os.Stat(upPath); err == nil {
			migrationFile = upFile
		} else {
			migrationFile = filename
		}
	}

	// Leer archivo de migración
	filePath := filepath.Join("migrations", migrationFile)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Iniciar transacción
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	// Ejecutar migración
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Registrar migración como aplicada
	insertQuery := `INSERT INTO schema_migrations (version) VALUES ($1)`
	if _, err := tx.Exec(insertQuery, version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transacción
	return tx.Commit()
}

func getVersionFromFilename(filename string) string {
	// Extraer versión del nombre del archivo
	// Soporta: "001_initial_schema.sql", "001_initial_schema.up.sql", "001_initial_schema.down.sql"
	parts := strings.Split(filename, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return filename
}

// findDownMigrationFile busca el archivo .down.sql correspondiente a una versión
func findDownMigrationFile(version string) (string, error) {
	migrationsDir := "migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if strings.HasSuffix(filename, ".down.sql") && getVersionFromFilename(filename) == version {
			return filename, nil
		}
	}

	return "", nil // No se encontró archivo .down.sql
}
