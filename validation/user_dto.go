package validation

import "gitee.com/NextEraAbyss/gin-template/models"

// UserQueryDTO 用户查询参数
type UserQueryDTO struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword" binding:"omitempty,max=50"`
	OrderBy  string `form:"order_by" binding:"omitempty,oneof=id username email created_at updated_at status"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// UserCreateDTO 创建用户参数
type UserCreateDTO struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"omitempty,max=32"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// UserUpdateDTO 更新用户参数
type UserUpdateDTO struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"omitempty,min=3,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname" binding:"omitempty,max=32"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// UserLoginDTO 用户登录参数
type UserLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserChangePasswordDTO 修改密码参数
type UserChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32,nefield=OldPassword"`
}

// UserResetPasswordDTO 重置密码参数
type UserResetPasswordDTO struct {
	Email string `json:"email" binding:"required,email"`
}

// UserResponseDTO 用户信息响应
type UserResponseDTO struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	LastLoginAt string `json:"last_login_at"`
}

// LoginResponseDTO 登录响应
type LoginResponseDTO struct {
	Token string          `json:"token"`
	User  UserResponseDTO `json:"user"`
}

// GetDefaultUserQuery 获取默认的用户查询参数
func GetDefaultUserQuery() UserQueryDTO {
	return UserQueryDTO{
		Page:     1,
		PageSize: 10,
		OrderBy:  "id",
		Order:    "desc",
	}
}

// ToModel 将UserCreateDTO转换为User模型
func (dto *UserCreateDTO) ToModel() models.User {
	return models.User{
		Username: dto.Username,
		Password: dto.Password,
		Email:    dto.Email,
		Nickname: dto.Nickname,
		Status:   dto.Status,
	}
}

// UpdateModel 将UserUpdateDTO的数据更新到User模型中
func (dto *UserUpdateDTO) UpdateModel(user *models.User) {
	if dto.Username != "" {
		user.Username = dto.Username
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.Nickname != "" {
		user.Nickname = dto.Nickname
	}
	if dto.Status != 0 {
		user.Status = dto.Status
	}
}

// FromUser 从User模型创建UserResponseDTO
func FromUser(user models.User) UserResponseDTO {
	return UserResponseDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   user.UpdatedAt.Format("2006-01-02 15:04:05"),
		LastLoginAt: user.LastLoginAt.Format("2006-01-02 15:04:05"),
	}
}
