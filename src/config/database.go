package config

import (
	apiKeys "dockerwizard-api/src/modules/apiKeys/models"
	billing "dockerwizard-api/src/modules/billing/model"
	ci "dockerwizard-api/src/modules/ci/models"
	deployments "dockerwizard-api/src/modules/deployments/models"
	github "dockerwizard-api/src/modules/github/models"
	logs "dockerwizard-api/src/modules/logs/models"
	notifications "dockerwizard-api/src/modules/notifications/model"
	plans "dockerwizard-api/src/modules/plans/models"
	projects "dockerwizard-api/src/modules/projects/models"
	secrets "dockerwizard-api/src/modules/secrets/models"
	teams "dockerwizard-api/src/modules/teams/models"
	templates "dockerwizard-api/src/modules/templates/models"
	users "dockerwizard-api/src/modules/users/models"
	webhooks "dockerwizard-api/src/modules/webhooks/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

// ConnectDatabase initializes and migrates the database.
func ConnectDatabase() *gorm.DB {
	// Retrieve connection info from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Format MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Connect to MySQL
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable SQL query logs
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = database
	fmt.Println("Connected to database!")

	// Perform database migrations
	if err := runMigrations(DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	return DB
}

func CheckConnection() bool {
	if DB == nil {
		return false
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Failed to get generic database object: %v", err)
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		log.Printf("Database ping failed: %v", err)
		return false
	}

	var result int
	if err := DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		log.Printf("Test query failed: %v", err)
		return false
	}
	return result == 1
}

// runMigrations runs all model migrations.
func runMigrations(db *gorm.DB) error {
	migrations := []func(*gorm.DB) error{
		plans.MigratePlan,
		users.MigrateUserCore,
		billing.MigrateBillingSubscriptions,
		teams.MigrateTeams,
		teams.MigrateTeamMembers,
		projects.MigrateProjects,
		templates.MigrateProjectTemplates,
		projects.MigrateProjectConfigs,
		projects.MigrateProjectFiles,
		ci.MigrateCIPipelines,
		ci.MigratePipelineSteps,
		deployments.MigrateDeployments,
		deployments.MigrateDeploymentTargets,
		github.MigrateGitHubIntegrations,
		logs.MigrateActivityLogs,
		apiKeys.MigrateAPIKeys,
		notifications.MigrateNotifications,
		webhooks.MigrateWebhooks,
		secrets.MigrateSecrets,
		templates.MigrateUsageMetrics,
	}

	// Iterate through all migrations
	for _, migrate := range migrations {
		if err := migrate(db); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	fmt.Println("All migrations completed successfully!")
	return nil
}
