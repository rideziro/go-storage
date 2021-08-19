package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPaginator(t *testing.T) {
	var defaultPaginator Paginator
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want Paginator
		err  error
	}{
		{
			"Empty",
			args{data: ""},
			defaultPaginator,
			nil,
		},
		{
			"Undecodable",
			args{data: "invalid string"},
			defaultPaginator,
			ErrUnableToDecode,
		},
		{
			"Page 1 Size 15",
			args{data: "MSwxNQ=="},
			Paginator{1, 15},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPaginator(tt.args.data)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestPaginator_String(t *testing.T) {
	tests := []struct {
		name string
		p    Paginator
		want string
	}{
		{
			"Page 1 Size 15",
			Paginator{1, 15},
			"MSwxNQ==",
		},
		{
			"Page 200 Size 30",
			Paginator{200, 30},
			"MjAwLDMw",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.p.String(), tt.want)
		})
	}
}

func TestPaginator_IncreasePage(t *testing.T) {
	tests := []struct {
		name string
		p    Paginator
		want Paginator
	}{
		{
			"Increase page",
			Paginator{1, 15},
			Paginator{2, 15},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.IncreasePage()
			assert.Equal(t, tt.p, tt.want)
		})
	}
}

func TestPaginator_MarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		paginator *Paginator
		want      []byte
		err       error
	}{
		{
			"Default",
			&Paginator{
				Page: 1,
				Size: PaginatorDefaultSize,
			},
			[]byte("{\"page\":\"MSwxNQ==\"}"),
			nil,
		},
		{
			"Nil Page",
			nil,
			nil,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.paginator.MarshalJSON()
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, got, tt.want)
		})
	}
}
