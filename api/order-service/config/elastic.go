package config

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/labstack/gommon/log"
)

func (cfg Config) NewElasticsearchClient() (*elasticsearch.Client, error) {
	connectionStr := fmt.Sprintf("http://%s:%s", cfg.Elasticsearch.Host, cfg.Elasticsearch.Port)
	configElastic := elasticsearch.Config{
		Addresses: []string{connectionStr},
	}

	es, err := elasticsearch.NewClient(configElastic)
	if err != nil {
		log.Errorf("[NewElasticsearchClient-1] failed to initialize elasticsearch: %v", err)
		return nil, err
	}

	return es, nil
}
