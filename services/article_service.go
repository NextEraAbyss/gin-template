package services

import (
	"errors"
	"fmt"

	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/repositories"
	"gorm.io/gorm"
)

// ArticleService 文章服务接口.
type ArticleService interface {
	GetArticleList(req *models.ArticleListRequest) ([]*models.Article, int64, error) // 获取文章列表.
	GetArticleByID(id uint) (*models.Article, error)                                 // 获取单个文章.
	CreateArticle(article *models.Article) error                                     // 创建文章.
	UpdateArticle(article *models.Article) error                                     // 更新文章.
	DeleteArticle(id uint) error                                                     // 删除文章.
	IncrementViewCount(id uint) error                                                // 增加浏览次数.
	IncrementLikeCount(id uint) error                                                // 增加点赞次数.
	IncrementCommentCount(id uint) error                                             // 增加评论次数.
}

// articleService 文章服务实现
type articleService struct {
	repo repositories.ArticleRepository
}

// NewArticleService 创建文章服务实例
func NewArticleService(repo repositories.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

// GetArticleList 获取文章列表
func (s *articleService) GetArticleList(req *models.ArticleListRequest) ([]*models.Article, int64, error) {
	query := &models.ArticleQueryDTO{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Status:   req.Status,
	}
	return s.repo.List(query)
}

// GetArticleByID 获取单个文章
func (s *articleService) GetArticleByID(id uint) (*models.Article, error) {
	return s.repo.FindByID(id)
}

// CreateArticle 创建文章
func (s *articleService) CreateArticle(article *models.Article) error {
	return s.repo.Create(article)
}

// UpdateArticle 更新文章
func (s *articleService) UpdateArticle(article *models.Article) error {
	// 检查文章是否存在
	existingArticle, err := s.repo.FindByID(article.ID)
	if err != nil {
		return err
	}

	if existingArticle == nil {
		return fmt.Errorf("文章不存在")
	}

	return s.repo.Update(article)
}

// DeleteArticle 删除文章
func (s *articleService) DeleteArticle(id uint) error {
	// 检查文章是否存在
	existingArticle, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if existingArticle == nil {
		return fmt.Errorf("文章不存在")
	}

	return s.repo.Delete(id)
}

// IncrementViewCount 增加文章浏览次数
func (s *articleService) IncrementViewCount(id uint) error {
	// 检查文章是否存在
	existingArticle, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if existingArticle == nil {
		return fmt.Errorf("文章不存在")
	}

	return s.repo.IncrementViewCount(id)
}

// IncrementLikeCount 增加文章点赞次数
func (s *articleService) IncrementLikeCount(id uint) error {
	// 先检查文章是否存在
	_, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("文章不存在")
		}
		return err
	}

	return s.repo.IncrementLikeCount(id)
}

// IncrementCommentCount 增加文章评论次数
func (s *articleService) IncrementCommentCount(id uint) error {
	// 先检查文章是否存在
	_, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("文章不存在")
		}
		return err
	}

	return s.repo.IncrementCommentCount(id)
}
