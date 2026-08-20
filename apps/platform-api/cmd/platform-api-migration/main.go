package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/blocknextai/go-packages/cast"
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/caarlos0/env/v11"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var modules = []string{
	"common",
	"account",
	"organizations",
	"executions",
	"workflows",
	"triggers",
	"credentials",
	"apikeys",
	"notifications",
	"eventbus",
}

func validateAndGetTargetModules(module string) []string {
	if module == "" {
		return modules
	}

	if slices.Contains(modules, module) {
		return []string{module}
	}

	slog.Error("Invalid module specified",
		"component", "migration",
		"module", module,
		"available_modules", strings.Join(modules, ", "),
	)
	os.Exit(1)
	return nil
}

func main() {
	var (
		direction = flag.String("direction", "up", "Migration direction: up or down")
		module    = flag.String("module", "", "Specific module to migrate (optional)")
		dryRun    = flag.Bool("dry-run", false, "Run in dry-run mode to check migrations")
	)
	flag.Parse()

	var dbOpts config.DatabaseOptions
	if err := env.ParseWithOptions(&dbOpts, env.Options{Prefix: "DATABASE_"}); err != nil {
		slog.Error("failed to parse database environment variables", "component", "migration", "error", err)
		os.Exit(1)
	}

	if strings.TrimSpace(dbOpts.Host) == "" || strings.TrimSpace(dbOpts.Name) == "" {
		slog.Error("DATABASE_HOST and DATABASE_NAME environment variables are required", "component", "migration")
		os.Exit(1)
	}

	databaseName := dbOpts.Name

	db, err := database.NewDB(
		dbOpts.DSN(),
		1,
		1,
		dbOpts.ConnMaxLifetime,
		dbOpts.ConnMaxIdleTime,
	)
	if err != nil {
		slog.Error("Failed to connect to database",
			"component", "migration",
			"error", err,
		)
		os.Exit(1)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("Failed to close SQL DB instance",
				"component", "migration",
				"error", err,
			)
		}
	}()

	targetModules := validateAndGetTargetModules(*module)

	switch *direction {
	case "up":
		err = migrateUp(db, databaseName, targetModules, *dryRun)
	case "down":
		err = migrateDown(db, databaseName, targetModules, *dryRun)
	default:
		slog.Error("Invalid direction. Use 'up' or 'down'", "component", "migration")
		os.Exit(1)
	}

	if err != nil {
		slog.Error("Migration failed",
			"component", "migration",
			"error", err,
		)
		os.Exit(1)
	}

	if !*dryRun {
		slog.Info("Migration completed successfully", "component", "migration")
	}
}

func isDockerEnvironment() bool {
	_, err := os.Stat("/app/migrations")
	return err == nil
}

func getModuleMigrationPath(module string) (string, string) {
	if isDockerEnvironment() {
		path := "/app/migrations/" + module
		return "file://" + path, path
	}
	path := "internal/" + module + "/infrastructure/database/migrations"
	return "file://" + path, path
}

func hasMigrationFiles(module string) bool {
	_, migrationPath := getModuleMigrationPath(module)

	info, err := os.Stat(migrationPath)
	if err != nil || !info.IsDir() {
		return false
	}

	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return false
	}

	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".up.sql") || strings.HasSuffix(file.Name(), ".down.sql")) {
			return true
		}
	}

	return false
}

func closeMigrator(m *migrate.Migrate, module string) {
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		slog.Warn("Failed to close migration source",
			"component", "migration",
			"module", module,
			"error", srcErr,
		)
	}
	if dbErr != nil {
		slog.Warn("Failed to close migration database driver",
			"component", "migration",
			"module", module,
			"error", dbErr,
		)
	}
}

func createMigrator(db *sql.DB, databaseName string, modulePath string, moduleName string) (*migrate.Migrate, error) {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithConnection(context.Background(), conn, &postgres.Config{
		DatabaseName:    databaseName,
		MigrationsTable: moduleName + "_migrations",
	})
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(modulePath, databaseName, driver)
	if err != nil {
		_ = driver.Close()
		return nil, err
	}

	return m, nil
}

func migrateUp(db *sql.DB, databaseName string, targetModules []string, dryRun bool) error {
	if dryRun {
		slog.Info("DRY RUN - Checking UP migrations", "component", "migration")
		return checkMigrationStatus(db, databaseName, targetModules, "up")
	}

	slog.Info("Running UP migrations", "component", "migration")
	for _, module := range targetModules {
		slog.Info("Processing module", "component", "migration", "module", module)

		if !hasMigrationFiles(module) {
			slog.Info("Skipping - no migration files found",
				"component", "migration",
				"module", module,
			)
			continue
		}

		if err := applyMigrationUp(db, databaseName, module); err != nil {
			return err
		}
	}

	return nil
}

