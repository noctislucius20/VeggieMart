package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"product-service/internal/core/domain/entity"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/labstack/gommon/log"
)

type ElasticRepositoryInterface interface {
	SearchProductElastic(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error)
}

type elasticRepository struct {
	esClient *elasticsearch.Client
	logger   *log.Logger
}

// SearchProductElastic implements [ElasticRepositoryInterface].
func (e *elasticRepository) SearchProductElastic(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
	offset := (query.Page - 1) * query.Limit

	categoryFilter := ""
	if query.CategoryID != 0 {
		categoryFilter = fmt.Sprintf(`{ "match": { "category_id": "%d" } },`, query.CategoryID)
	}

	priceFilter := ""
	if query.StartPrice > 0 && query.EndPrice > 0 {
		priceFilter = fmt.Sprintf(`{ "range": { "regular_price": { "gte": %d, "lte": %d } } }`, query.StartPrice, query.EndPrice)
	}

	searchFilter := `{ "match_all": {} }`
	if query.Search != "" {
		searchFilter = fmt.Sprintf(`{ "multi_match": { "query": "%s", "fields": ["name", "description"] } }`, query.Search)
	}

	mainQuery := fmt.Sprintf(`{
		"from": %d,
		"size": %d,
		"query": {
			"bool": {
				"must": [
					%s
					%s
				],
				"filter": [
					%s
				]
			}
		},
		"sort": [
			{ "id": "asc" }
		]
	}`, offset, query.Limit, categoryFilter, searchFilter, priceFilter)

	res, err := e.esClient.Search(
		e.esClient.Search.WithContext(ctx),
		e.esClient.Search.WithIndex("products"),
		e.esClient.Search.WithBody(strings.NewReader(mainQuery)),
		e.esClient.Search.WithPretty(),
	)
	if err != nil {
		e.logger.Errorf("[ProductRepository-1] SearchProducts: %v", err)
		return nil, 0, 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err := errors.New(res.Status())
		e.logger.Errorf("[ProductRepository-2] SearchProducts: %v", err)
		return nil, 0, 0, err
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		e.logger.Errorf("[ProductRepository-3] SearchProducts: %v", err)
		return nil, 0, 0, err
	}

	totalData := int64(0)
	if hitsTotal, found := result["hits"].(map[string]any)["total"].(map[string]any); found {
		totalData = int64(hitsTotal["value"].(float64))
	}

	totalPage := int64(0)
	if query.Limit > 0 {
		totalPage = int64(math.Ceil(float64(totalData) / float64(query.Limit)))
	}

	products := []entity.ProductEntity{}
	hits, found := result["hits"].(map[string]any)["hits"].([]any)
	if found {
		for _, hit := range hits {
			product := entity.ProductEntity{}

			source := hit.(map[string]any)["_source"]
			data, _ := json.Marshal(source)
			json.Unmarshal(data, &product)

			products = append(products, product)
		}
	}

	return products, totalData, totalPage, nil
}

func NewElasticRepository(esClient *elasticsearch.Client, logger *log.Logger) ElasticRepositoryInterface {
	return &elasticRepository{
		esClient: esClient,
		logger:   logger,
	}
}
