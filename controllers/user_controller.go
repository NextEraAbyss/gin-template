package controllers

import (
	"strconv"

	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/services"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
)

// UserController 接口定义了用户相关的处理器
type UserController interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Login(c *gin.Context)
}

// userController 实现 UserController 接口
type userController struct {
	userService services.UserService
}

// NewUserController 创建一个新的 UserController 实例
func NewUserController(userService services.UserService) UserController {
	return &userController{userService: userService}
}

// Create 创建一个新用户
func (ctrl *userController) Create(c *gin.Context) {
	var dto models.UserCreateDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.InvalidParams(c, "无效的请求参数："+err.Error())
		return
	}

	// 转换DTO为User模型
	user := dto.ToUser()

	if err := ctrl.userService.CreateUser(&user); err != nil {
		utils.Failure(c, err)
		return
	}

	utils.Success(c, user)
}

// GetByID 根据ID获取用户
func (ctrl *userController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.InvalidParams(c, "无效的用户ID")
		return
	}

	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

// GetAll 获取所有用户
func (ctrl *userController) GetAll(c *gin.Context) {
	// 绑定查询参数
	var query models.UserQueryDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.InvalidParams(c, "无效的查询参数："+err.Error())
		return
	}

	// 设置默认值
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	users, err := ctrl.userService.GetAllUsers()
	if err != nil {
		utils.Failure(c, err)
		return
	}

	if len(users) == 0 {
		utils.NoContent(c)
		return
	}

	utils.Success(c, users)
}

// Update 更新用户信息
func (ctrl *userController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.InvalidParams(c, "无效的用户ID")
		return
	}

	// 获取现有用户
	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 绑定更新DTO
	var dto models.UserUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.InvalidParams(c, "无效的请求参数："+err.Error())
		return
	}

	// 更新用户模型
	dto.UpdateUser(user)

	if err := ctrl.userService.UpdateUser(user); err != nil {
		utils.Failure(c, err)
		return
	}

	utils.SuccessWithMessage(c, "用户信息更新成功", user)
}

// Delete 删除用户
func (ctrl *userController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.InvalidParams(c, "无效的用户ID")
		return
	}

	if err := ctrl.userService.DeleteUser(uint(id)); err != nil {
		utils.Failure(c, err)
		return
	}

	utils.SuccessWithMessage(c, "用户删除成功", nil)
}

// Login 用户登录
func (ctrl *userController) Login(c *gin.Context) {
	var dto models.LoginDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.InvalidParams(c, "无效的登录信息："+err.Error())
		return
	}

	user, err := ctrl.userService.Authenticate(dto.Username, dto.Password)
	if err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.Failure(c, err)
		return
	}

	utils.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}
