package controllers

import (
	"strconv"

	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/services"
	"gitee.com/NextEraAbyss/gin-template/utils"
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
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码，默认为1" default(1)
// @Param page_size query int false "每页数量，默认为10" default(10)
// @Success 200 {object} models.UserListResponse "用户列表获取成功"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users [get]
func (ctrl *UserController) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 创建查询参数对象
	query := &models.UserQueryDTO{
		Page:     page,
		PageSize: pageSize,
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
// @Summary 获取用户详情
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} models.UserResponse "用户详情获取成功"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "用户不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/{id} [get]
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
// @Summary 更新用户信息
// @Description 更新指定用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Param user body models.User true "用户信息"
// @Success 200 {object} models.UserResponse "用户信息更新成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "用户不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/{id} [put]
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
// @Summary 删除用户
// @Description 删除指定用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response "用户删除成功"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 404 {object} utils.Response "用户不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/{id} [delete]
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
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param passwordForm body models.ChangePasswordRequest true "密码信息"
// @Success 200 {object} utils.Response "密码修改成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 401 {object} utils.Response "未授权"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/change-password [post]
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
