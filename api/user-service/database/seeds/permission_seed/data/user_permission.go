package data

import "user-service/internal/core/domain/model"

var UserPermissions = []model.Permission{
	{
		Resource: "users",
		Action:   "read",
		Scope:    "own",
	},
	{
		Resource: "users",
		Action:   "write",
		Scope:    "own",
	},
	{
		Resource: "users",
		Action:   "update",
		Scope:    "own",
	},
	{
		Resource: "users",
		Action:   "delete",
		Scope:    "own",
	},
	{
		Resource: "users",
		Action:   "read",
		Scope:    "all",
	},
	{
		Resource: "users",
		Action:   "write",
		Scope:    "all",
	},
	{
		Resource: "users",
		Action:   "update",
		Scope:    "all",
	},
	{
		Resource: "users",
		Action:   "delete",
		Scope:    "all",
	},
}
