package repositories

import (
	"context"

	"gitee.com/NextEraAbyss/gin-template/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
// 提供用户数据的持久化操作，包括增删改查和特定条件查询
type UserRepository interface {
	// Create 创建用户
	// 参数：
	// - ctx: 上下文，用于控制数据库操作的超时和取消
	// - user: 用户实体，包含需要存储的用户数据
	// 返回：
	// - error: 数据库操作错误，如唯一索引冲突等
	Create(ctx context.Context, user *models.User) error

	// GetByID 根据ID获取用户
	// 参数：
	// - ctx: 上下文，用于控制数据库操作的超时和取消
	// - id: 用户ID，数据库主键
	// 返回：
	// - *models.User: 查询到的用户实体，若不存在则为nil
	// - error: 数据库操作错误，如记录不存在等
	GetByID(ctx context.Context, id uint) (*models.User, error)

	// Update 更新用户信息
	// 参数：
	// - ctx: 上下文，用于控制数据库操作的超时和取消
	// - user: 用户实体，包含需要更新的字段
	// 返回：
	// - error: 数据库操作错误，如记录不存在、字段验证失败等
	Update(ctx context.Context, user *models.User) error

	// Delete 删除用户
	// 参数：
	// - ctx: 上下文，用于控制数据库操作的超时和取消
	// - id: 用户ID，数据库主键
	// 返回：
	// - error: 数据库操作错误，如记录不存在等
	Delete(ctx context.Context, id uint) error

	// List 获取用户列表
	// 参数：
	// - ctx: 上下文，用于控制数据库操作的超时和取消
	// - query: 查询参数，包含分页、排序和筛选条件
	// 返回：
	// - []*models.User: 用户列表
	// - int64: 符合条件的总记录数（不考虑分页）
	// - error: 数据库操作错误
	List(ctx context.Context, query *models.UserQueryDTO) ([]*models.User, int64, error)

	// GetByUsername 根据用户名获取用户
	// 参数：
	// - ctx: 上下文，用于控制数据库操作的超时和取消
	// - username: 用户名，唯一标识
	// 返回：
	// - *models.User: 查询到的用户实体，若不存在则为nil
	// - error: 数据库操作错误，如记录不存在等
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// GetByEmail 根据邮箱获取用户
	// 参数：
	// - ctx: 上下文，用于控制数据库操作的超时和取消
	// - email: 邮箱地址，唯一标识
	// 返回：
	// - *models.User: 查询到的用户实体，若不存在则为nil
	// - error: 数据库操作错误，如记录不存在等
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

// userRepository 实现 UserRepository 接口
type userRepository struct {
	db *gorm.DB
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
func (r *userRepository) List(ctx context.Context, query *models.UserQueryDTO) ([]*models.User, int64, error) {
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
