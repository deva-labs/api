package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	apiKeys "skypipe/src/modules/apiKeys/models"
	billing "skypipe/src/modules/billing/model"
	ci "skypipe/src/modules/ci/models"
	deployments "skypipe/src/modules/deployments/models"
	github "skypipe/src/modules/github/models"
	key_token "skypipe/src/modules/key_token/models"
	logs "skypipe/src/modules/logs/models"
	notifications "skypipe/src/modules/notifications/models"
	plans "skypipe/src/modules/plans/models"
	projects "skypipe/src/modules/projects/models"
	roles "skypipe/src/modules/roles/models"
	secrets "skypipe/src/modules/secrets/models"
	teams "skypipe/src/modules/teams/models"
	templates "skypipe/src/modules/templates/models"
	users "skypipe/src/modules/users/models"
	verifications "skypipe/src/modules/verifications/models"
	webhooks "skypipe/src/modules/webhooks/models"
	"skypipe/src/utils"
)

var DB *gorm.DB

// ConnectDatabase initializes and migrates the database.
func ConnectDatabase() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dbHost, dbPort, dbUser, dbPass, dbName)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = database
	fmt.Println("Connected to PostgreSQL database!")

	// Enable uuid on db
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
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

// runMigrations runs all models migrations.
func runMigrations(db *gorm.DB) error {
	migrations := []func(*gorm.DB) error{
		plans.MigratePlan,
		users.MigrateUserCore,
		users.MigrateProfile,
		roles.MigrateRole,
		roles.MigratePermission,
		roles.MigrateRolePermission,
		users.MigrateUserRole,
		key_token.MigrateAccessTokens,
		key_token.MigrateRefreshTokens,
		billing.MigrateBillingSubscriptions,
		teams.MigrateTeams,
		teams.MigrateTeamMembers,
		projects.MigrateProjects,
		projects.MigrateProjectDeployment,
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
		verifications.MigrateVerificationCode,
	}

	// Iterate through all migrations
	for _, migrate := range migrations {
		if err := migrate(db); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	// Seed for new system
	err := SeedForNewSystem(db)
	if err != nil {
		return err
	}
	fmt.Println("All migrations completed successfully!")
	return nil
}

func SeedForNewSystem(db *gorm.DB) error {
	log.Println("🔄 Starting system seeding...")

	if err := seedPlans(db); err != nil {
		return fmt.Errorf("seeding plans failed: %w", err)
	}

	if err := seedRoles(db); err != nil {
		return fmt.Errorf("seeding roles failed: %w", err)
	}

	createdPermissions, err := seedPermissions(db)
	if err != nil {
		return fmt.Errorf("seeding permissions failed: %w", err)
	}

	if err := assignPermissionsToAdmin(db, createdPermissions); err != nil {
		return fmt.Errorf("assigning permissions to admin failed: %w", err)
	}

	if err := seedRootUser(db); err != nil {
		return fmt.Errorf("seeding root user failed: %w", err)
	}

	log.Println("✅ System seeding completed successfully")
	return nil
}

func seedPlans(db *gorm.DB) error {
	log.Println("🔹 Seeding plans...")
	defaultPlans := []plans.Plan{
		{Name: "Free", Price: 0, MaxProjects: 1, MaxDeploymentsPerDay: 5},
		{Name: "Pro", Price: 19.99, MaxProjects: 10, MaxDeploymentsPerDay: 50},
		{Name: "Enterprise", Price: 99.99, MaxProjects: 100, MaxDeploymentsPerDay: 1000},
	}

	for _, plan := range defaultPlans {
		if err := db.FirstOrCreate(&plan, plans.Plan{Name: plan.Name}).Error; err != nil {
			return err
		}
	}
	return nil
}
func seedRoles(db *gorm.DB) error {
	log.Println("🔹 Seeding roles...")
	rolesToSeed := []roles.Role{
		{Name: "admin", Description: "Administrator with full access"},
		{Name: "user", Description: "Regular user"},
		{Name: "devops", Description: "DevOps engineer"},
	}

	for _, role := range rolesToSeed {
		if err := db.FirstOrCreate(&role, roles.Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}
	return nil
}
func seedPermissions(db *gorm.DB) ([]roles.Permission, error) {
	log.Println("🔹 Seeding permissions...")
	permissions := []roles.Permission{
		{Resource: "user", Action: "create"},
		{Resource: "user", Action: "read"},
		{Resource: "user", Action: "update"},
		{Resource: "user", Action: "delete"},
		{Resource: "project", Action: "create"},
		{Resource: "project", Action: "read"},
		{Resource: "project", Action: "update"},
		{Resource: "project", Action: "delete"},
		{Resource: "deployment", Action: "create"},
		{Resource: "deployment", Action: "read"},
		{Resource: "deployment", Action: "update"},
		{Resource: "deployment", Action: "delete"},
	}

	var created []roles.Permission
	for _, perm := range permissions {
		var p roles.Permission
		if err := db.FirstOrCreate(&p, roles.Permission{Resource: perm.Resource, Action: perm.Action}).Error; err != nil {
			return nil, err
		}
		created = append(created, p)
	}
	return created, nil
}

func assignPermissionsToAdmin(db *gorm.DB, permissions []roles.Permission) error {
	log.Println("🔹 Assigning permissions to admin role...")
	var admin roles.Role
	if err := db.Where("name = ?", "admin").First(&admin).Error; err != nil {
		return err
	}

	for _, perm := range permissions {
		rp := roles.RolePermission{
			RoleID:       admin.ID,
			PermissionID: perm.ID,
		}
		if err := db.FirstOrCreate(&rp, rp).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedRootUser(db *gorm.DB) error {
	log.Println("🔹 Creating root user...")
	hashedPassword, err := utils.HashPassword("skypipe@root-admin")
	if err != nil {
		return err
	}

	var enterprisePlan plans.Plan
	if err := db.Where("name = ?", "Enterprise").First(&enterprisePlan).Error; err != nil {
		return err
	}

	rootUser := users.User{
		Name:     "root",
		Email:    "ties.node@outlook.com",
		PlanID:   enterprisePlan.ID,
		Password: hashedPassword,
	}
	if err := db.FirstOrCreate(&rootUser, users.User{Email: rootUser.Email}).Error; err != nil {
		return err
	}

	// Create root profile
	rootProfile := users.Profile{
		UserID:    rootUser.ID,
		FullName:  "SkyPipe Root Admin",
		Bio:       "System administrator account",
		Gender:    "Other",
		Country:   "Unknown",
		City:      "Unknown",
		UpdatedBy: rootUser.ID,
	}
	if err := db.FirstOrCreate(&rootProfile, users.Profile{UserID: rootUser.ID}).Error; err != nil {
		return fmt.Errorf("failed to create profile for root user: %w", err)
	}

	var admin roles.Role
	if err := db.Where("name = ?", "admin").First(&admin).Error; err != nil {
		return err
	}

	userRole := users.UserRole{
		UserID:    rootUser.ID,
		RoleID:    admin.ID,
		UpdatedBy: rootUser.ID,
	}
	if err := db.FirstOrCreate(&userRole, userRole).Error; err != nil {
		return err
	}

	return nil
}
