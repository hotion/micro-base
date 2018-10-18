package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/shiguanghuxian/micro-base/internal/config"
	"github.com/shiguanghuxian/micro-base/internal/log"
	"github.com/shiguanghuxian/micro-base/internal/register"
	"github.com/shiguanghuxian/micro-base/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/kit/middleware/logging"
	"github.com/shiguanghuxian/micro-base/kit/service"
	"github.com/shiguanghuxian/micro-base/kit/transport"
	"github.com/shiguanghuxian/micro-base/pb"
	"google.golang.org/grpc"
)

var (
	svcName  = config.GetSvcName()
	VERSION  string // 程序版本
	GIT_HASH string // git hash
)

func main() {
	// 日志
	logger := log.Logger

	// 服务对象
	var s service.Service
	{
		s = service.NewBasicService()
		// 注册日志中间件
		s = logging.LoggingMiddleware(logger)(s)
	}

	// 创建所有端点
	endpoints := endpoint.MakeServerEndpoints(s)

	// transport grpc方式
	var grpcServer pb.HelloServer
	{
		grpcServer = transport.NewGRPCServer(endpoints, s, logger)
	}

	// 启动grpc监听
	go func() {
		grpcListener, err := net.Listen("tcp", config.GetGRPCAddr())
		if err != nil {
			logger.Panicf("TCP监听错误:%v", err)
			return
		}
		logger.Infof("TCP监听成功:%s", config.GetGRPCAddr())

		baseServer := grpc.NewServer()
		pb.RegisterHelloServer(baseServer, grpcServer)

		err = baseServer.Serve(grpcListener)
		if err != nil {
			logger.Panicf("TCP监听错误:%v", err)
			return
		}
	}()

	// 注册服务
	go register.Register(config.GetETCDAddr(), svcName, config.GetGRPCAddr(), 5)

	// 监听退出信号
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-ch
	register.UnRegister(svcName, config.GetGRPCAddr())
	if i, ok := sig.(syscall.Signal); ok {
		os.Exit(int(i))
	} else {
		os.Exit(0)
	}
}
