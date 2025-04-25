package controllers

import (
	"gitee.com/NextEraAbyss/gin-template/internal/errors"
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
// @Summary      用户注册
// @Description  新用户注册接口
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Param        register  body      models.RegisterRequest  true  "注册信息"
// @Success      200       {object}  utils.Response{data=models.UserResponse}
// @Failure      400       {object}  utils.Response "参数错误"
// @Failure      500       {object}  utils.Response "服务器内部错误"
// @Router       /api/v1/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, errors.CodeInvalidParams, err.Error())
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
		utils.ResponseError(ctx, errors.CodeRegisterFailed, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, user.ToResponse())
}

// Login 用户登录
// @Summary      用户登录
// @Description  用户登录接口，返回JWT令牌
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Param        login  body      models.LoginRequest  true  "登录信息"
// @Success      200    {object}  utils.Response{data=models.LoginResponse}
// @Failure      400    {object}  utils.Response "参数错误"
// @Failure      401    {object}  utils.Response "用户名或密码错误"
// @Router       /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(ctx, errors.CodeInvalidParams, err.Error())
		return
	}

	token, user, err := c.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		utils.ResponseError(ctx, errors.CodeLoginFailed, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, models.LoginResponse{
		Token: token,
		User:  *user,
	})
}
