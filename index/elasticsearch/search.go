package elasticsearch

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/rideziro/go-storage/index"
	"io"
)

const (
	ListSize       = 15
	MinimumArticle = 1
)

type Search struct {
	Query         SearchQuery          `json:"query"`
	Sort          JsonArray            `json:"sort"`
	Aggregations  JsonObject           `json:"aggs"`
	FunctionScore JsonObject           `json:"function_score"`
	Size          int                  `json:"size"`
	Paginator     index.PaginatorIndex `json:"paginator"`

	multipleQuery []SearchQuery
	indexNames    []string
	ctx           context.Context
	client        *elasticsearch.Client
	isClone       bool
}

func NewSearch(client *elasticsearch.Client, indexNames ...string) *Search {
	return &Search{
		indexNames: indexNames,
		client:     client,
	}
}

func (s *Search) WithContext(ctx context.Context) *Search {
	sData := s.getInstance()
	sData.ctx = ctx
	return sData
}

func (s *Search) SetSize(size int) *Search {
	sData := s.getInstance()
	sData.Size = size
	return sData
}

func (s *Search) SetPaginator(paginator index.PaginatorIndex) *Search {
	sData := s.getInstance()
	sData.Paginator = paginator
	return sData
}

func (s *Search) AddSort(sort interface{}) *Search {
	sData := s.getInstance()
	sData.Sort = append(sData.Sort, sort)
	return sData
}

func (s *Search) AddMustQuery(query interface{}) *Search {
	sData := s.getInstance()
	sData.Query.Must = append(sData.Query.Must, query)
	return sData
}

func (s *Search) AddShouldQuery(query interface{}) *Search {
	sData := s.getInstance()
	sData.Query.Should = append(sData.Query.Should, query)
	return sData
}

func (s *Search) AddMustNotQuery(query interface{}) *Search {
	sData := s.getInstance()
	sData.Query.MustNot = append(sData.Query.MustNot, query)
	return sData
}

func (s *Search) AddFilterQuery(query interface{}) *Search {
	sData := s.getInstance()
	sData.Query.Filter = append(sData.Query.Filter, query)
	return sData
}

func (s *Search) SetAggregationQuery(query interface{}) *Search {
	sData := s.getInstance()
	sData.Aggregations = JsonObject{"aggs": query}
	return sData
}

func (s *Search) SetFunctionScore(query interface{}) *Search {
	sData := s.getInstance()
	sData.FunctionScore = JsonObject{"function_score": query}
	return sData
}

func (s *Search) AddMultipleQuery(query SearchQuery) *Search {
	sData := s.getInstance()
	sData.multipleQuery = append(sData.multipleQuery, query)
	return sData
}

func (s *Search) ToData() (io.Reader, error) {
	sData := s.getInstance()
	script := JsonObject{}
	if len(sData.FunctionScore) > 0 {
		script["query"] = sData.FunctionScore
	} else {
		script["query"] = JsonObject{
			"bool": sData.Query,
		}
	}

	if len(sData.multipleQuery) > 0 {
		boolArray := JsonArray{}
		for _, query := range sData.multipleQuery {
			boolArray = append(boolArray, JsonObject{
				"bool": query,
			})
		}
		script["query"] = JsonObject{
			"bool": JsonObject{
				"filter": JsonObject{
					"bool": sData.Query,
				},
				"should": boolArray,
			},
		}
	}

	sort, err := sData.Paginator.ToSort()
	if err == nil && sort != nil {
		script["search_after"] = sort
	}
	if len(sData.Sort) > 0 {
		script["sort"] = sData.Sort
	}
	if sData.Size >= 0 {
		script["size"] = s.Size
	}
	if len(sData.Aggregations) > 0 {
		script["aggs"] = sData.Aggregations["aggs"]
	}

	data, err := json.Marshal(script)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

func (s *Search) Do() (SearchResponse, index.PaginatorIndex, error) {
	sData := s.getInstance()
	var (
		response SearchResponse
		next     index.PaginatorIndex

		ctx    = sData.ctx
		client = sData.client
	)

	if sData.Size == 0 {
		sData.Size = ListSize
	} else if sData.Size < 0 {
		sData.Size = 0
	}
	size := sData.Size

	data, err := sData.ToData()
	if err != nil {
		return SearchResponse{}, "", err
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(sData.indexNames...),
		client.Search.WithBody(data),
	)
	if err != nil {
		return response, next, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var errorRes JsonObject
		if err := json.NewDecoder(res.Body).Decode(&errorRes); err != nil {
			return response, next, err
		}
		return response, next, fmt.Errorf("[%s] %s: %s", res.Status(), errorRes["error"].(map[string]interface{})["type"], errorRes["error"].(map[string]interface{})["reason"])
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return response, next, err
	}
	hits := response.Hits.Hits
	if len(hits) > 0 {
		if len(hits) == size {
			sort := hits[len(hits)-1].Sort
			next = index.NewPaginatorIndex(sort.String())
		}
	}
	return response, next, nil
}

func (s *Search) getInstance() *Search {
	if s.isClone {
		return s
	}

	return &Search{
		Query:         s.Query,
		Sort:          s.Sort,
		Size:          s.Size,
		Paginator:     s.Paginator,
		multipleQuery: s.multipleQuery,
		indexNames:    s.indexNames,
		ctx:           s.ctx,
		client:        s.client,
		isClone:       true,
	}
}
