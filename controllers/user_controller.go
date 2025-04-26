package controllers

import (
	"strconv"

	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/services"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"gitee.com/NextEraAbyss/gin-template/validation"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService services.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// List 获取用户列表
// @Summary      用户列表查询
// @Description  支持分页、排序和关键词搜索的用户列表查询接口
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page        query    int     false  "页码，从1开始计数"              default(1)
// @Param        page_size   query    int     false  "每页记录数，默认10条"            default(10)
// @Param        keyword     query    string  false  "搜索关键词，支持用户名、邮箱和昵称模糊搜索"
// @Param        order_by    query    string  false  "排序字段，支持id、username、created_at等"  default(id)
// @Param        order       query    string  false  "排序方向: asc(升序)或desc(降序)"     default(desc)
// @Success      200         {object}  models.UserListResponse  "用户列表数据，包含总数和分页记录"
// @Failure      400         {object}  utils.Response           "请求参数错误"
// @Failure      401         {object}  utils.Response           "未授权，请先登录"
// @Failure      500         {object}  utils.Response           "服务器内部错误"
// @Router       /api/v1/users [get]
func (ctrl *UserController) List(c *gin.Context) {
	// 验证查询参数
	var queryDTO validation.UserQueryDTO
	if !utils.ValidateQuery(c, &queryDTO) {
		return
	}

	// 使用默认值
	if queryDTO.Page <= 0 {
		queryDTO.Page = 1
	}
	if queryDTO.PageSize <= 0 {
		queryDTO.PageSize = 10
	}
	if queryDTO.OrderBy == "" {
		queryDTO.OrderBy = "id"
	}
	if queryDTO.Order == "" {
		queryDTO.Order = "desc"
	}

	// 转换为模型层查询对象
	query := &models.UserQueryDTO{
		Page:     queryDTO.Page,
		PageSize: queryDTO.PageSize,
		Keyword:  queryDTO.Keyword,
		OrderBy:  queryDTO.OrderBy,
		Order:    queryDTO.Order,
	}

	// 获取用户列表
	users, total, err := ctrl.userService.List(c.Request.Context(), query)
	if err != nil {
		utils.ResponseError(c, utils.CodeInternalError, err.Error())
		return
	}

	// 返回结果
	utils.ResponseSuccess(c, models.UserListResponse{
		Total: total,
		Items: users,
	})
}

// Get 获取用户详情
// @Summary      获取单个用户信息
// @Description  根据用户ID获取用户的详细资料信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id          path     int  true   "用户ID (必填)"
// @Success      200         {object}  models.UserResponse  "用户详细信息"
// @Failure      400         {object}  utils.Response       "用户ID格式错误"
// @Failure      401         {object}  utils.Response       "未授权，请先登录"
// @Failure      404         {object}  utils.Response       "用户不存在"
// @Failure      500         {object}  utils.Response       "服务器内部错误"
// @Router       /api/v1/users/{id} [get]
func (ctrl *UserController) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, utils.CodeInvalidParams, "无效的用户ID")
		return
	}

	user, err := ctrl.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.ResponseError(c, utils.CodeUserNotFound, "用户不存在")
		return
	}

	utils.ResponseSuccess(c, user.ToResponse())
}

// Update 更新用户信息
// @Summary      更新用户信息
// @Description  根据用户ID更新用户资料，支持部分字段更新
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id          path     int           true   "用户ID (必填)"
// @Param        user        body     validation.UserUpdateDTO  true   "用户信息更新内容"
// @Success      200         {object}  models.UserResponse       "更新后的用户信息"
// @Failure      400         {object}  utils.Response            "请求参数错误"
// @Failure      401         {object}  utils.Response            "未授权，请先登录"
// @Failure      403         {object}  utils.Response            "无权限操作此用户"
// @Failure      404         {object}  utils.Response            "用户不存在"
// @Failure      500         {object}  utils.Response            "服务器内部错误"
// @Router       /api/v1/users/{id} [put]
func (ctrl *UserController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, utils.CodeInvalidParams, "无效的用户ID")
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ResponseError(c, utils.CodeInvalidParams, err.Error())
		return
	}

	user.ID = uint(id)
	if err := ctrl.userService.Update(c.Request.Context(), &user); err != nil {
		utils.ResponseError(c, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(c, user.ToResponse())
}

// Delete 删除用户
// @Summary      删除用户
// @Description  根据用户ID删除指定用户记录
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id          path     int        true   "用户ID (必填)"
// @Success      200         {object}  utils.Response   "删除成功"
// @Failure      400         {object}  utils.Response   "用户ID格式错误"
// @Failure      401         {object}  utils.Response   "未授权，请先登录"
// @Failure      403         {object}  utils.Response   "无权删除此用户"
// @Failure      404         {object}  utils.Response   "用户不存在"
// @Failure      500         {object}  utils.Response   "服务器内部错误"
// @Router       /api/v1/users/{id} [delete]
func (ctrl *UserController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, utils.CodeInvalidParams, "无效的用户ID")
		return
	}

	if err := ctrl.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		utils.ResponseError(c, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(c, nil)
}

// ChangePassword 修改密码
// @Summary      修改用户密码
// @Description  修改当前登录用户的密码，需要提供旧密码和新密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        passwordData  body     models.ChangePasswordRequest  true   "密码修改数据，包含旧密码和新密码"
// @Success      200           {object}  utils.Response               "密码修改成功"
// @Failure      400           {object}  utils.Response               "请求参数错误"
// @Failure      401           {object}  utils.Response               "未授权，请先登录"
// @Failure      403           {object}  utils.Response               "旧密码验证失败"
// @Failure      422           {object}  utils.Response               "新密码格式不符合要求"
// @Failure      500           {object}  utils.Response               "服务器内部错误"
// @Router       /api/v1/users/change-password [post]
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.CodeInvalidParams, err.Error())
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ResponseError(c, utils.CodeUnauthorized, "未登录")
		return
	}

	if err := ctrl.userService.ChangePassword(c.Request.Context(), userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		utils.ResponseError(c, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(c, nil)
}
