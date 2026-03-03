package conv

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
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

func GenerateOrderCode() string {
	return fmt.Sprintf("ORD-%s-%s", time.Now().Format("20060102150405"), random.New().String(5, random.Alphanumeric))
}

func ParseLatLngToFloat64(lat string, lng string) (float64, float64, error) {
	latParsed, err := strconv.ParseFloat(lat, 64)
	lngParsed, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		return 0, 0, err
	}

	return latParsed, lngParsed, nil

}
