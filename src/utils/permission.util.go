package utils

import "skypipe/src/lib/dto"

var Permissions = map[string]dto.PermissionInterface{
	// User permissions
	"USER_CREATE": {Resource: "user", Action: "create"},
	"USER_READ":   {Resource: "user", Action: "read"},
	"USER_UPDATE": {Resource: "user", Action: "update"},
	"USER_DELETE": {Resource: "user", Action: "delete"},

	// Project permissions
	"PROJECT_CREATE": {Resource: "project", Action: "create"},
	"PROJECT_READ":   {Resource: "project", Action: "read"},
	"PROJECT_UPDATE": {Resource: "project", Action: "update"},
	"PROJECT_DELETE": {Resource: "project", Action: "delete"},

	// Deployment permissions
	"DEPLOYMENT_CREATE": {Resource: "deployment", Action: "create"},
	"DEPLOYMENT_READ":   {Resource: "deployment", Action: "read"},
	"DEPLOYMENT_UPDATE": {Resource: "deployment", Action: "update"},
	"DEPLOYMENT_DELETE": {Resource: "deployment", Action: "delete"},
}
