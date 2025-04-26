package controllers

import (
	"gitee.com/NextEraAbyss/gin-template/internal/errors"
	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/services"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"gitee.com/NextEraAbyss/gin-template/validation"
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
// @Param        register  body      validation.UserCreateDTO  true  "注册信息"
// @Success      200       {object}  utils.Response{data=models.UserResponse}
// @Failure      400       {object}  utils.Response "参数错误"
// @Failure      500       {object}  utils.Response "服务器内部错误"
// @Router       /api/v1/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	// 使用验证工具验证请求参数
	var createDTO validation.UserCreateDTO
	if !utils.ValidateJSON(ctx, &createDTO) {
		return
	}

	// 创建用户
	user := &models.User{
		Username: createDTO.Username,
		Password: createDTO.Password,
		Email:    createDTO.Email,
		Nickname: createDTO.Nickname,
		Status:   createDTO.Status,
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
// @Param        login  body      validation.UserLoginDTO  true  "登录信息"
// @Success      200    {object}  utils.Response{data=models.LoginResponse}
// @Failure      400    {object}  utils.Response "参数错误"
// @Failure      401    {object}  utils.Response "用户名或密码错误"
// @Router       /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	// 使用验证工具验证请求参数
	var loginDTO validation.UserLoginDTO
	if !utils.ValidateJSON(ctx, &loginDTO) {
		return
	}

	token, user, err := c.userService.Login(ctx, loginDTO.Username, loginDTO.Password)
	if err != nil {
		utils.ResponseError(ctx, errors.CodeLoginFailed, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, models.LoginResponse{
		Token: token,
		User:  *user,
	})
}
