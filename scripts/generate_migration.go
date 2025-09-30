package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/generate_migration.go <migration_name>")
		fmt.Println("Example: go run scripts/generate_migration.go add_user_table")
		os.Exit(1)
	}

	migrationName := toSnakeCase(os.Args[1])

	// Validate migration name (only alphanumeric and underscores)
	if !isValidMigrationName(migrationName) {
		fmt.Println("Error: Migration name can only contain letters, numbers, and underscores")
		os.Exit(1)
	}

	// Get next migration number
	nextNumber, err := getNextMigrationNumber()
	if err != nil {
		fmt.Printf("Error getting next migration number: %v\n", err)
		os.Exit(1)
	}

	// Create migration folder name
	migrationFolderName := fmt.Sprintf("%03d_%s", nextNumber, migrationName)
	migrationDir := filepath.Join("migrations", migrationFolderName)

	// Create migration directory
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		fmt.Printf("Error creating migration directory: %v\n", err)
		os.Exit(1)
	}

	// Create up.sql file
	upFilePath := filepath.Join(migrationDir, "up.sql")
	upContent := fmt.Sprintf(`-- migrations/%s/up.sql
-- Created at: %s

-- Add your migration SQL here

`, migrationFolderName, time.Now().Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(upFilePath, []byte(upContent), 0644); err != nil {
		fmt.Printf("Error creating up.sql file: %v\n", err)
		os.Exit(1)
	}

	// Create down.sql file
	downFilePath := filepath.Join(migrationDir, "down.sql")
	downContent := fmt.Sprintf(`-- migrations/%s/down.sql
-- Created at: %s

-- Add your rollback SQL here

`, migrationFolderName, time.Now().Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(downFilePath, []byte(downContent), 0644); err != nil {
		fmt.Printf("Error creating down.sql file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Migration created successfully!\n")
	fmt.Printf("📁 Directory: %s\n", migrationDir)
	fmt.Printf("📝 Files created:\n")
	fmt.Printf("   - %s\n", upFilePath)
	fmt.Printf("   - %s\n", downFilePath)
	fmt.Printf("\n💡 Next steps:\n")
	fmt.Printf("   1. Edit %s with your migration SQL\n", upFilePath)
	fmt.Printf("   2. Edit %s with your rollback SQL\n", downFilePath)
	fmt.Printf("   3. Run: make migrate-up\n")
}

func toSnakeCase(str string) string {
	// Replace spaces and hyphens with underscores
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ReplaceAll(str, "-", "_")

	// Convert camelCase and PascalCase to snake_case
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	str = re.ReplaceAllString(str, `${1}_${2}`)

	// Convert to lowercase
	str = strings.ToLower(str)

	// Remove multiple consecutive underscores
	re = regexp.MustCompile(`_+`)
	str = re.ReplaceAllString(str, "_")

	// Remove leading and trailing underscores
	str = strings.Trim(str, "_")

	return str
}

func isValidMigrationName(name string) bool {
	// Allow letters, numbers, and underscores only
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, name)
	return matched
}

func getNextMigrationNumber() (int, error) {
	migrationsDir := "migrations"

	// Create migrations directory if it doesn't exist
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(migrationsDir, 0755); err != nil {
			return 0, fmt.Errorf("failed to create migrations directory: %w", err)
		}
		return 1, nil
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationNumbers []int

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		// Extract number from directory name (e.g., "001_initial_schema" -> 1)
		parts := strings.Split(file.Name(), "_")
		if len(parts) > 0 {
			if num, err := strconv.Atoi(parts[0]); err == nil {
				migrationNumbers = append(migrationNumbers, num)
			}
		}
	}

	if len(migrationNumbers) == 0 {
		return 1, nil
	}

	sort.Ints(migrationNumbers)
	return migrationNumbers[len(migrationNumbers)-1] + 1, nil
}