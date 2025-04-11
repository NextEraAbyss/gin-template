package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gitee.com/NextEraAbyss/gin-template/internal/redis"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Condition 查询条件
type Condition struct {
	Query interface{}
	Args  []interface{}
}

// QueryStats 查询统计
type QueryStats struct {
	SQL       string
	Duration  time.Duration
	Rows      int64
	Timestamp time.Time
	IsSlow    bool
	Error     error
}

// QueryBuilder 查询构建器
type QueryBuilder struct {
	db      *gorm.DB
	query   *gorm.DB
	context context.Context
	cache   redis.Cache
	// 慢查询阈值
	slowQueryThreshold time.Duration
	// 是否启用缓存
	enableCache bool
	// 缓存过期时间
	cacheExpiration time.Duration
	// 是否忽略软删除
	unscoped bool
	// 是否启用乐观锁
	optimisticLock bool
	// 查询条件
	conditions []Condition
	// 选择的字段
	selectedFields []string
	// 查询统计
	stats *QueryStats
	// 最大结果集大小
	maxResultSize int
	// 是否启用查询计划分析
	enableExplain bool
	// 连接池配置
	poolConfig *PoolConfig
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{
		db:                 db,
		query:              db,
		context:            context.Background(),
		slowQueryThreshold: 100 * time.Millisecond,
		enableCache:        false,
		cacheExpiration:    time.Hour,
		unscoped:           false,
		optimisticLock:     false,
		conditions:         make([]Condition, 0),
		selectedFields:     make([]string, 0),
		maxResultSize:      1000,
		enableExplain:      false,
		poolConfig: &PoolConfig{
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: time.Minute * 30,
		},
	}
}

// WithContext 设置上下文
func (qb *QueryBuilder) WithContext(ctx context.Context) *QueryBuilder {
	qb.context = ctx
	qb.query = qb.query.WithContext(ctx)
	return qb
}

// Where 添加查询条件
func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
	qb.conditions = append(qb.conditions, Condition{
		Query: query,
		Args:  args,
	})
	qb.query = qb.query.Where(query, args...)
	return qb
}

// Or 添加OR条件
func (qb *QueryBuilder) Or(query interface{}, args ...interface{}) *QueryBuilder {
	qb.query = qb.query.Or(query, args...)
	return qb
}

// Order 添加排序
func (qb *QueryBuilder) Order(value interface{}) *QueryBuilder {
	qb.query = qb.query.Order(value)
	return qb
}

// Limit 设置限制
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.query = qb.query.Limit(limit)
	return qb
}

// Offset 设置偏移
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.query = qb.query.Offset(offset)
	return qb
}

// Preload 预加载关联
func (qb *QueryBuilder) Preload(query string, args ...interface{}) *QueryBuilder {
	qb.query = qb.query.Preload(query, args...)
	return qb
}

// Select 选择字段
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.selectedFields = append(qb.selectedFields, fields...)
	qb.query = qb.query.Select(fields)
	return qb
}

// Group 分组
func (qb *QueryBuilder) Group(name string) *QueryBuilder {
	qb.query = qb.query.Group(name)
	return qb
}

// Having 分组条件
func (qb *QueryBuilder) Having(query interface{}, args ...interface{}) *QueryBuilder {
	qb.query = qb.query.Having(query, args...)
	return qb
}

// Joins 连接
func (qb *QueryBuilder) Joins(query string, args ...interface{}) *QueryBuilder {
	qb.query = qb.query.Joins(query, args...)
	return qb
}

// WithCache 启用缓存
func (qb *QueryBuilder) WithCache(cache redis.Cache, expiration time.Duration) *QueryBuilder {
	qb.cache = cache
	qb.enableCache = true
	qb.cacheExpiration = expiration
	return qb
}

// WithSlowQueryThreshold 设置慢查询阈值
func (qb *QueryBuilder) WithSlowQueryThreshold(threshold time.Duration) *QueryBuilder {
	qb.slowQueryThreshold = threshold
	return qb
}

// WithMaxResultSize 设置最大结果集大小
func (qb *QueryBuilder) WithMaxResultSize(size int) *QueryBuilder {
	qb.maxResultSize = size
	return qb
}

// WithExplain 启用查询计划分析
func (qb *QueryBuilder) WithExplain() *QueryBuilder {
	qb.enableExplain = true
	return qb
}

// WithPoolConfig 配置连接池
func (qb *QueryBuilder) WithPoolConfig(config *PoolConfig) *QueryBuilder {
	qb.poolConfig = config
	sqlDB, err := qb.db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
		sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}
	return qb
}

// getCacheKey 生成缓存键
func (qb *QueryBuilder) getCacheKey(operation string, args ...interface{}) string {
	return fmt.Sprintf("query:%s:%v", operation, args)
}

