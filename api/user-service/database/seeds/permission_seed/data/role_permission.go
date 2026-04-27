package data

import "user-service/internal/core/domain/model"

var RolePermissions = []model.Permission{
	{
		Resource: "roles",
		Action:   "read",
		Scope:    "own",
	},
	{
		Resource: "roles",
		Action:   "write",
		Scope:    "own",
	},
	{
		Resource: "roles",
		Action:   "update",
		Scope:    "own",
	},
	{
		Resource: "roles",
		Action:   "delete",
		Scope:    "own",
	},
	{
		Resource: "roles",
		Action:   "read",
		Scope:    "all",
	},
	{
		Resource: "roles",
		Action:   "write",
		Scope:    "all",
	},
	{
		Resource: "roles",
		Action:   "update",
		Scope:    "all",
	},
	{
		Resource: "roles",
		Action:   "delete",
		Scope:    "all",
	},
}
