package elasticsearch

import (
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	jsoniniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

var (
	json = jsoniniter.ConfigCompatibleWithStandardLibrary
)

func NewElasticsearchClient() (*elasticsearch.Client, error) {
	url := viper.GetString("ELASTICSEARCH_URL")
	username := viper.GetString("ELASTICSEARCH_USERNAME")
	password := viper.GetString("ELASTICSEARCH_PASSWORD")
	caCert, _ := os.ReadFile(viper.GetString("ELASTICSEARCH_CA_CERT"))

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
		Username:  username,
		Password:  password,
		CACert:    caCert,
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
