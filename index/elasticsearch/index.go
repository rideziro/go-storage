package elasticsearch

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Index struct {
	Data       interface{}
	DocumentID string

	ctx       context.Context
	indexName string
	client    *elasticsearch.Client
	isClone   bool
}

func NewIndex(client *elasticsearch.Client, indexName string) *Index {
	return &Index{
		indexName: indexName,
		client:    client,
	}
}

func (i *Index) SetContext(ctx context.Context) *Index {
	sData := i.getInstance()
	sData.ctx = ctx
	return sData
}

func (i *Index) SetData(data interface{}) *Index {
	sData := i.getInstance()
	sData.Data = data
	return sData
}

func (i *Index) SetDocumentID(id string) *Index {
	sData := i.getInstance()
	sData.DocumentID = id
	return sData
}

func (i *Index) Do() error {
	sData := i.getInstance()
	data, err := json.Marshal(sData.Data)
	if err != nil {
		return fmt.Errorf("es index: %w", err)
	}

	request := esapi.IndexRequest{
		Index: sData.indexName,
		Body:  bytes.NewReader(data),
	}
	if sData.DocumentID != "" {
		request.DocumentID = sData.DocumentID
	}
	res, err := request.Do(sData.ctx, sData.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e JsonObject
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	return nil
}

func (i *Index) getInstance() *Index {
	if i.isClone {
		return i
	}

	return &Index{
		Data:       i.Data,
		indexName:  i.indexName,
		DocumentID: i.DocumentID,
		ctx:        i.ctx,
		client:     i.client,
		isClone:    true,
	}
}
