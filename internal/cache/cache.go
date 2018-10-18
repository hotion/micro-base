package cache

import (
	"errors"
	"time"

	"atisafe.com/tools/logger"
	"github.com/shiguanghuxian/micro-base/internal/config"
	"github.com/shiguanghuxian/micro-base/internal/etcdcli"
	goredis "gopkg.in/redis.v5"
)

/* 缓存连接 redis */

var (
	// Client redis 连接对象
	Client goredis.Cmdable
)

// 初始化redis
func init() {
	cfg, err := config.GetRedisConfg(etcdcli.EtcdCli)
	if err != nil {
		// ### log
	}
	Client, err = NewClient(cfg)
	if err != nil {
		// ### log
	}
}

// NewClient 创建客户端连接
func NewClient(cfg *config.RedisConfg) (client goredis.Cmdable, err error) {
	if cfg == nil {
		err = errors.New("The redis configuration file can not be empty.")
		return
	}
	logger.Infoln(logger.INFO_REDIS_CONN, "start")
	if cfg.IsCluster == true {
		// redis集群
		client = goredis.NewClusterClient(&goredis.ClusterOptions{
			Addrs:    cfg.Address,
			Password: cfg.Password,
			PoolSize: cfg.PoolSize,
		})
	} else {
		// redis单机
		client = goredis.NewClient(&goredis.Options{
			Addr:     cfg.Address[0],
			Password: cfg.Password,
			DB:       cfg.Db,
			PoolSize: cfg.PoolSize,
		})
	}
	// ping 防止断开
	go func() {
		for {
			err := client.Ping().Err()
			if err != nil {
				logger.Errorln(logger.ERROR_REDIS_CONN, err)
			}
			time.Sleep(time.Second * 30)
		}
	}()
	logger.Infoln(logger.INFO_REDIS_CONN, "end")
	return
}
