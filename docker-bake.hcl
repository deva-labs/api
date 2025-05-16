# Define variables with direct values (fallbacks)
APP_NAME = "skypipe-api"
FRAMEWORK = "fiber"
VERSION = "1.0.0"
DB_PASS = "skypipe-admin"
DB_NAME = "skypipe"
DB_USER = "admin"

group "default" {
  targets = ["app", "db"]
}

target "app" {
  context = "."
  dockerfile = "Dockerfile"
  tags = ["${APP_NAME}-${FRAMEWORK}:${VERSION}"]
  platforms = ["linux/amd64"]
  args = {
    APP_NAME = APP_NAME
    VERSION = VERSION
  }
}

target "db" {
  context = "."
  dockerfile = "Dockerfile.postgres"
  tags = ["${APP_NAME}-postgres:${VERSION}"]
  platforms = ["linux/amd64"]
  args = {
    POSTGRES_DB: DB_NAME
    POSTGRES_USER: DB_USER
    POSTGRES_PASSWORD: DB_PASS
  }
}