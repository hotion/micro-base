package config

import "os"

// 全局配置
const (
	// DEFAULT_ETCD_ADDR 默认etcd地址
	DEFAULT_ETCD_ADDR string = "127.0.0.1:2379"
	// DEFAULT_GRPC_ADDR 默认grpc监听地址
	DEFAULT_GRPC_ADDR string = "127.0.0.1:28080"
	// DEFAULT_SVC_NAME 默认服务名
	DEFAULT_SVC_NAME string = "project/default"
	// 当前运行环境，dev or pro
	DEFAULT_MODE string = "dev"
)

// GetETCDAddr 读取etcd服务地址
func GetETCDAddr() string {
	etcdAddr := os.Getenv("ETCD_ADDR")
	if etcdAddr == "" {
		return DEFAULT_ETCD_ADDR
	}
	return etcdAddr
}

// GetGRPCAddr 读取grpc地址
func GetGRPCAddr() string {
	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		return DEFAULT_GRPC_ADDR
	}
	return grpcAddr
}

// GetSvcName 获取服务名
func GetSvcName() string {
	svcName := os.Getenv("SVC_NAME")
	if svcName == "" {
		return DEFAULT_SVC_NAME
	}
	return svcName
}

// GetMode 当前运行环境
func GetMode() string {
	mode := os.Getenv("MODE")
	if mode == "" {
		return DEFAULT_MODE
	}
	return mode
}
