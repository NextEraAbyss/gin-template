package services

import (
	"context"
	"errors"

	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/repositories"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int) (*models.UserListResponse, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, email string) error
	Login(ctx context.Context, username, password string) (string, error)
}

// userService 用户服务实现
type userService struct {
	repo repositories.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, user *models.User) error {
	// 检查用户名是否已存在
	existingUser, err := s.repo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.repo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.New("邮箱已存在")
	}

	// 验证密码强度
	if err := utils.ValidatePasswordStrength(user.Password); err != nil {
		return err
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// 设置默认状态
	if user.Status == 0 {
		user.Status = 1
	}

	return s.repo.Create(ctx, user)
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// Update 更新用户信息
func (s *userService) Update(ctx context.Context, user *models.User) error {
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 如果修改了用户名，检查新用户名是否已存在
	if user.Username != "" && user.Username != existingUser.Username {
		userWithSameUsername, err := s.repo.GetByUsername(ctx, user.Username)
		if err == nil && userWithSameUsername != nil && userWithSameUsername.ID != user.ID {
			return errors.New("用户名已存在")
		}
	}

	// 如果修改了邮箱，检查新邮箱是否已存在
	if user.Email != "" && user.Email != existingUser.Email {
		userWithSameEmail, err := s.repo.GetByEmail(ctx, user.Email)
		if err == nil && userWithSameEmail != nil && userWithSameEmail.ID != user.ID {
			return errors.New("邮箱已存在")
		}
	}

	// 如果修改了密码，验证密码强度并加密
	if user.Password != "" {
		if err := utils.ValidatePasswordStrength(user.Password); err != nil {
			return err
		}
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	return s.repo.Update(ctx, user)
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// GetByUsername 根据用户名获取用户
func (s *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

// GetByEmail 根据邮箱获取用户
func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// List 获取用户列表
func (s *userService) List(ctx context.Context, page, pageSize int) (*models.UserListResponse, error) {
	users, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 将 []*models.User 转换为 []models.User
	items := make([]models.User, len(users))
	for i, user := range users {
		items[i] = *user
	}

	return &models.UserListResponse{
		Total: total,
		Items: items,
	}, nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	// 获取用户信息
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("旧密码错误")
	}

	// 验证新密码强度
	if err := utils.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = hashedPassword
	return s.repo.Update(ctx, user)
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(ctx context.Context, email string) error {
	// 获取用户信息
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	// 生成随机密码
	newPassword := utils.GenerateRandomPassword()
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	// TODO: 发送邮件通知用户新密码
	return nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	// 获取用户信息
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("用户名或密码错误")
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}
