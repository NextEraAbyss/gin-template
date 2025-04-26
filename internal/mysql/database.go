package mysql

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitee.com/NextEraAbyss/gin-template/config"
	"gitee.com/NextEraAbyss/gin-template/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 数据库连接
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(config *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	// 设置自定义Logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢查询阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色输出
		},
	)

	// MySQL连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database connection: %v", err))
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置连接的最大可复用时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 迁移数据库（开发环境自动迁移）
	if config.Env != "production" {
		err = db.AutoMigrate(
			&models.User{},
			// 添加其他模型...
		)
		if err != nil {
			log.Printf("数据库迁移失败: %v", err)
		}
	}

	// 设置全局变量
	DB = db

	return db
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// Transaction 事务处理
func Transaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(fn)
}
