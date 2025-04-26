package repositories

import (
	"context"

	"gitee.com/NextEraAbyss/gin-template/models"

	"gorm.io/gorm"
)

// UserQueryParams 用户查询参数
type UserQueryParams struct {
	Page     int    `json:"page"`     // 页码
	PageSize int    `json:"pageSize"` // 每页数量
	Keyword  string `json:"keyword"`  // 搜索关键词
	Status   int    `json:"status"`   // 状态
	OrderBy  string `json:"orderBy"`  // 排序字段
	Order    string `json:"order"`    // 排序方向
}

// BaseRepository 基础仓库接口，定义通用操作
type BaseRepository interface {
	// RepositoryName 获取仓库名称
	RepositoryName() string
}

// UserRepository 用户仓库接口
type UserRepository interface {
	BaseRepository

	// Create 创建用户
	Create(ctx context.Context, user *models.User) error

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id uint) (*models.User, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *models.User) error

	// Delete 删除用户
	Delete(ctx context.Context, id uint) error

	// List 获取用户列表
	List(ctx context.Context, query *UserQueryParams) ([]*models.User, int64, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

// userRepository 实现 UserRepository 接口
type userRepository struct {
	db *gorm.DB
}

// RepositoryName 获取仓库名称
func (r *userRepository) RepositoryName() string {
	return "UserRepository"
}

// NewUserRepository 创建一个新的 UserRepository 实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, query *UserQueryParams) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64
	var order string

	// 构建查询
	db := r.db.WithContext(ctx).Model(&models.User{})

	// 添加查询条件
	if query.Keyword != "" {
		db = db.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+query.Keyword+"%",
			"%"+query.Keyword+"%",
			"%"+query.Keyword+"%")
	}

	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}

	// 排序
	if query.OrderBy != "" {
		order = "desc"
		if query.Order == "asc" {
			order = "asc"
		}
		db = db.Order(query.OrderBy + " " + order)
	} else {
		db = db.Order("created_at DESC")
	}

	// 执行查询
	if err := db.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
