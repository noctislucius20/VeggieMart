package data

import "user-service/internal/core/domain/model"

var PaymentPermissions = []model.Permission{
	{
		Resource: "payments",
		Action:   "read",
		Scope:    "own",
	},
	{
		Resource: "payments",
		Action:   "write",
		Scope:    "own",
	},
	{
		Resource: "payments",
		Action:   "update",
		Scope:    "own",
	},
	{
		Resource: "payments",
		Action:   "delete",
		Scope:    "own",
	},
	{
		Resource: "payments",
		Action:   "read",
		Scope:    "all",
	},
	{
		Resource: "payments",
		Action:   "write",
		Scope:    "all",
	},
	{
		Resource: "payments",
		Action:   "update",
		Scope:    "all",
	},
	{
		Resource: "payments",
		Action:   "delete",
		Scope:    "all",
	},
}
