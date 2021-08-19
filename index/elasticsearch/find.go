package elasticsearch

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
)

type Find struct {
	indexName string
	ctx       context.Context
	client    *elasticsearch.Client

	isClone bool
}

func NewFind(client *elasticsearch.Client, indexName string) *Find {
	return &Find{
		indexName: indexName,
		client:    client,
	}
}

func (s *Find) getInstance() *Find {
	if s.isClone {
		return s
	}

	return &Find{
		indexName: s.indexName,
		ctx:       s.ctx,
		client:    s.client,
		isClone:   true,
	}
}

func (s *Find) WithContext(ctx context.Context) *Find {
	sData := s.getInstance()
	sData.ctx = ctx
	return sData
}

func (s *Find) Do(id string) (SearchHit, error) {
	sData := s.getInstance()
	var (
		response SearchHit
		ctx      = sData.ctx
		client   = sData.client
	)
	res, err := client.Get(sData.indexName, id,
		client.Get.WithContext(ctx))
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return response, ErrNotFound
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return response, err
	}
	return response, nil
}
