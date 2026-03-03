package conv

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}

func ParseInt64QueryParam(
	c echo.Context,
	param string,
	defaultVal int64,
) (int64, error) {

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
