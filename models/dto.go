package models

// UserCreateDTO 用户创建DTO
type UserCreateDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8,max=100"`
	FullName string `json:"full_name" binding:"max=100"`
}

// UserUpdateDTO 用户更新DTO
type UserUpdateDTO struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Password string `json:"password" binding:"omitempty,min=8,max=100"`
	FullName string `json:"full_name" binding:"omitempty,max=100"`
}

// LoginDTO 登录DTO
type LoginDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=1,max=100"`
}

// UserQueryDTO 用户查询DTO
type UserQueryDTO struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Sort     string `form:"sort" binding:"omitempty"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
	Search   string `form:"search" binding:"omitempty"`
}

// ToUser 转换为User模型
func (dto *UserCreateDTO) ToUser() User {
	return User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		FullName: dto.FullName,
	}
}

// UpdateUser 将DTO数据更新到用户模型中
func (dto *UserUpdateDTO) UpdateUser(user *User) {
	if dto.Username != "" {
		user.Username = dto.Username
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.Password != "" {
		user.Password = dto.Password
	}
	if dto.FullName != "" {
		user.FullName = dto.FullName
	}
}
