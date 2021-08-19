package scopes

import (
	"github.com/rideziro/go-storage/storage"
	"gorm.io/gorm"
)

const (
	MaxPaginationSize     = 50
	DefaultPaginationSize = 15
)

func Paginate(paginator *storage.Paginator) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := paginator.Page
		pageSize := paginator.Size
		switch {
		case pageSize > MaxPaginationSize:
			pageSize = MaxPaginationSize
		case pageSize <= 0:
			pageSize = DefaultPaginationSize
		}

		if page <= 0 {
			page = 1
		}

		offset := (page - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		paginator.Page = page
		paginator.Size = pageSize
		paginator.IncreasePage()

		return db.Offset(offset).Limit(pageSize)
	}
}