// First 获取第一条记录
func (qb *QueryBuilder) First(dest interface{}) error {
	start := time.Now()
	qb.stats = &QueryStats{}

	// 如果启用查询计划分析，先分析查询计划
	if qb.enableExplain {
		explain, err := qb.Explain()
		if err != nil {
			utils.Warn("查询计划分析失败: %v", err)
		} else {
			utils.Debug("查询计划: %s", explain)
		}
	}

	// 如果启用缓存，尝试从缓存获取
	if qb.enableCache && qb.cache != nil {
		cacheKey := qb.getCacheKey("First", qb.query.Statement.SQL.String())
		if err := qb.cache.Get(qb.context, cacheKey, dest); err == nil {
			return nil
		}
	}

	err := qb.query.First(dest).Error
	qb.stats.Duration = time.Since(start)
	qb.stats.Rows = qb.query.RowsAffected

	// 记录查询统计
	qb.logQueryStats("First", err)

	if err != nil {
		return err
	}

	// 如果启用缓存，将结果存入缓存
	if qb.enableCache && qb.cache != nil {
		cacheKey := qb.getCacheKey("First", qb.query.Statement.SQL.String())
		qb.cache.Set(qb.context, cacheKey, dest, qb.cacheExpiration)
	}

	return nil
}

// Find 获取多条记录
func (qb *QueryBuilder) Find(dest interface{}) error {
	start := time.Now()
	qb.stats = &QueryStats{}

	// 如果启用查询计划分析，先分析查询计划
	if qb.enableExplain {
		explain, err := qb.Explain()
		if err != nil {
			utils.Warn("查询计划分析失败: %v", err)
		} else {
			utils.Debug("查询计划: %s", explain)
		}
	}

	// 限制结果集大小
	if qb.maxResultSize > 0 {
		qb.query = qb.query.Limit(qb.maxResultSize)
	}

	// 执行查询
	err := qb.query.Find(dest).Error
	qb.stats.Duration = time.Since(start)
	qb.stats.Rows = qb.query.RowsAffected

	// 记录查询统计
	qb.logQueryStats("Find", err)

	return err
}

// Count 获取记录数
func (qb *QueryBuilder) Count(count *int64) error {
	start := time.Now()
	err := qb.query.Count(count).Error
	utils.Debug("查询耗时: %v", time.Since(start))
	return err
}

// Create 创建记录
func (qb *QueryBuilder) Create(value interface{}) error {
	start := time.Now()
	err := qb.query.Create(value).Error
	utils.Debug("创建耗时: %v", time.Since(start))
	return err
}

// Updates 更新记录
func (qb *QueryBuilder) Updates(attrs interface{}) error {
	start := time.Now()
	defer func() {
		qb.logSlowQuery("Updates", time.Since(start))
	}()

	if qb.optimisticLock {
		// 添加乐观锁条件
		qb.query = qb.query.Clauses(clause.Locking{})
	}

	err := qb.query.Updates(attrs).Error
	if err != nil {
		return err
	}

	// 如果启用缓存，清除相关缓存
	if qb.enableCache && qb.cache != nil {
		// TODO: 实现缓存清理逻辑
	}

	return nil
}

// Delete 删除记录
func (qb *QueryBuilder) Delete(value interface{}) error {
	start := time.Now()
	defer func() {
		qb.logSlowQuery("Delete", time.Since(start))
	}()

	if qb.unscoped {
		// 物理删除
		err := qb.query.Unscoped().Delete(value).Error
		if err != nil {
			return err
		}
	} else {
		// 软删除
		err := qb.query.Delete(value).Error
		if err != nil {
			return err
		}
	}

	// 如果启用缓存，清除相关缓存
	if qb.enableCache && qb.cache != nil {
		// TODO: 实现缓存清理逻辑
	}

	return nil
}

// Transaction 事务
func (qb *QueryBuilder) Transaction(fc func(tx *gorm.DB) error) error {
	start := time.Now()
	err := qb.query.Transaction(fc)
	utils.Debug("事务耗时: %v", time.Since(start))
	return err
}

// Paginate 分页
func (qb *QueryBuilder) Paginate(page, pageSize int, dest interface{}) (total int64, err error) {
	start := time.Now()
	defer func() {
		utils.Debug("分页查询耗时: %v", time.Since(start))
	}()

	// 获取总数
	if err = qb.query.Count(&total).Error; err != nil {
		return 0, fmt.Errorf("获取总数失败: %v", err)
	}

	// 分页查询
	if err = qb.query.Offset((page - 1) * pageSize).Limit(pageSize).Find(dest).Error; err != nil {
		return 0, fmt.Errorf("分页查询失败: %v", err)
	}

	return total, nil
}

// BatchCreate 批量创建
func (qb *QueryBuilder) BatchCreate(values interface{}, batchSize int) error {
	start := time.Now()
	defer func() {
		qb.logSlowQuery("BatchCreate", time.Since(start))
	}()

	return qb.Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(values, batchSize).Error
	})
}

// BatchUpdate 批量更新
func (qb *QueryBuilder) BatchUpdate(values interface{}, batchSize int) error {
	start := time.Now()
	defer func() {
		qb.logSlowQuery("BatchUpdate", time.Since(start))
	}()

	return qb.Transaction(func(tx *gorm.DB) error {
		return tx.Save(values).Error
	})
}

