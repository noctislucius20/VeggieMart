package service

import (
	"encoding/json"
	"fmt"
	"io"
	"payment-service/config"
	httpclient "payment-service/internal/adapter/http_client"
	"payment-service/internal/core/domain/entity"
	"payment-service/utils"
	"strconv"
	"strings"
)

type HttpServiceInterface interface {
	HttpOrderByIdService(orderId int64, jwtUserData entity.JwtUserData, roleName string) (*entity.OrderDetailResponseEntity, error)
}

type httpService struct {
	cfg        *config.Config
	httpClient httpclient.HttpClientInterface
}

// HttpOrderByIdService implements [HttpServiceInterface].
func (h *httpService) HttpOrderByIdService(orderId int64, jwtUserData entity.JwtUserData, roleName string) (*entity.OrderDetailResponseEntity, error) {
	baseUrlOrder := fmt.Sprintf("%s%s", h.cfg.App.OrderServiceUrl, "/internal/orders/"+strconv.Itoa(int(orderId))+"/admin")
	if strings.ToLower(roleName) == "customer" {
		baseUrlOrder = fmt.Sprintf("%s%s", h.cfg.App.OrderServiceUrl, "/internal/orders/"+strconv.Itoa(int(orderId)))
	}

	header := map[string]string{
		"Authorization": "Bearer " + jwtUserData.Token,
		"Content-Type":  "application/json",
	}

	orderFetch, err := h.httpClient.CallURL("GET", baseUrlOrder, header, nil)
	if err != nil {
		return nil, err
	}

	if orderFetch.StatusCode != 200 {
		switch orderFetch.StatusCode {
		case 403:
			err := utils.ErrAccessForbidden
			return nil, err
		case 404:
			err := utils.ErrRelationDataNotFound
			return nil, err
		default:
			err := utils.ErrInternalServerError
			return nil, err
		}
	}

	body, err := io.ReadAll(orderFetch.Body)
	if err != nil {
		return nil, err
	}

	defer orderFetch.Body.Close()

	orderResponse := entity.OrderHttpClientResponse{}
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		return nil, err
	}

	return &orderResponse.Data, nil
}

func NewHttpService(cfg *config.Config, httpClient httpclient.HttpClientInterface) HttpServiceInterface {
	return &httpService{
		cfg:        cfg,
		httpClient: httpClient,
	}
}
