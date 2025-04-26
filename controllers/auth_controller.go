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

// NewAuthController 创建认证控制器
func NewAuthController(userService services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

// Register 用户注册
// @Summary      用户账号注册
// @Description  新用户注册接口，提供用户名、密码、邮箱等信息进行账号创建
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Param        userData  body      validation.UserCreateDTO  true  "用户注册信息，包含用户名、密码、邮箱等"
// @Success      200       {object}  utils.Response{data=models.UserResponse}  "注册成功，返回用户信息"
// @Failure      400       {object}  utils.Response  "请求参数错误，包括格式不正确或必填字段缺失"
// @Failure      409       {object}  utils.Response  "用户名或邮箱已被占用"
// @Failure      422       {object}  utils.Response  "密码强度不符合要求"
// @Failure      500       {object}  utils.Response  "服务器内部错误"
// @Router       /api/v1/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	// 验证请求参数
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
// @Summary      用户账号登录
// @Description  用户登录接口，通过用户名和密码验证身份并返回JWT访问令牌
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Param        credentials  body      validation.UserLoginDTO  true  "登录凭证，包含用户名和密码"
// @Success      200          {object}  utils.Response{data=models.LoginResponse}  "登录成功，返回访问令牌和用户信息"
// @Failure      400          {object}  utils.Response  "请求参数错误，包括格式不正确或必填字段缺失"
// @Failure      401          {object}  utils.Response  "用户名或密码错误"
// @Failure      403          {object}  utils.Response  "账号已被禁用"
// @Failure      429          {object}  utils.Response  "登录尝试次数过多，请稍后再试"
// @Failure      500          {object}  utils.Response  "服务器内部错误"
// @Router       /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	// 验证请求参数
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
