package repositories

import (
	"gitee.com/NextEraAbyss/gin-template/models"
	"gorm.io/gorm"
)

// ArticleRepository 文章仓库接口
// 定义了文章相关的数据库操作
type ArticleRepository interface {
	Create(article *models.Article) error                                // 创建文章
	Update(article *models.Article) error                                // 更新文章
	Delete(id uint) error                                                // 删除文章
	FindByID(id uint) (*models.Article, error)                           // 根据ID查找文章
	List(query *models.ArticleQueryDTO) ([]models.Article, int64, error) // 获取文章列表
	IncrementViewCount(id uint) error                                    // 增加文章浏览次数
	IncrementLikeCount(id uint) error                                    // 增加文章点赞次数
	IncrementCommentCount(id uint) error                                 // 增加文章评论次数
}

// articleRepository 文章仓库实现
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓库实例
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// Create 创建文章
func (r *articleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

// Update 更新文章
func (r *articleRepository) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

// Delete 删除文章
func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Article{}, id).Error
}

// FindByID 根据ID查找文章
func (r *articleRepository) FindByID(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Author").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// List 获取文章列表
func (r *articleRepository) List(query *models.ArticleQueryDTO) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	db := r.db.Model(&models.Article{})

	// 添加查询条件
	if query.Keyword != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.AuthorID > 0 {
		db = db.Where("author_id = ?", query.AuthorID)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (query.Page - 1) * query.PageSize
	err := db.Preload("Author").
		Offset(offset).
		Limit(query.PageSize).
		Order("created_at DESC").
		Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// IncrementViewCount 增加文章浏览次数
func (r *articleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).
		Error
}

// IncrementLikeCount 增加文章点赞次数
func (r *articleRepository) IncrementLikeCount(id uint) error {
	return r.db.Model(&models.Article{}).
		Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).
		Error
}

// IncrementCommentCount 增加文章评论次数
func (r *articleRepository) IncrementCommentCount(id uint) error {
	return r.db.Model(&models.Article{}).
		Where("id = ?", id).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).
		Error
}
