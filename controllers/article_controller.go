package controllers

import (
	"strconv"

	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/services"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// ArticleController 文章控制器
// 处理文章相关的HTTP请求
type ArticleController struct {
	service services.ArticleService
}

// NewArticleController 创建文章控制器实例
func NewArticleController(service services.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

// List 获取文章列表
// @Summary 获取文章列表
// @Description 分页获取文章列表，支持关键词搜索、状态过滤和排序
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param keyword query string false "搜索关键词"
// @Param status query int false "文章状态(1:草稿,2:已发布)"
// @Param order_by query string false "排序字段(created_at,view_count)"
// @Param order query string false "排序方向(asc,desc)"
// @Success 200 {object} models.ArticleListResponse
// @Router /api/articles [get]
func (c *ArticleController) List(ctx *gin.Context) {
	var req models.ArticleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	articles, total, err := c.service.GetArticleList(&req)
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, gin.H{
		"total": total,
		"items": articles,
	})
}

// Get 获取单个文章
// @Summary 获取单个文章
// @Description 根据ID获取文章详情
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} models.Article
// @Router /api/articles/{id} [get]
func (c *ArticleController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, "无效的文章ID")
		return
	}

	article, err := c.service.GetArticleByID(uint(id))
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	// 增加浏览次数
	go c.service.IncrementViewCount(uint(id))

	utils.ResponseSuccess(ctx, article)
}

// Create 创建文章
// @Summary 创建文章
// @Description 创建新文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param article body models.CreateArticleRequest true "文章信息"
// @Success 201 {object} models.Article
// @Router /api/articles [post]
func (c *ArticleController) Create(ctx *gin.Context) {
	var req models.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, err.Error())
		return
	}

	// 获取当前用户ID
	authorID := ctx.GetUint("user_id")
	if authorID == 0 {
		utils.ResponseError(ctx, utils.CodeUnauthorized, "未登录")
		return
	}

	article := &models.Article{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
		Status:   req.Status,
	}

	if err := c.service.CreateArticle(article); err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, article)
}

// Update 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param article body models.UpdateArticleRequest true "文章信息"
// @Success 200 {object} models.Article
// @Router /api/articles/{id} [put]
func (c *ArticleController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, "无效的文章ID")
		return
	}

	var req models.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, err.Error())
		return
	}

	article := &models.Article{
		Title:   req.Title,
		Content: req.Content,
		Status:  *req.Status,
	}

	if err := c.service.UpdateArticle(uint(id), article); err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, nil)
}

// Delete 删除文章
// @Summary 删除文章
// @Description 删除指定文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} string
// @Router /api/articles/{id} [delete]
func (c *ArticleController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, "无效的文章ID")
		return
	}

	if err := c.service.DeleteArticle(uint(id)); err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, nil)
}

// Like 点赞文章
func (c *ArticleController) Like(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, "无效的文章ID")
		return
	}

	if err := c.service.IncrementLikeCount(uint(id)); err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, nil)
}
