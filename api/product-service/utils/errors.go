package utils

import "errors"

var (
	ErrIDRequired           = errors.New("400: ID_REQUIRED")
	ErrProductIDRequired    = errors.New("400: PRODUCT_ID_REQUIRED")
	ErrSlugRequired         = errors.New("400: SLUG_REQUIRED")
	ErrSessionExpired       = errors.New("401: SESSION_EXPIRED")
	ErrTokenInvalid         = errors.New("401: TOKEN_INVALID")
	ErrTokenExpired         = errors.New("401: TOKEN_EXPIRED")
	ErrLoginInvalid         = errors.New("401: LOGIN_INVALID")
	ErrAccessForbidden      = errors.New("403: ACCESS_FORBIDDEN")
	ErrGatewayRequired      = errors.New("403: GATEWAY_REQUIRED")
	ErrGatewaySecretInvalid = errors.New("403: GATEWAY_SECRET_INVALID")
	ErrServiceNotAllowed    = errors.New("403: SERVICE_NOT_ALLOWED")
	ErrServiceSecretInvalid = errors.New("403: SERVICE_SECRET_INVALID")
	ErrDataNotFound         = errors.New("404: DATA_NOT_FOUND")
	ErrEmailAlreadyExists   = errors.New("409: EMAIL_ALREADY_EXISTS")
	ErrEmailNotVerified     = errors.New("409: EMAIL_NOT_VERIFIED")
	ErrDataAlreadyExists    = errors.New("409: DATA_ALREADY_EXISTS")
	ErrStockUnavailable     = errors.New("409: STOCK_UNAVAILABLE")
	ErrDataStillInUsed      = errors.New("409: DATA_STILL_IN_USED")
	ErrRelationDataNotFound = errors.New("422: RELATION_DATA_NOT_FOUND")
	ErrProductIDInvalid     = errors.New("422: PRODUCT_ID_INVALID")
	ErrIDInvalid            = errors.New("422: ID_INVALID")
	ErrQuantityInvalid      = errors.New("422: QUANTITY_INVALID")
	ErrPriceRangeInvalid    = errors.New("422: PRICE_RANGE_INVALID")

	ErrInternalServerError  = errors.New("500: INTERNAL_SERVER_ERROR")
	ErrServiceUnavailable   = errors.New("503: SERVICE_UNAVAILABLE")
	ErrTimeoutLimitExceeded = errors.New("504: TIMEOUT_LIMIT_EXCEEDED")
)
