package database

import (
	"fmt"
	"log"
	"time"

	"gitee.com/NextEraAbyss/gin-template/config"
	"gitee.com/NextEraAbyss/gin-template/models"
	"gitee.com/NextEraAbyss/gin-template/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init(config *config.Config) *gorm.DB {
	// MySQL DSN格式: username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	// 设置日志级别
	var logLevel logger.LogLevel
	if config.Env == "production" {
		logLevel = logger.Error
	} else {
		logLevel = logger.Info
	}

	// 自定义日志配置
	customLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 200 * time.Millisecond, // 慢查询阈值
			LogLevel:      logLevel,               // 日志级别
			Colorful:      true,                   // 彩色打印
		},
	)

	// 配置GORM选项
	gormConfig := &gorm.Config{
		Logger: customLogger,
		// 禁用默认事务
		SkipDefaultTransaction: true,
		// 命名策略
		NamingStrategy: nil, // 默认表名和字段名命名规则
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB instance: %v", err)
	}

	// 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 设置连接最大空闲时间
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	// 自动迁移数据库表
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建默认用户
	createDefaultUsers(db)

	DB = db
	return db
}

// createDefaultUsers 创建默认用户（仅当数据库中没有用户时）
func createDefaultUsers(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)

	// 只有当数据库中没有用户时才创建默认用户
	if count == 0 {
		log.Println("Creating default users...")

		// 默认用户数据
		defaultUsers := []models.User{
			{
				Username: "admin",
				Email:    "admin@example.com",
				Password: "Admin123!", // 将在保存前加密
			},
			{
				Username: "user1",
				Email:    "user1@example.com",
				Password: "Password123!",
			},
			{
				Username: "user2",
				Email:    "user2@example.com",
				Password: "Password123!",
			},
		}

		// 加密密码并保存用户
		for _, user := range defaultUsers {
			// 密码加密
			hashedPassword, err := utils.HashPassword(user.Password)
			if err != nil {
				log.Printf("Failed to hash password for user %s: %v", user.Username, err)
				continue
			}
			user.Password = hashedPassword

			// 保存用户
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create default user %s: %v", user.Username, err)
			} else {
				log.Printf("Default user created: %s", user.Username)
			}
		}
	}
}
