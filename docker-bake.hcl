# Define variables with direct values (fallbacks)
APP_NAME = "dockerwizard-api"
FRAMEWORK = "fiber"
VERSION = "1.0.0"
DB_PASS = "dockerwizard-admin"
DB_NAME = "dockerwizard"
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
  dockerfile = "Dockerfile.mysql"
  tags = ["${APP_NAME}-mysql:${VERSION}"]
  platforms = ["linux/amd64"]
  args = {
    MYSQL_ROOT_PASSWORD = DB_PASS
    MYSQL_DATABASE = DB_NAME
    MYSQL_USER = DB_USER
    MYSQL_PASSWORD = DB_PASS
  }
}