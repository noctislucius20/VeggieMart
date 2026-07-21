package data

import "user-service/internal/core/domain/model"

var NotificationPermissions = []model.Permission{
	{
		Resource: "notifications",
		Action:   "read",
		Scope:    "own",
	},
	{
		Resource: "notifications",
		Action:   "write",
		Scope:    "own",
	},
	{
		Resource: "notifications",
		Action:   "update",
		Scope:    "own",
	},
	{
		Resource: "notifications",
		Action:   "delete",
		Scope:    "own",
	},
	{
		Resource: "notifications",
		Action:   "read",
		Scope:    "all",
	},
	{
		Resource: "notifications",
		Action:   "write",
		Scope:    "all",
	},
	{
		Resource: "notifications",
		Action:   "update",
		Scope:    "all",
	},
	{
		Resource: "notifications",
		Action:   "delete",
		Scope:    "all",
	},
}
