package elasticsearch

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ESScript struct {
	Source string      `json:"source"`
	Lang   string      `json:"lang"`
	Params interface{} `json:"params,omitempty"`
	Query  interface{} `json:"-"`

	id        string
	conflict  bool
	indexName string
	ctx       context.Context
	client    *elasticsearch.Client
	isClone   bool
}

func NewScriptPainless(client *elasticsearch.Client, index string) *ESScript {
	return &ESScript{
		Lang:      "painless",
		indexName: index,
		client:    client,
	}
}

func (s *ESScript) WithContext(ctx context.Context) *ESScript {
	sData := s.getInstance()
	sData.ctx = ctx
	return sData
}

func (s *ESScript) SetID(id string) *ESScript {
	sData := s.getInstance()
	sData.id = id
	return sData
}

func (s *ESScript) SetSource(source string) *ESScript {
	sData := s.getInstance()
	sData.Source = source
	return sData
}

func (s *ESScript) SetConflictProceed(conflict bool) *ESScript {
	sData := s.getInstance()
	sData.conflict = conflict
	return sData
}

func (s *ESScript) SetParams(params interface{}) *ESScript {
	sData := s.getInstance()
	sData.Params = params
	return sData
}

func (s *ESScript) SetQuery(query interface{}) *ESScript {
	sData := s.getInstance()
	sData.Query = query
	return sData
}

func (s *ESScript) ToData() (io.Reader, error) {
	sData := s.getInstance()
	script := JsonObject{}
	if sData.Query != nil {
		script["query"] = sData.Query
	}
	script["script"] = sData
	data, err := json.Marshal(script)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

func (s *ESScript) Do() error {
	sData := s.getInstance()
	var (
		res    *esapi.Response
		client = sData.client
		ctx    = sData.ctx
	)

	data, err := sData.ToData()
	if err != nil {
		return err
	}

	conflict := "abort"
	if sData.conflict {
		conflict = "proceed"
	}

	if sData.id != "" {
		res, err = client.Update(sData.indexName, sData.id, data, client.Update.WithContext(ctx))
	} else {
		res, err = client.UpdateByQuery([]string{sData.indexName},
			client.UpdateByQuery.WithContext(ctx),
			client.UpdateByQuery.WithBody(data),
			client.UpdateByQuery.WithConflicts(conflict),
		)
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("[%v]", res.StatusCode)
	}

	return nil
}

func (s *ESScript) getInstance() *ESScript {
	if s.isClone {
		return s
	}

	return &ESScript{
		Source:    s.Source,
		Lang:      s.Lang,
		Params:    s.Params,
		Query:     s.Query,
		id:        s.id,
		conflict:  s.conflict,
		indexName: s.indexName,
		ctx:       s.ctx,
		client:    s.client,
		isClone:   true,
	}
}
