package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"order-service/config"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type HttpClientInterface interface {
	Connect()
	CallURL(method string, url string, header map[string]string, payload []byte) (*http.Response, error)
	RoundTrip(req *http.Request) (*http.Response, error)
}

type options struct {
	timeout int
	http    *http.Client
	logger  echo.Logger
}

// RoundTrip implements [HttpClientInterface].
func (o *options) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		o.logger.Infof("Request failed: %v", err.Error())
		return nil, err
	}

	o.logger.Infof("%s %d %s", req.Method, resp.StatusCode, req.URL)

	return resp, nil
}

// CallURL implements [HttpClientInterface].
func (o *options) CallURL(method string, url string, header map[string]string, payload []byte) (*http.Response, error) {
	o.Connect()
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		o.logger.Errorj(log.JSON{
			"message": "[HttpClient-1] CallURL: failed to pre request client http",
			"error":   err.Error(),
		})
		return nil, err
	}

	if len(header) > 0 {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	resp, err := o.http.Do(req)
	if err != nil {
		o.logger.Errorj(log.JSON{
			"message": "[HttpClient-2] CallURL: failed to call request http",
			"error":   err.Error(),
		})
		return nil, err
	}

	return resp, nil
}

// Connect implements [HttpClientInterface].
func (o *options) Connect() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetHeader(`${time_rfc3339} | ${level} | ${short_file}:${line} |`)

	httpClient := &http.Client{
		Timeout:   time.Duration(o.timeout) * time.Second,
		Transport: &options{logger: e.Logger},
	}

	o.http = httpClient
	o.logger = e.Logger
}

func NewHttpClient(cfg *config.Config) HttpClientInterface {
	opt := new(options)
	opt.timeout = cfg.App.ServerTimeout
	return opt
}
