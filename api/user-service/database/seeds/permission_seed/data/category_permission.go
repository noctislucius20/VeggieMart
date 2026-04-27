package data

import "user-service/internal/core/domain/model"

var CategoryPermissions = []model.Permission{
	{
		Resource: "categories",
		Action:   "read",
		Scope:    "own",
	},
	{
		Resource: "categories",
		Action:   "write",
		Scope:    "own",
	},
	{
		Resource: "categories",
		Action:   "update",
		Scope:    "own",
	},
	{
		Resource: "categories",
		Action:   "delete",
		Scope:    "own",
	},
	{
		Resource: "categories",
		Action:   "read",
		Scope:    "all",
	},
	{
		Resource: "categories",
		Action:   "write",
		Scope:    "all",
	},
	{
		Resource: "categories",
		Action:   "update",
		Scope:    "all",
	},
	{
		Resource: "categories",
		Action:   "delete",
		Scope:    "all",
	},
}
