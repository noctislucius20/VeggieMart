package data

import "user-service/internal/core/domain/model"

var ProductPermissions = []model.Permission{
	{
		Resource: "products",
		Action:   "read",
		Scope:    "own",
	},
	{
		Resource: "products",
		Action:   "write",
		Scope:    "own",
	},
	{
		Resource: "products",
		Action:   "update",
		Scope:    "own",
	},
	{
		Resource: "products",
		Action:   "delete",
		Scope:    "own",
	},
	{
		Resource: "products",
		Action:   "read",
		Scope:    "all",
	},
	{
		Resource: "products",
		Action:   "write",
		Scope:    "all",
	},
	{
		Resource: "products",
		Action:   "update",
		Scope:    "all",
	},
	{
		Resource: "products",
		Action:   "delete",
		Scope:    "all",
	},
}
