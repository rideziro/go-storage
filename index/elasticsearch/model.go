package elasticsearch

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/rideziro/go-storage/index"
)

type SearchQuery struct {
	Must    JsonArray `json:"must,omitempty"`
	Should  JsonArray `json:"should,omitempty"`
	MustNot JsonArray `json:"must_not,omitempty"`
	Filter  JsonArray `json:"filter,omitempty"`
}

type SearchResponse struct {
	Took int64 `json:"took"`
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []*SearchHit `json:"hits"`
	} `json:"hits"`
	Aggregations jsoniter.RawMessage `json:"aggregations"`
}

type SearchHit struct {
	Source  jsoniter.RawMessage `json:"_source"`
	Score   float64             `json:"_score,omitempty"`
	Index   string              `json:"_index"`
	Type    string              `json:"_type"`
	Version int64               `json:"_version,omitempty"`
	Sort    index.Sort          `json:"sort,omitempty"`
}
