package service

import (
	"encoding/json"
	"errors"
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
	HttpOrdersAllService(orderIds []int64, userData string) (map[int64]entity.OrderDetailResponseEntity, error)
	HttpOrderByIdService(orderId int64, userData string) (*entity.OrderDetailResponseEntity, error)
	HttpOrderIdByOrderCodePublicService(orderCode string) (uint, error)
}

type httpService struct {
	cfg        *config.Config
	httpClient httpclient.HttpClientInterface
}

// HttpOrderIdByOrderCodePublicService implements [HttpServiceInterface].
func (h *httpService) HttpOrderIdByOrderCodePublicService(orderCode string) (uint, error) {
	baseUrlOrder := fmt.Sprintf("%s%s", h.cfg.App.OrderServiceUrl, "/public/orders/"+orderCode+"/code")

	header := map[string]string{
		"Content-Type": "application/json",
	}

	orderFetch, err := h.httpClient.CallURL("GET", baseUrlOrder, header, nil)
	if err != nil {
		return 0, err
	}

	if orderFetch.StatusCode != 200 {
		err := errors.New(strconv.Itoa(orderFetch.StatusCode))
		return 0, err
	}

	body, err := io.ReadAll(orderFetch.Body)
	if err != nil {
		return 0, err
	}
	defer orderFetch.Body.Close()

	orderResponse := entity.OrderIDHttpResponseEntity{}
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		return 0, err
	}

	return orderResponse.Data.OrderID, nil
}

// HttpOrderByIdService implements [HttpServiceInterface].
func (h *httpService) HttpOrderByIdService(orderId int64, userData string) (*entity.OrderDetailResponseEntity, error) {
	userDataEntity := entity.JwtUserData{}
	if err := json.Unmarshal([]byte(userData), &userDataEntity); err != nil {
		return nil, err
	}

	baseUrlOrder := fmt.Sprintf("%s%s", h.cfg.App.OrderServiceUrl, "/admin/orders/"+strconv.Itoa(int(orderId)))
	if strings.ToLower(userDataEntity.RoleName) == "customer" {
		baseUrlOrder = fmt.Sprintf("%s%s", h.cfg.App.OrderServiceUrl, "/auth/orders/"+strconv.Itoa(int(orderId)))
	}

	header := map[string]string{
		"Authorization": "Bearer " + userDataEntity.Token,
		"Content-Type":  "application/json",
	}

	orderFetch, err := h.httpClient.CallURL("GET", baseUrlOrder, header, nil)
	if err != nil {
		return nil, err
	}

	if orderFetch.StatusCode != 200 {
		switch orderFetch.StatusCode {
		case 404:
			err := errors.New(utils.DATA_NOT_FOUND)
			return nil, err
		default:
			err := errors.New(utils.INTERNAL_SERVER_ERROR)
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

func (h *httpService) HttpOrdersAllService(orderIds []int64, userData string) (map[int64]entity.OrderDetailResponseEntity, error) {
	baseUrlOrder := fmt.Sprintf("%s%s", h.cfg.App.OrderServiceUrl, "/auth/orders/batch")

	userDataEntity := entity.JwtUserData{}
	if err := json.Unmarshal([]byte(userData), &userDataEntity); err != nil {
		return nil, err
	}

	header := map[string]string{
		"Authorization": "Bearer " + userDataEntity.Token,
		"Content-Type":  "application/json",
	}

	payload, err := json.Marshal(map[string][]int64{"id_orders": orderIds})
	if err != nil {
		return nil, err
	}

	orderFetch, err := h.httpClient.CallURL("POST", baseUrlOrder, header, payload)
	if err != nil {
		return nil, err
	}

	if orderFetch.StatusCode != 200 {
		err := errors.New(strconv.Itoa(orderFetch.StatusCode))
		return nil, err
	}

	body, err := io.ReadAll(orderFetch.Body)
	if err != nil {
		return nil, err
	}
	defer orderFetch.Body.Close()

	orderResponse := entity.OrderHttpClientResponseList{}
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		return nil, err
	}

	ordersMap := map[int64]entity.OrderDetailResponseEntity{}
	for _, o := range orderResponse.Data {
		ordersMap[o.ID] = o
	}

	return ordersMap, nil
}

func NewHttpService(cfg *config.Config, httpClient httpclient.HttpClientInterface) HttpServiceInterface {
	return &httpService{
		cfg:        cfg,
		httpClient: httpClient,
	}
}
