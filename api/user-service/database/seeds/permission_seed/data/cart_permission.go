package data

import "user-service/internal/core/domain/model"

var CartPermissions = []model.Permission{
	{
		Resource: "carts",
		Action:   "read",
		Scope:    "own",
	},
	{
		Resource: "carts",
		Action:   "write",
		Scope:    "own",
	},
	{
		Resource: "carts",
		Action:   "update",
		Scope:    "own",
	},
	{
		Resource: "carts",
		Action:   "delete",
		Scope:    "own",
	},
	{
		Resource: "carts",
		Action:   "read",
		Scope:    "all",
	},
	{
		Resource: "carts",
		Action:   "write",
		Scope:    "all",
	},
	{
		Resource: "carts",
		Action:   "update",
		Scope:    "all",
	},
	{
		Resource: "carts",
		Action:   "delete",
		Scope:    "all",
	},
}
