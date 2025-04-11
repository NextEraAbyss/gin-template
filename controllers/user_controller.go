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

// Create 创建用户
func (ctrl *UserController) Create(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ResponseError(c, utils.CodeInvalidParams, err.Error())
		return
	}

	if err := ctrl.userService.Create(c.Request.Context(), &user); err != nil {
		utils.ResponseError(c, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(c, user.ToResponse())
}

// List 获取用户列表
func (ctrl *UserController) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	result, err := ctrl.userService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		utils.ResponseError(c, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(c, result)
}

// Get 获取用户详情
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

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	var loginForm struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginForm); err != nil {
		utils.ResponseError(c, utils.CodeInvalidParams, err.Error())
		return
	}

	token, err := ctrl.userService.Login(c.Request.Context(), loginForm.Username, loginForm.Password)
	if err != nil {
		utils.ResponseError(c, utils.CodeUnauthorized, "用户名或密码错误")
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"token": token,
	})
}

// ChangePassword 修改密码
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

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

// ResetPassword 重置密码
func (ctrl *UserController) ResetPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.CodeInvalidParams, err.Error())
		return
	}

	if err := ctrl.userService.ResetPassword(c.Request.Context(), req.Email); err != nil {
		utils.ResponseError(c, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(c, nil)
}
