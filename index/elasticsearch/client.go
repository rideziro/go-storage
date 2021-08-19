package elasticsearch

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	jsoniniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

var (
	json = jsoniniter.ConfigCompatibleWithStandardLibrary
)

func NewElasticsearchClient() (*elasticsearch.Client, error) {
	username := viper.GetString("ELASTICSEARCH_USERNAME")
	password := viper.GetString("ELASTICSEARCH_PASSWORD")

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	ping, err := client.Ping()
	if err != nil {
		return nil, err
	}
	if ping.IsError() {
		return nil, fmt.Errorf("can't connect: [%v] %v", ping.StatusCode, ping.String())
	}
	return client, nil
}

type JsonObject map[string]interface{}
type JsonArray []interface{}

type ESClient struct {
	ESFind   *Find
	ESScript *ESScript
	ESIndex  *Index
	ESSearch *Search

	Client    *elasticsearch.Client
	IndexName string
}

func NewESClient(client *elasticsearch.Client, indexName string) ESClient {
	find := NewFind(client, indexName)
	painless := NewScriptPainless(client, indexName)
	index := NewIndex(client, indexName)
	search := NewSearch(client, indexName)

	return ESClient{
		ESFind:   find,
		ESScript: painless,
		ESIndex:  index,
		ESSearch: search,

		Client:    client,
		IndexName: indexName,
	}
}
