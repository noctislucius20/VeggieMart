package config

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/labstack/gommon/log"
)

func (cfg Config) NewElasticSearchClient() (*elasticsearch.Client, error) {
	connectionStr := fmt.Sprintf("http://%s:%s", cfg.ElasticSearch.Host, cfg.ElasticSearch.Port)
	configElastic := elasticsearch.Config{
		Addresses: []string{connectionStr},
	}

	es, err := elasticsearch.NewClient(configElastic)
	if err != nil {
		log.Errorf("[NewElasticSearchClient-1] failed to initialize elasticsearch: %v", err)
		return nil, err
	}

	return es, nil
}
