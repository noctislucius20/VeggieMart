package utils

import "errors"

var (
	ErrSessionExpired       = errors.New("401: SESSION_EXPIRED")
	ErrTokenInvalid         = errors.New("401: TOKEN_INVALID")
	ErrTokenExpired         = errors.New("401: TOKEN_EXPIRED")
	ErrLoginInvalid         = errors.New("401: LOGIN_INVALID")
	ErrAccessForbidden      = errors.New("403: ACCESS_FORBIDDEN")
	ErrDataNotFound         = errors.New("404: DATA_NOT_FOUND")
	ErrEmailAlreadyExists   = errors.New("409: EMAIL_ALREADY_EXISTS")
	ErrEmailNotVerified     = errors.New("409: EMAIL_NOT_VERIFIED")
	ErrDataAlreadyExists    = errors.New("409: DATA_ALREADY_EXISTS")
	ErrDataStillInUsed      = errors.New("409: DATA_STILL_IN_USED")
	ErrRelationDataNotFound = errors.New("422: RELATION_DATA_NOT_FOUND")
	ErrIDInvalid            = errors.New("422: ID_INVALID")
	ErrTooManyRequests      = errors.New("429: TOO_MANY_REQUESTS")

	ErrInternalServerError  = errors.New("500: INTERNAL_SERVER_ERROR")
	ErrServiceUnavailable   = errors.New("503: SERVICE_UNAVAILABLE")
	ErrTimeoutLimitExceeded = errors.New("504: TIMEOUT_LIMIT_EXCEEDED")
)
