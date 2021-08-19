package storage

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	PaginatorDefaultSize = 15
)

type Paginator struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

var (
	ErrUnableToDecode    = errors.New("unable to decode given string")
	ErrInvalidPagination = errors.New("invalid pagination string given")
)

func (p *Paginator) String() string {
	data := []int{p.Page, p.Size}
	return base64.StdEncoding.EncodeToString([]byte(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(data)), ","), "[]")))
}

func (p *Paginator) MarshalJSON() ([]byte, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(&struct {
		Page string `json:"page"`
	}{
		p.String(),
	})
}

func (p *Paginator) IncreasePage() {
	p.Page++
}

func NewPaginator(data string) (Paginator, error) {
	var paginator Paginator
	if data == "" {
		return paginator, nil
	}
	decodeString, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return paginator, ErrUnableToDecode
	}
	allValues := strings.Split(string(decodeString), ",")
	if len(allValues) != 2 {
		return paginator, ErrInvalidPagination
	}
	page, err := strconv.Atoi(allValues[0])
	if err != nil {
		return paginator, ErrInvalidPagination
	}
	size, err := strconv.Atoi(allValues[1])
	if err != nil {
		return paginator, ErrInvalidPagination
	}

	paginator.Page = page
	paginator.Size = size
	return paginator, nil
}
