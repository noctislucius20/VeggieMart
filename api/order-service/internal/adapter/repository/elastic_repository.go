package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"order-service/internal/core/domain/entity"
	"order-service/utils"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/labstack/gommon/log"
)

type ElasticRepositoryInterface interface {
	SearchOrderElastic(ctx context.Context, queryString entity.OrderQueryString) ([]entity.OrderEntity, int64, int64, error)
}

type elasticRepository struct {
	esClient *elasticsearch.Client
	logger   *log.Logger
}

// SearchOrderElastic implements [ElasticRepositoryInterface].
func (e *elasticRepository) SearchOrderElastic(ctx context.Context, queryString entity.OrderQueryString) ([]entity.OrderEntity, int64, int64, error) {
	offset := (queryString.Page - 1) * queryString.Limit

	statusFilter := ""
	if queryString.Status != "" {
		statusFilter = fmt.Sprintf(`{ "match": { "status": "%s" } },`, queryString.Status)
	}

	buyerIdFilter := ""
	if queryString.BuyerID != 0 {
		buyerIdFilter = fmt.Sprintf(`{ "match": { "buyer_id": %d } },`, queryString.BuyerID)
	}

	searchFilter := `{ "match_all": {} }`
	if queryString.Search != "" {
		searchFilter = fmt.Sprintf(`{ "multi_match": { "query": "%s", "fields": ["order_code", "status", "buyer_name"] } }`, queryString.Search)
	}

	mainQuery := fmt.Sprintf(`{
		"from": %d,
		"size": %d,
		"query": {
			"bool": {
				"must": [
					%s
					%s
					%s
				]
			}
		},
		"sort": [
			{ "id": "asc" }
		]
	}`, offset, queryString.Limit, statusFilter, buyerIdFilter, searchFilter)

	res, err := e.esClient.Search(
		e.esClient.Search.WithContext(ctx),
		e.esClient.Search.WithIndex("orders"),
		e.esClient.Search.WithBody(strings.NewReader(mainQuery)),
		e.esClient.Search.WithPretty(),
	)

	if err != nil {
		e.logger.Errorf("[ElasticRepository-1] SearchOrderElastic: %v", err)
		return nil, 0, 0, err
	}

	defer res.Body.Close()

	if res.IsError() && res.StatusCode == 404 {
		err := errors.New(utils.DATA_NOT_FOUND)
		e.logger.Errorf("[ElasticRepository-2] SearchOrderElastic: %v", err)
		return nil, 0, 0, err
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		e.logger.Errorf("[ElasticRepository-3] SearchOrderElastic: %v", err)
		return nil, 0, 0, err
	}

	totalData := int64(0)
	if hitsTotal, found := result["hits"].(map[string]any)["total"].(map[string]any); found {
		totalData = int64(hitsTotal["value"].(float64))
	}

	totalPage := int64(0)
	if queryString.Limit > 0 {
		totalPage = int64(math.Ceil(float64(totalData) / float64(queryString.Limit)))
	}

	orders := []entity.OrderEntity{}
	hits, found := result["hits"].(map[string]any)["hits"].([]any)
	if found {
		for _, hit := range hits {
			order := entity.OrderEntity{}
			source := hit.(map[string]any)["_source"]
			data, _ := json.Marshal(source)
			json.Unmarshal(data, &order)
			orders = append(orders, order)
		}
	}

	return orders, totalData, totalPage, nil
}

func NewElasticRepository(esClient *elasticsearch.Client, logger *log.Logger) ElasticRepositoryInterface {
	return &elasticRepository{
		esClient: esClient,
		logger:   logger,
	}
}
