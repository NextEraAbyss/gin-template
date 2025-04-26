package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/repositories"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	// Create 创建新用户
	Create(ctx context.Context, user *models.User) error

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id uint) (*models.User, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *models.User) error

	// Delete 删除用户
	Delete(ctx context.Context, id uint) error

	// List 获取用户列表
	List(ctx context.Context, query *models.UserQueryDTO) ([]models.User, int64, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error

	// ResetPassword 重置密码
	ResetPassword(ctx context.Context, email string) error

	// Login 用户登录
	Login(ctx context.Context, username, password string) (string, *models.User, error)

	// Register 用户注册
	Register(ctx context.Context, user *models.User) error
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

	// 设置最后登录时间为当前时间
	user.LastLoginAt = time.Now()

	return s.repo.Create(ctx, user)
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// Update 更新用户信息
func (s *userService) Update(ctx context.Context, user *models.User) error {
	// 检查用户是否存在
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return fmt.Errorf("用户不存在")
	}

	// 如果更新了密码，需要重新加密
	if user.Password != "" && user.Password != existingUser.Password {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}

		user.Password = hashedPassword
	}

	// 更新用户信息
	return s.repo.Update(ctx, user)
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id uint) error {
	// 先检查用户是否存在
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

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
func (s *userService) List(ctx context.Context, query *models.UserQueryDTO) ([]models.User, int64, error) {
	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}

	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	users, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// 转换为非指针切片
	items := make([]models.User, 0, len(users))

	for _, user := range users {
		if user != nil {
			items = append(items, *user)
		}
	}

	return items, total, nil
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 生成随机密码
	newPassword, err := utils.GenerateRandomPassword()
	if err != nil {
		return err
	}

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
func (s *userService) Login(ctx context.Context, username, password string) (string, *models.User, error) {
	// 获取用户信息
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("用户不存在")
		}
		return "", nil, err
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return "", nil, errors.New("密码错误")
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now()
	if err := s.repo.Update(ctx, user); err != nil {
		return "", nil, err
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// Register 注册用户
func (s *userService) Register(ctx context.Context, user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.Create(ctx, user)
}
