package repositories

// BaseRepository 定义基础仓库接口
type BaseRepository[T any, ID any] interface {
	Create(entity *T) error
	GetByID(id ID) (*T, error)
	GetAll(page, pageSize int, sort, order, search string) ([]T, int64, error)
	Update(entity *T) error
	Delete(id ID) error
}

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
	Search   string
}

// NewPaginationParams 创建分页参数
func NewPaginationParams(page, pageSize int, sort, order, search string) PaginationParams {
	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
		Search:   search,
	}
}
