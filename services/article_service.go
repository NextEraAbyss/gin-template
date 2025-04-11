package services

import (
	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/repositories"
)

// ArticleService 文章服务接口
// 定义了文章相关的业务逻辑操作
type ArticleService interface {
	GetArticleList(req *models.ArticleListRequest) ([]models.Article, int64, error) // 获取文章列表
	GetArticleByID(id uint) (*models.Article, error)                                // 获取单个文章
	CreateArticle(article *models.Article) error                                    // 创建文章
	UpdateArticle(id uint, article *models.Article) error                           // 更新文章
	DeleteArticle(id uint) error                                                    // 删除文章
	IncrementViewCount(id uint) error                                               // 增加文章浏览次数
	IncrementLikeCount(id uint) error                                               // 增加文章点赞次数
	IncrementCommentCount(id uint) error                                            // 增加文章评论次数
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
func (s *articleService) GetArticleList(req *models.ArticleListRequest) ([]models.Article, int64, error) {
	query := &models.ArticleQueryDTO{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		AuthorID: req.AuthorID,
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
func (s *articleService) UpdateArticle(id uint, article *models.Article) error {
	// 先获取现有文章
	existingArticle, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 更新文章信息
	existingArticle.Title = article.Title
	existingArticle.Content = article.Content
	existingArticle.Status = article.Status

	return s.repo.Update(existingArticle)
}

// DeleteArticle 删除文章
func (s *articleService) DeleteArticle(id uint) error {
	return s.repo.Delete(id)
}

// IncrementViewCount 增加文章浏览次数
func (s *articleService) IncrementViewCount(id uint) error {
	return s.repo.IncrementViewCount(id)
}

// IncrementLikeCount 增加文章点赞次数
func (s *articleService) IncrementLikeCount(id uint) error {
	return s.repo.IncrementLikeCount(id)
}

// IncrementCommentCount 增加文章评论次数
func (s *articleService) IncrementCommentCount(id uint) error {
	return s.repo.IncrementCommentCount(id)
}
