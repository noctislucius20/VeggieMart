package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"order-service/config"
	httpclient "order-service/internal/adapter/http_client"
	"order-service/internal/core/domain/entity"
	"order-service/utils"
	"strconv"
)

type HttpServiceInterface interface {
	HttpProductsAllService(productIds []int64, userData string) (map[int64]entity.ProductResponseEntity, error)
	HttpUsersAllAdminService(userId []int64, userData string) (map[int64]entity.UserResponseEntity, error)
	HttpUserByIdAdminService(userId int64, userData string) (*entity.UserResponseEntity, error)
	HttpUserByIdService(userData string) (*entity.UserResponseEntity, error)
}

type httpService struct {
	cfg        *config.Config
	httpClient httpclient.HttpClientInterface
}

// HttpUserByIdService implements [HttpServiceInterface].
func (h *httpService) HttpUserByIdService(userData string) (*entity.UserResponseEntity, error) {
	baseUrlUser := fmt.Sprintf("%s%s", h.cfg.App.UserServiceUrl, "/auth/profile")

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
			err := errors.New(utils.DATA_NOT_FOUND)
			return nil, err
		default:
			err := errors.New(utils.INTERNAL_SERVER_ERROR)
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

// HttpUserByIdAdminService implements [HttpServiceInterface].
func (h *httpService) HttpUserByIdAdminService(userId int64, userData string) (*entity.UserResponseEntity, error) {
	baseUrlUser := fmt.Sprintf("%s%s", h.cfg.App.UserServiceUrl, "/admin/customers/"+strconv.Itoa(int(userId)))

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
			err := errors.New(utils.DATA_NOT_FOUND)
			return nil, err
		default:
			err := errors.New(utils.INTERNAL_SERVER_ERROR)
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

// HttpUsersAllAdminService implements [HttpServiceInterface].
func (h *httpService) HttpUsersAllAdminService(userIds []int64, userData string) (map[int64]entity.UserResponseEntity, error) {
	baseUrlUser := fmt.Sprintf("%s%s", h.cfg.App.UserServiceUrl, "/admin/customers/batch")

	userDataEntity := entity.JwtUserData{}
	if err := json.Unmarshal([]byte(userData), &userDataEntity); err != nil {
		return nil, err
	}

	header := map[string]string{
		"Authorization": "Bearer " + userDataEntity.Token,
		"Content-Type":  "application/json",
	}

	payload, err := json.Marshal(map[string][]int64{"id_users": userIds})
	if err != nil {
		return nil, err
	}

	userFetch, err := h.httpClient.CallURL("POST", baseUrlUser, header, payload)
	if err != nil {
		return nil, err
	}

	if userFetch.StatusCode != 200 {
		err := errors.New(strconv.Itoa(userFetch.StatusCode))
		return nil, err
	}

	body, err := io.ReadAll(userFetch.Body)
	if err != nil {
		return nil, err
	}
	defer userFetch.Body.Close()

	userResponse := entity.UsersHttpClientResponse{}
	if err := json.Unmarshal(body, &userResponse); err != nil {
		return nil, err
	}

	usersMap := map[int64]entity.UserResponseEntity{}
	for _, p := range userResponse.Data {
		usersMap[p.ID] = p
	}

	return usersMap, nil
}

// HttpProductsAllService implements [HttpServiceInterface].
func (h *httpService) HttpProductsAllService(productIds []int64, userData string) (map[int64]entity.ProductResponseEntity, error) {
	baseUrlProducts := fmt.Sprintf("%s%s", h.cfg.App.ProductServiceUrl, "/auth/products/batch")

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
		err := errors.New(strconv.Itoa(productsFetch.StatusCode))
		return nil, err
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
