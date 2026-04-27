package data

import "user-service/internal/core/domain/model"

var OrderPermissions = []model.Permission{
	{
		Resource: "orders",
		Action:   "read",
		Scope:    "own",
	},
	{
		Resource: "orders",
		Action:   "write",
		Scope:    "own",
	},
	{
		Resource: "orders",
		Action:   "update",
		Scope:    "own",
	},
	{
		Resource: "orders",
		Action:   "delete",
		Scope:    "own",
	},
	{
		Resource: "orders",
		Action:   "read",
		Scope:    "all",
	},
	{
		Resource: "orders",
		Action:   "write",
		Scope:    "all",
	},
	{
		Resource: "orders",
		Action:   "update",
		Scope:    "all",
	},
	{
		Resource: "orders",
		Action:   "delete",
		Scope:    "all",
	},
}
