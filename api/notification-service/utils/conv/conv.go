package conv

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func ParseStringToDateTime(dateStr string, timeStr string) (*time.Time, *time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, nil, err
	}

	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return nil, nil, err
	}

	return &parsedDate, &parsedTime, nil
}

func ParseStringToBool(s string) (bool, error) {
	parsedBool, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}

	return parsedBool, err
}

func ParseInt64QueryParam(c echo.Context, param string, defaultVal int64) (int64, error) {
	strVal := c.QueryParam(param)

	if strVal == "" {
		return defaultVal, nil
	}

	val, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return 0, err
	}

	if val <= 0 {
		return defaultVal, nil
	}

	return val, nil
}
