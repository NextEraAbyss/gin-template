package services

import (
	"gitee.com/NextEraAbyss/gin-template/repositories"
)

// BaseService 定义基础服务接口
type BaseService[T any, ID any] interface {
	Create(entity *T) error
	GetByID(id ID) (*T, error)
	GetAll(page, pageSize int, sort, order, search string) ([]T, int64, error)
	Update(entity *T) error
	Delete(id ID) error
}

// 使用仓库层定义的PaginationParams
type PaginationParams = repositories.PaginationParams

// NewPaginationParams 创建分页参数
func NewPaginationParams(page, pageSize int, sort, order, search string) PaginationParams {
	return repositories.NewPaginationParams(page, pageSize, sort, order, search)
}
