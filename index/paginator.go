package index

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrUnableToDecode = errors.New("unable to decode given string")
)

type Sort []interface{}
type PaginatorIndex string

func (p *Sort) String() string {
	if p == nil {
		return ""
	}
	var allValues []string
	for _, datum := range *p {
		allValues = append(allValues, fmt.Sprintf("%v", datum))
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(allValues, ",")))
}

func (p PaginatorIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Next string `json:"next"`
	}{
		string(p),
	})
}

func (p *PaginatorIndex) ToSort() (Sort, error) {
	var sort Sort
	data := p
	if data == nil || *data == "" {
		return sort, nil
	}
	decodeString, err := base64.StdEncoding.DecodeString(string(*data))
	if err != nil {
		return sort, ErrUnableToDecode
	}
	allValues := strings.Split(string(decodeString), ",")
	for _, datum := range allValues {
		float, err := strconv.ParseFloat(datum, 64)
		if err == nil {
			sort = append(sort, float)
			continue
		}
		sort = append(sort, datum)
	}
	return sort, nil
}

func NewPaginatorIndex(data string) PaginatorIndex {
	return PaginatorIndex(data)
}
