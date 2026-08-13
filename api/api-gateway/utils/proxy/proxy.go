package proxy

import (
	"api-gateway/response"
	"api-gateway/utils/helper"
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func ForwardRequest(c echo.Context, targetUrl string) error {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return c.JSON(http.StatusBadGateway, response.ResponseFailed("invalid target url"))
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed("failed to read request body"))
	}

	req, err := http.NewRequest(c.Request().Method, target.String(), bytes.NewReader(body))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed("Failed to create request"))
	}

	for k, v := range c.Request().Header {
		for _, vl := range v {
			req.Header.Add(k, vl)
		}
	}

	helper.SetGatewayHeaders(c, &req.Header)
	req.URL.RawQuery = c.Request().URL.RawQuery

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed("Failed to forward request"))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed("Failed to read response body"))
	}

	for k, v := range resp.Header {
		for _, vl := range v {
			c.Response().Header().Add(k, vl)
		}
	}

	return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
