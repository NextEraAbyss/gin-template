package controllers

import (
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// BaseController 定义基础控制器接口
type BaseController[T any, D any] interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// PaginatedResponse 分页响应结构
type PaginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// GetPaginationParams 从请求中获取分页参数
func GetPaginationParams(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 10

	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		if p, err := utils.ParseInt(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.DefaultQuery("page_size", "10"); pageSizeStr != "" {
		if ps, err := utils.ParseInt(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return page, pageSize
}

// GetSortParams 从请求中获取排序参数
func GetSortParams(c *gin.Context) (sort, order string) {
	sort = c.DefaultQuery("sort", "id")
	order = c.DefaultQuery("order", "asc")

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	return sort, order
}

// GetSearchParam 从请求中获取搜索参数
func GetSearchParam(c *gin.Context) string {
	return c.DefaultQuery("search", "")
}

// SendPaginatedResponse 发送分页响应
func SendPaginatedResponse[T any](c *gin.Context, items []T, total int64, page, pageSize int) {
	totalPages := (int(total) + pageSize - 1) / pageSize

	response := PaginatedResponse[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	utils.Success(c, response)
}
