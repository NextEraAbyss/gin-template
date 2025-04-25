package controllers

import (
	"strconv"

	"gitee.com/NextEraAbyss/gin-template/internal/errors"
	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/services"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

const (
	// 错误消息常量
	ErrArticleNotFound  = "文章不存在"
	ErrInvalidArticleID = "无效的文章ID"
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
// @Summary      获取文章列表
// @Description  分页获取文章列表，支持搜索和排序
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Param        page       query    int     false  "页码"  default(1)
// @Param        page_size  query    int     false  "每页数量"  default(10)
// @Param        keyword    query    string  false  "搜索关键词"
// @Param        status     query    int     false  "文章状态(1:草稿,2:已发布)"
// @Param        order_by   query    string  false  "排序字段(created_at,view_count)"
// @Param        order      query    string  false  "排序方向(asc,desc)"
// @Success      200  {object}  utils.Response{data=models.ArticleListResponse}
// @Router       /api/v1/articles [get]
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

	// 转换为响应格式
	responses := make([]*models.ArticleResponse, len(articles))
	for i := range articles {
		responses[i] = articles[i].ToResponse()
	}

	utils.ResponseSuccess(ctx, &models.ArticleListResponse{
		Total: total,
		Items: responses,
	})
}

// Get 获取文章详情
// @Summary      获取文章详情
// @Description  根据ID获取文章详情
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "文章ID"
// @Success      200  {object}  utils.Response{data=models.ArticleResponse}
// @Router       /api/v1/articles/{id} [get]
func (c *ArticleController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, ErrInvalidArticleID)
		return
	}

	article, err := c.service.GetArticleByID(uint(id))
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	// 增加浏览次数
	if err := c.service.IncrementViewCount(uint(id)); err != nil {
		utils.Errorf("Failed to increment view count: %v", err)
	}

	utils.ResponseSuccess(ctx, article.ToResponse())
}

// Create 创建文章
// @Summary      创建文章
// @Description  创建新文章
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        article  body      models.CreateArticleRequest  true  "文章信息"
// @Success      200      {object}  utils.Response{data=models.ArticleResponse}
// @Failure      400      {object}  utils.Response "参数错误"
// @Failure      401      {object}  utils.Response "未授权"
// @Failure      500      {object}  utils.Response "服务器内部错误"
// @Router       /api/v1/articles [post]
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

	utils.ResponseSuccess(ctx, article.ToResponse())
}

// Update 更新文章
// @Summary      更新文章
// @Description  更新文章信息
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id       path      int                        true  "文章ID"
// @Param        article  body      models.UpdateArticleRequest  true  "文章信息"
// @Success      200      {object}  utils.Response{data=models.ArticleResponse}
// @Failure      400      {object}  utils.Response "参数错误"
// @Failure      401      {object}  utils.Response "未授权"
// @Failure      404      {object}  utils.Response "文章不存在"
// @Failure      500      {object}  utils.Response "服务器内部错误"
// @Router       /api/v1/articles/{id} [put]
func (c *ArticleController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, errors.CodeInvalidParams, "无效的文章ID")
		return
	}

	var req models.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, errors.CodeInvalidParams, err.Error())
		return
	}

	article := &models.Article{
		ID:      uint(id),
		Title:   req.Title,
		Content: req.Content,
		Status:  *req.Status,
	}

	if err := c.service.UpdateArticle(article); err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, nil)
}

// Delete 删除文章
// @Summary      删除文章
// @Description  删除指定文章
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path  int  true  "文章ID"
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response "未授权"
// @Failure      404  {object}  utils.Response "文章不存在"
// @Failure      500  {object}  utils.Response "服务器内部错误"
// @Router       /api/v1/articles/{id} [delete]
func (c *ArticleController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, ErrInvalidArticleID)
		return
	}

	if err := c.service.DeleteArticle(uint(id)); err != nil {
		if err.Error() == ErrArticleNotFound {
			utils.ResponseError(ctx, utils.CodeArticleNotFound, err.Error())
			return
		}
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, nil)
}

// Like 点赞文章
func (c *ArticleController) Like(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, errors.CodeInvalidParams, "无效的文章ID")
		return
	}

	err = c.service.IncrementLikeCount(uint(id))
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, nil)
}

// View 查看文章
func (c *ArticleController) View(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, errors.CodeInvalidParams, "无效的文章ID")
		return
	}

	err = c.service.IncrementViewCount(uint(id))
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, nil)
}

// Comment 评论文章
func (c *ArticleController) Comment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(ctx, errors.CodeInvalidParams, "无效的文章ID")
		return
	}

	err = c.service.IncrementCommentCount(uint(id))
	if err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, nil)
}
