package models

// UserCreateDTO 用户创建DTO
type UserCreateDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8,max=100"`
	Nickname string `json:"nickname" binding:"max=100"`
}

// UserUpdateDTO 用户更新DTO
type UserUpdateDTO struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Password string `json:"password" binding:"omitempty,min=8,max=100"`
	Nickname string `json:"nickname" binding:"omitempty,max=100"`
}

// LoginDTO 登录DTO
type LoginDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=1,max=100"`
}

// ToUser 转换为User模型
func (dto *UserCreateDTO) ToUser() User {
	return User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		Nickname: dto.Nickname,
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
	if dto.Nickname != "" {
		user.Nickname = dto.Nickname
	}
}