func migrateDown(db *sql.DB, databaseName string, targetModules []string, dryRun bool) error {
	if dryRun {
		slog.Info("DRY RUN - Checking DOWN migrations", "component", "migration")
		return checkMigrationStatus(db, databaseName, targetModules, "down")
	}

	slog.Info("Running DOWN migrations", "component", "migration")
	for _, module := range slices.Backward(targetModules) {
		slog.Info("Processing module", "component", "migration", "module", module)

		if !hasMigrationFiles(module) {
			slog.Info("Skipping - no migration files found",
				"component", "migration",
				"module", module,
			)
			continue
		}

		if err := applyMigrationDown(db, databaseName, module); err != nil {
			return err
		}
	}

	return nil
}

func checkMigrationStatus(db *sql.DB, databaseName string, targetModules []string, direction string) error {
	modulesToCheck := targetModules
	if direction == "down" {
		reversed := make([]string, 0, len(targetModules))
		for _, targetModule := range slices.Backward(targetModules) {
			reversed = append(reversed, targetModule)
		}
		modulesToCheck = reversed
	}

	for _, module := range modulesToCheck {
		slog.Info("Checking module status",
			"component", "migration",
			"module", module,
		)
		if !hasMigrationFiles(module) {
			slog.Info("Skipping - no migration files found",
				"component", "migration",
				"module", module,
			)
			continue
		}

		modulePath, _ := getModuleMigrationPath(module)
		m, err := createMigrator(db, databaseName, modulePath, module)
		if err != nil {
			slog.Error("Error creating migrator",
				"component", "migration",
				"module", module,
				"error", err,
			)
			continue
		}

		version, dirty, err := m.Version()
		closeMigrator(m, module)
		statusMsg := getMigrationStatusMessage(version, dirty, err, direction)
		slog.Info(statusMsg,
			"component", "migration",
			"module", module,
		)

		if direction == "up" && !dirty && (err == nil || err == migrate.ErrNilVersion) {
			pending := checkPendingMigrations(version, err == nil, module)
			if pending > 0 {
				slog.Warn("Pending migrations found",
					"component", "migration",
					"module", module,
					"pending_count", pending,
				)
			}
		}
	}

	return nil
}

func getMigrationStatusMessage(version uint, dirty bool, err error, direction string) string {
	if err == migrate.ErrNilVersion {
		return "Clean (no migrations applied)"
	}

	if err != nil {
		return "Error: " + err.Error()
	}

	if dirty {
		return "Dirty state at version " + cast.ToString(version)
	}

	if direction == "down" {
		return "Version " + cast.ToString(version) + " (would rollback 1 migration step)"
	}

	return "Version " + cast.ToString(version)
}

func checkPendingMigrations(currentVersion uint, hasAppliedVersion bool, module string) int {
	_, migrationPath := getModuleMigrationPath(module)

	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return 0
	}

	pending := 0
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}

		underscoreIdx := strings.Index(file.Name(), "_")
		if underscoreIdx <= 0 {
			continue
		}

		fileVersion, parseErr := strconv.ParseUint(file.Name()[:underscoreIdx], 10, 64)
		if parseErr != nil {
			continue
		}

		if !hasAppliedVersion || uint(fileVersion) > currentVersion {
			pending++
		}
	}

	return pending
}

func applyMigrationUp(db *sql.DB, databaseName string, module string) error {
	modulePath, _ := getModuleMigrationPath(module)
	m, err := createMigrator(db, databaseName, modulePath, module)
	if err != nil {
		slog.Error("Error creating migrator",
			"component", "migration",
			"module", module,
			"error", err,
		)
		return err
	}
	defer closeMigrator(m, module)

	upErr := m.Up()
	if upErr != nil && upErr != migrate.ErrNoChange {
		slog.Error("Migration failed",
			"component", "migration",
			"module", module,
			"error", upErr,
		)
		return upErr
	}

	if upErr == migrate.ErrNoChange {
		slog.Info("Already up to date",
			"component", "migration",
			"module", module,
		)
	} else {
		slog.Info("Migrations applied successfully",
			"component", "migration",
			"module", module,
		)
	}

	return nil
}

func applyMigrationDown(db *sql.DB, databaseName string, module string) error {
	modulePath, _ := getModuleMigrationPath(module)
	m, err := createMigrator(db, databaseName, modulePath, module)
	if err != nil {
		slog.Error("Error creating migrator",
			"component", "migration",
			"module", module,
			"error", err,
		)
		return err
	}
	defer closeMigrator(m, module)

	slog.Info("Rolling back 1 migration step",
		"component", "migration",
		"module", module,
	)
	downErr := m.Steps(-1)
	if downErr != nil && downErr != migrate.ErrNoChange {
		slog.Error("Rollback failed",
			"component", "migration",
			"module", module,
			"error", downErr,
		)
		return downErr
	}

	if downErr == migrate.ErrNoChange {
		slog.Info("No migrations to rollback",
			"component", "migration",
			"module", module,
		)
	} else {
		slog.Info("Rollback successful",
			"component", "migration",
			"module", module,
		)
	}

	return nil
}