// BatchDelete 批量删除
func (qb *QueryBuilder) BatchDelete(values interface{}) error {
	start := time.Now()
	defer func() {
		qb.logSlowQuery("BatchDelete", time.Since(start))
	}()

	return qb.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(values).Error
	})
}

// Raw 执行原生SQL
func (qb *QueryBuilder) Raw(sql string, values ...interface{}) *QueryBuilder {
	qb.query = qb.query.Raw(sql, values...)
	return qb
}

// Exec 执行SQL
func (qb *QueryBuilder) Exec(sql string, values ...interface{}) error {
	start := time.Now()
	err := qb.query.Exec(sql, values...).Error
	utils.Debug("执行SQL耗时: %v", time.Since(start))
	return err
}

// Scan 扫描结果
func (qb *QueryBuilder) Scan(dest interface{}) error {
	start := time.Now()
	err := qb.query.Scan(dest).Error
	utils.Debug("扫描结果耗时: %v", time.Since(start))
	return err
}

// Pluck 获取单个字段
func (qb *QueryBuilder) Pluck(column string, dest interface{}) error {
	start := time.Now()
	err := qb.query.Pluck(column, dest).Error
	utils.Debug("获取字段耗时: %v", time.Since(start))
	return err
}

// Debug 开启调试模式
func (qb *QueryBuilder) Debug() *QueryBuilder {
	qb.query = qb.query.Debug()
	return qb
}

// Unscoped 忽略软删除
func (qb *QueryBuilder) Unscoped() *QueryBuilder {
	qb.unscoped = true
	qb.query = qb.query.Unscoped()
	return qb
}

// Clauses 添加子句
func (qb *QueryBuilder) Clauses(conds ...clause.Expression) *QueryBuilder {
	qb.query = qb.query.Clauses(conds...)
	return qb
}

// Reset 重置查询
func (qb *QueryBuilder) Reset() *QueryBuilder {
	qb.query = qb.db.WithContext(qb.context)
	qb.conditions = make([]Condition, 0)
	qb.selectedFields = make([]string, 0)
	qb.unscoped = false
	qb.optimisticLock = false
	return qb
}

// WithOptimisticLock 启用乐观锁
func (qb *QueryBuilder) WithOptimisticLock() *QueryBuilder {
	qb.optimisticLock = true
	return qb
}

// GetConditions 获取查询条件
func (qb *QueryBuilder) GetConditions() []Condition {
	return qb.conditions
}

// GetSelectedFields 获取选择的字段
func (qb *QueryBuilder) GetSelectedFields() []string {
	return qb.selectedFields
}

// Explain 分析查询计划
func (qb *QueryBuilder) Explain() (string, error) {
	var result string
	err := qb.query.Raw("EXPLAIN " + qb.query.Statement.SQL.String()).Scan(&result).Error
	return result, err
}

// GetPoolStats 获取连接池统计信息
func (qb *QueryBuilder) GetPoolStats() (*sql.DBStats, error) {
	sqlDB, err := qb.db.DB()
	if err != nil {
		return nil, err
	}
	stats := sqlDB.Stats()
	return &stats, nil
}

// MonitorPool 监控连接池
func (qb *QueryBuilder) MonitorPool() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			stats, err := qb.GetPoolStats()
			if err != nil {
				utils.Error("获取连接池统计信息失败: %v", err)
				continue
			}

			// 记录连接池统计信息
			utils.Debug("连接池统计 - 打开连接数: %d, 使用中连接数: %d, 空闲连接数: %d, 等待连接数: %d",
				stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)

			// 检查连接池健康状态
			if stats.OpenConnections >= qb.poolConfig.MaxOpenConns*8/10 {
				utils.Warn("连接池接近最大连接数 - 打开连接数: %d, 最大连接数: %d",
					stats.OpenConnections, qb.poolConfig.MaxOpenConns)
			}

			if stats.WaitCount > 0 {
				utils.Warn("连接池有等待连接 - 等待连接数: %d", stats.WaitCount)
			}
		}
	}()
}

// logQueryStats 记录查询统计
func (qb *QueryBuilder) logQueryStats(operation string, err error) {
	if qb.stats != nil {
		qb.stats.SQL = qb.query.Statement.SQL.String()
		qb.stats.Timestamp = time.Now()
		qb.stats.Error = err
		qb.stats.IsSlow = qb.stats.Duration > qb.slowQueryThreshold

		// 记录慢查询
		if qb.stats.IsSlow {
			utils.Warn("慢查询警告 - 操作: %s, SQL: %s, 耗时: %v, 行数: %d",
				operation, qb.stats.SQL, qb.stats.Duration, qb.stats.Rows)
		}

		// 记录查询统计
		utils.Debug("查询统计 - 操作: %s, SQL: %s, 耗时: %v, 行数: %d",
			operation, qb.stats.SQL, qb.stats.Duration, qb.stats.Rows)
	}
}

// logSlowQuery 记录慢查询
func (qb *QueryBuilder) logSlowQuery(operation string, duration time.Duration) {
	if duration > qb.slowQueryThreshold {
		utils.Warn("慢查询警告 - 操作: %s, 耗时: %v", operation, duration)
	}
}
