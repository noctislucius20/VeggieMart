package service

import (
	"encoding/json"
	"fmt"
	"io"
	"order-service/config"
	httpclient "order-service/internal/adapter/http_client"
	"order-service/internal/core/domain/entity"
	"order-service/utils"
)

type HttpServiceInterface interface {
	HttpGetProductsService(productIds []int64, userData string) (map[int64]entity.ProductResponseEntity, error)
	HttpUpdateStockProductsService(products []entity.OrderItemEntity, userData string) error
	HttpUserByIdService(userData string) (*entity.UserResponseEntity, error)
}

type httpService struct {
	cfg        *config.Config
	httpClient httpclient.HttpClientInterface
}

// HttpUpdateStockProductsService implements [HttpServiceInterface].
func (h *httpService) HttpUpdateStockProductsService(products []entity.OrderItemEntity, userData string) error {
	baseUrlProducts := fmt.Sprintf("%s%s", h.cfg.App.ProductServiceUrl, "/internal/products/stock")

	userDataEntity := entity.JwtUserData{}
	if err := json.Unmarshal([]byte(userData), &userDataEntity); err != nil {
		return err
	}

	header := map[string]string{
		"Authorization": "Bearer " + userDataEntity.Token,
		"Content-Type":  "application/json",
	}

	req := []map[string]int64{}
	for _, rp := range products {
		req = append(req, map[string]int64{
			"product_id": rp.ProductID,
			"quantity":   rp.Quantity,
		})
	}

	payload, _ := json.Marshal(req)

	productsFetch, err := h.httpClient.CallURL("POST", baseUrlProducts, header, payload)
	if err != nil {
		return err
	}

	if productsFetch.StatusCode != 200 {
		switch productsFetch.StatusCode {
		case 409:
			return utils.ErrStockUnavailable
		default:
			return utils.ErrInternalServerError
		}
	}

	return nil
}

// HttpUserByIdService implements [HttpServiceInterface].
func (h *httpService) HttpUserByIdService(userData string) (*entity.UserResponseEntity, error) {
	baseUrlUser := fmt.Sprintf("%s%s", h.cfg.App.UserServiceUrl, "/internal/users/profile")

	userDataEntity := entity.JwtUserData{}
	if err := json.Unmarshal([]byte(userData), &userDataEntity); err != nil {
		return nil, err
	}

	header := map[string]string{
		"Authorization": "Bearer " + userDataEntity.Token,
		"Content-Type":  "application/json",
	}

	userFetch, err := h.httpClient.CallURL("GET", baseUrlUser, header, nil)
	if err != nil {
		return nil, err
	}

	if userFetch.StatusCode != 200 {
		switch userFetch.StatusCode {
		case 404:
			err := utils.ErrRelationDataNotFound
			return nil, err
		default:
			err := utils.ErrInternalServerError
			return nil, err
		}
	}

	body, err := io.ReadAll(userFetch.Body)
	if err != nil {
		return nil, err
	}
	defer userFetch.Body.Close()

	userResponse := entity.UserHttpClientResponse{}
	if err := json.Unmarshal(body, &userResponse); err != nil {
		return nil, err
	}

	return &userResponse.Data, nil
}

// HttpGetProductsService implements [HttpServiceInterface].
func (h *httpService) HttpGetProductsService(productIds []int64, userData string) (map[int64]entity.ProductResponseEntity, error) {
	baseUrlProducts := fmt.Sprintf("%s%s", h.cfg.App.ProductServiceUrl, "/internal/products/batch")

	userDataEntity := entity.JwtUserData{}
	if err := json.Unmarshal([]byte(userData), &userDataEntity); err != nil {
		return nil, err
	}

	header := map[string]string{
		"Authorization": "Bearer " + userDataEntity.Token,
		"Content-Type":  "application/json",
	}

	payload, _ := json.Marshal(map[string][]int64{"id_products": productIds})

	productsFetch, err := h.httpClient.CallURL("POST", baseUrlProducts, header, payload)
	if err != nil {
		return nil, err
	}

	if productsFetch.StatusCode != 200 {
		switch productsFetch.StatusCode {
		case 404:
			return nil, utils.ErrRelationDataNotFound
		default:
			return nil, utils.ErrInternalServerError
		}
	}

	body, err := io.ReadAll(productsFetch.Body)
	if err != nil {
		return nil, err
	}
	defer productsFetch.Body.Close()

	productsResponse := entity.ProductHttpClientResponse{}
	if err := json.Unmarshal(body, &productsResponse); err != nil {
		return nil, err
	}

	productsMap := map[int64]entity.ProductResponseEntity{}
	for _, p := range productsResponse.Data {
		productsMap[p.ID] = p
	}

	return productsMap, nil
}

func NewHttpService(cfg *config.Config, httpClient httpclient.HttpClientInterface) HttpServiceInterface {
	return &httpService{
		cfg:        cfg,
		httpClient: httpClient,
	}
}
