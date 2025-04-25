package models

import (
	"time"
)

// Article 文章模型.
// 包含文章的基本信息、内容、统计信息和时间戳
type Article struct {
	ID           uint       `gorm:"primarykey" json:"id"`                       // 文章ID.
	CreatedAt    time.Time  `gorm:"type:datetime;not null" json:"created_at"`   // 创建时间.
	UpdatedAt    time.Time  `gorm:"type:datetime;not null" json:"updated_at"`   // 更新时间.
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at"`                    // 删除时间.
	Title        string     `gorm:"type:varchar(255);not null" json:"title"`    // 文章标题.
	Content      string     `gorm:"type:text;not null" json:"content"`          // 文章内容.
	AuthorID     uint       `gorm:"not null;index" json:"author_id"`            // 作者ID.
	Author       User       `gorm:"foreignKey:AuthorID" json:"author"`          // 作者信息.
	Status       int        `gorm:"type:tinyint;default:1;index" json:"status"` // 文章状态：1-已发布，0-草稿.
	PublishedAt  time.Time  `gorm:"type:timestamp;index" json:"published_at"`   // 发布时间
	ViewCount    int        `gorm:"default:0" json:"view_count"`                // 浏览次数.
	LikeCount    int        `gorm:"default:0" json:"like_count"`                // 点赞次数.
	CommentCount int        `gorm:"default:0" json:"comment_count"`             // 评论次数.
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}

// ArticleQueryDTO 文章查询参数
type ArticleQueryDTO struct {
	Page     int    `form:"page" binding:"required,min=1"`      // 页码，从1开始
	PageSize int    `form:"page_size" binding:"required,min=1"` // 每页数量
	Keyword  string `form:"keyword"`                            // 搜索关键词
	AuthorID uint   `form:"author_id"`                          // 作者ID
	Status   *int   `form:"status"`                             // 文章状态
}

// ArticleListRequest 文章列表请求
type ArticleListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`      // 页码，从1开始
	PageSize int    `form:"page_size" binding:"required,min=1"` // 每页数量
	Keyword  string `form:"keyword"`                            // 搜索关键词
	AuthorID uint   `form:"author_id"`                          // 作者ID
	Status   *int   `form:"status"`                             // 文章状态
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"` // 文章标题，1-200字符
	Content string `json:"content" binding:"required,min=1"`       // 文章内容，不能为空
	Status  int    `json:"status" binding:"oneof=0 1"`             // 文章状态：1-发布，0-草稿
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=200"` // 文章标题，可选，1-200字符
	Content string `json:"content" binding:"omitempty,min=1"`       // 文章内容，可选，不能为空
	Status  *int   `json:"status" binding:"omitempty,oneof=0 1"`    // 文章状态，可选：1-发布，0-草稿
}

// ArticleResponse 文章信息响应
type ArticleResponse struct {
	ID           uint      `json:"id"`            // 文章ID
	Title        string    `json:"title"`         // 文章标题
	Content      string    `json:"content"`       // 文章内容
	AuthorID     uint      `json:"author_id"`     // 作者ID
	Author       User      `json:"author"`        // 作者信息
	Status       int       `json:"status"`        // 文章状态
	ViewCount    int       `json:"view_count"`    // 浏览次数
	LikeCount    int       `json:"like_count"`    // 点赞次数
	CommentCount int       `json:"comment_count"` // 评论次数
	PublishedAt  time.Time `json:"published_at"`  // 发布时间
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`    // 更新时间
}

// ArticleListResponse 文章列表响应
type ArticleListResponse struct {
	Total int64              `json:"total"` // 总数
	Items []*ArticleResponse `json:"items"` // 文章列表
}

// ToResponse 将Article模型转换为ArticleResponse
func (a *Article) ToResponse() *ArticleResponse {
	return &ArticleResponse{
		ID:           a.ID,
		Title:        a.Title,
		Content:      a.Content,
		AuthorID:     a.AuthorID,
		Author:       a.Author,
		Status:       a.Status,
		ViewCount:    a.ViewCount,
		LikeCount:    a.LikeCount,
		CommentCount: a.CommentCount,
		PublishedAt:  a.PublishedAt,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
