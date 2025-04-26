package validation

// UserQueryDTO 用户查询参数DTO
type UserQueryDTO struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword" binding:"omitempty,max=50"`
	OrderBy  string `form:"order_by" binding:"omitempty,oneof=id username email created_at updated_at status"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// UserCreateDTO 创建用户参数DTO
type UserCreateDTO struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"omitempty,max=32"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// UserUpdateDTO 更新用户参数DTO
type UserUpdateDTO struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"omitempty,min=3,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname" binding:"omitempty,max=32"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// UserLoginDTO 用户登录参数DTO
type UserLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserChangePasswordDTO 修改密码参数DTO
type UserChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32,nefield=OldPassword"`
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
