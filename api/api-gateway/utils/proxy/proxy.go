package proxy

import (
	"api-gateway/response"
	"api-gateway/utils/helper"
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

func ForwardRequest(c echo.Context, targetUrl string) error {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return c.JSON(http.StatusBadGateway, response.ResponseFailed("invalid target url"))
	}

	if strings.ToLower(c.Request().Header.Get("Upgrade")) == "websocket" {
		proxy := &httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				// 1. Ini pengganti fungsi originalDirector(), otomatis mengarahkan ke target URL
				// pr.SetURL(target)

				pr.Out.URL = target
				pr.Out.URL.RawQuery = pr.In.URL.RawQuery
				// 2. Pertahankan header Upgrade & Connection dari client asli (pr.In) ke request backend (pr.Out)
				pr.Out.Header.Set("Connection", pr.In.Header.Get("Connection"))
				pr.Out.Header.Set("Upgrade", pr.In.Header.Get("Upgrade"))

				// 3. Tambahkan Custom Gateway Header ke pr.Out
				// helper.SetGatewayHeaders(c, pr.Out)
			},
		}

		// Jalankan proxy WebSocket, lalu langsung return nil (selesai)
		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
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
