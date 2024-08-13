package datasource

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

type MySQLDatasource struct {
	engine *gorm.DB
}

func NewMySQLDatasource() (Datasource, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       envMySQL(), // DSN data source name
		DefaultStringSize:         256,        // string 类型字段的默认长度
		DisableDatetimePrecision:  true,       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,      // 根据当前 MySQL 版本自动配置

	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	if !db.Migrator().HasTable("t_config") {
		// 自动创建和更新表结构
		err = db.Migrator().CreateTable(&Config{})
		if err != nil {
			log.Panicf("Failed to auto migrate database: %v", err)
		}
	}
	return &MySQLDatasource{engine: db}, nil
}

func envMySQL() string {
	return os.Getenv("MYSQL_URL")
}

func (m *MySQLDatasource) Get(ctx context.Context, key, field string) (string, error) {
	config := Config{}
	err := m.engine.WithContext(ctx).Where("c_key=? and c_field=?", key, field).First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}
	return config.Value, nil
}

func (m *MySQLDatasource) Set(ctx context.Context, key, field string, value string) error {
	return m.engine.WithContext(ctx).Create(&Config{
		Key:   key,
		Field: field,
		Value: value,
	}).Error
}

func (m *MySQLDatasource) Delete(ctx context.Context, key, field string) error {
	err := m.engine.WithContext(ctx).Where("c_key=? and c_field=?", key, field).Delete(&Config{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (m *MySQLDatasource) DeleteAll(ctx context.Context, key string) error {
	err := m.engine.WithContext(ctx).Where("c_key=?", key).Delete(&Config{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (m *MySQLDatasource) Select(ctx context.Context, key string, options *SelectOptions) ([]string, error) {
	keys := make([]string, 0)
	configs := make([]Config, 0)
	err := m.engine.WithContext(ctx).Where("c_key=?", key).Limit(int(options.Count)).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	for _, k := range configs {
		keys = append(keys, k.Value)
	}
	return keys, err
}

func (m *MySQLDatasource) SelectAll(ctx context.Context, key string) (map[string]string, error) {
	keys := make(map[string]string)
	configs := make([]Config, 0)
	err := m.engine.WithContext(ctx).Where("c_key=?", key).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	for _, k := range configs {
		keys[k.Field] = k.Value
	}
	return keys, err
}

func (m *MySQLDatasource) Count(ctx context.Context, key string) (int64, error) {
	var count int64
	err := m.engine.WithContext(ctx).Model(&Config{}).Where("c_key=?", key).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *MySQLDatasource) Incr(ctx context.Context, key, field string, value int64) error {
	config := Config{}
	tx := m.engine.Begin()
	defer tx.Commit()
	err := tx.WithContext(ctx).
		//Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("c_key=? and c_field=?", key, field).First(&config).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if config.Value == "" {
		return nil
	}
	parseInt, err := strconv.ParseInt(config.Value, 10, 64)
	if err != nil {
		return err
	}
	return tx.WithContext(ctx).Model(&Config{}).
		Where("c_key=? and c_field=?", key, field).
		Update("c_value", fmt.Sprintf("%d", parseInt+value)).Error
}
