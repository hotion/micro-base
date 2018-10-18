package config

import (
	"context"
	"errors"

	"github.com/coreos/etcd/clientv3"
	"github.com/naoina/toml"
)

/* 读取各种配置 */

/* 读取mysql配置 */

// DbConfig 数据库配置
type DbConfig struct {
	Debug        bool   `toml:"debug"`          // 是否调试模式
	Address      string `toml:"address"`        // 数据库连接地址
	Port         int    `toml:"port"`           // 数据库端口
	MaxIdleConns int    `toml:"max_idle_conns"` // 连接池最大连接数
	MaxOpenConns int    `toml:"max_open_conns"` // 默认打开连接数
	User         string `toml:"user"`           // 数据库用户名
	Passwd       string `toml:"passwd"`         // 数据库密码
	DbName       string `toml:"db_name"`        // 数据库名
	Prefix       string `toml:"prefix"`         // 数据库表前缀
}

// GetDBConfig 获取数据库配置
func GetDBConfig(cli *clientv3.Client, node string) (*DbConfig, error) {
	if node != "master" && node != "slave" {
		return nil, errors.New("Node can only be master or slave.")
	}
	// 初始化master连接
	etcdResp, err := cli.Get(context.Background(), "home/config/mysql/account/"+node+".toml", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if len(etcdResp.Kvs) == 0 {
		return nil, errors.New("The database configuration MySQL master node is configured to be empty.")
	}
	dbConfig := new(DbConfig)
	err = toml.Unmarshal(etcdResp.Kvs[0].Value, dbConfig)
	if err != nil {
		return nil, err
	}
	return dbConfig, nil
}

/* redis配置 */

// RedisConfg redis配置
type RedisConfg struct {
	Address   []string `toml:"address"`    // redis 服务器地址,包括地址和端口 127.0.0.1:6379
	Password  string   `toml:"password"`   // redis 密码
	Db        int      `toml:"db"`         // 连接的数据库
	PoolSize  int      `toml:"pool_size"`  // 连接池大小
	IsCluster bool     `toml:"is_cluster"` // 是否集群模式
}

// GetRedisConfg 获取redis配置
func GetRedisConfg(cli *clientv3.Client) (*RedisConfg, error) {
	// 初始化master连接
	etcdResp, err := cli.Get(context.Background(), "home/config/redis/cfg.toml", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if len(etcdResp.Kvs) == 0 {
		return nil, errors.New("The redis configuration configured to be empty.")
	}
	redisConfg := new(RedisConfg)
	err = toml.Unmarshal(etcdResp.Kvs[0].Value, redisConfg)
	if err != nil {
		return nil, err
	}
	return redisConfg, nil
}
