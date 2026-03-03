package conv

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LatLngToString(f float64) string {
	if f == 0 {
		return ""
	}

	return strconv.FormatFloat(f, 'g', -1, 64)
}

func ToJSON(data any) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonData, err
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
