package controllers

import (
	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/services"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	userService services.UserService
}

// NewAuthController 创建新的认证控制器
func NewAuthController(userService services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, err.Error())
		return
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Nickname: req.FullName, // 使用 FullName 作为昵称
	}

	if err := c.userService.Create(ctx, user); err != nil {
		utils.ResponseError(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, user.ToResponse())
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, utils.CodeInvalidParams, err.Error())
		return
	}

	token, err := c.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		utils.ResponseError(ctx, utils.CodeUnauthorized, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, gin.H{
		"token": token,
	})
}
