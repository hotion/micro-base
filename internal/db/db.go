package db

import (
	"errors"
	"fmt"
	"time"

	"atisafe.com/tools/logger"
	"github.com/jinzhu/gorm"
	"github.com/shiguanghuxian/micro-base/internal/config"
	"github.com/shiguanghuxian/micro-base/internal/etcdcli"
)

var (
	// MasterDB db主节点
	MasterDB *gorm.DB
	// SlaveDB db从节点
	SlaveDB *gorm.DB
)

// 初始化主节点
func init() {
	// 初始化MasterDB
	cfg, err := config.GetDBConfig(etcdcli.EtcdCli, "master")
	if err != nil {
		// ### log
	}
	logger.Infoln(logger.INFO_MySQL_CONN, "db conn master")
	// 创建数据库连接
	MasterDB, err = NewDbClient(cfg)
	if err != nil {
		// ### log
	}
	logger.Infoln(logger.INFO_MySQL_CONN, "db conn master end")
	// 初始化SlaveDB
	cfg, err = config.GetDBConfig(etcdcli.EtcdCli, "slave")
	if err != nil {
		// ### log
	}
	// 创建数据库连接
	SlaveDB, err = NewDbClient(cfg)
	if err != nil {
		// ### log
	}
}

// NewDbClient 创建数据库连接
func NewDbClient(cfg *config.DbConfig) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errors.New("The database configuration file can not be empty.")
	}
	// 拼接连接数据库字符串
	connStr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.User,
		cfg.Passwd,
		cfg.Address,
		cfg.Port,
		cfg.DbName)
	// 连接数据库
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	// 设置表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return cfg.Prefix + defaultTableName
	}
	// 禁用表名多元化
	db.SingularTable(true)
	// 是否开启debug模式
	if cfg.Debug {
		// debug 模式
		db = db.Debug()
	}
	// 连接池最大连接数
	db.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	// 默认打开连接数
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns)
	// 开启协程ping MySQL数据库查看连接状态
	go func() {
		for {
			// ping
			err = db.DB().Ping()
			if err != nil {
				logger.Infoln(logger.ERROR_MySQL_CONN, err)
			}
			// 间隔30s ping一次
			time.Sleep(time.Second * 30)
		}
	}()

	return db, err
}
